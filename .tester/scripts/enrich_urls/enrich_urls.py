import csv
import os
import re
import time
import argparse
import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry
from urllib.parse import urlparse
from bs4 import BeautifulSoup

# --- Global Session with Retries ---
session = requests.Session()
retries = Retry(total=3, backoff_factor=1, status_forcelist=[ 500, 502, 503, 504 ])
session.mount('http://', HTTPAdapter(max_retries=retries))
session.mount('https://', HTTPAdapter(max_retries=retries))

# --- Configuration ---
BRAVE_API_KEY = os.environ.get("BRAVE_API_KEY", "")
BRAVE_SEARCH_URL = "https://api.search.brave.com/res/v1/web/search"
EXCLUDED_DOMAINS = [
    "yelp.com", "yellowpages.com", "linkedin.com", "mapquest.com", 
    "alignable.com", "eventbrite.com", "rccgna.com", "afrifinder.com"
]
TIMEOUT_SECONDS = 10

def extract_email_from_fallback(fallback_url):
    """Fetches the fallback RCCG directory URL and extracts an email if present."""
    return None
    try:
        headers = {'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'}
        response = session.get(fallback_url, headers=headers, timeout=TIMEOUT_SECONDS)
        response.raise_for_status()
        
        # Regex for standard emails
        emails = set(re.findall(r'[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}', response.text))
        
        # Remove common false positives (like webmaster, info from framework, etc) - though any email helps
        if emails:
            return list(emails)[0]
    except Exception as e:
        print(f"  [Fallback Error] {e}")
    return None

# --- Normalization Helpers ---
def normalize_phone(phone_str):
    """Strip all non-numeric characters from a string."""
    if not phone_str:
        return ""
    return re.sub(r'\D', '', phone_str)

def normalize_text(text_str):
    """Lowercase and strip whitespace for flexible comparison."""
    if not text_str:
        return ""
    return str(text_str).lower().strip()

# --- Verification Logic ---
def verify_website(candidate_url, title, address, phone):
    """
    Fetches the candidate URL and checks if the page content contains the parish's phone, title, or address.
    Returns (is_verified, reason).
    """
    try:
        # Some servers reject requests without a standard User-Agent
        headers = {'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'}
        response = session.get(candidate_url, headers=headers, timeout=TIMEOUT_SECONDS)
        response.raise_for_status() 
        html_content = response.text
        soup = BeautifulSoup(html_content, 'html.parser')
        page_text = soup.get_text(separator=' ')
        
        norm_page_text = normalize_text(page_text)
        norm_phone = normalize_phone(phone)
        norm_title = normalize_text(title)
        
        # Phone check (strongest signal)
        if len(norm_phone) >= 10 and norm_phone in normalize_phone(page_text):
            return True, "Phone Match"
        
        # Title check (weaker signal)
        title_tag = soup.find('title')
        page_title = normalize_text(title_tag.string) if title_tag else ""
        
        if norm_title and norm_title in page_title:
             return True, "Title Match"
             
        # Address check (zip code is usually a solid signal if present)
        # Extract zip code from address if possible
        zip_match = re.search(r'\b\d{5}(?:-\d{4})?\b', address)
        if zip_match:
            zip_code = zip_match.group(0)
            if zip_code in norm_page_text:
                if norm_title in norm_page_text:
                    return True, "Title + Zip Match"
                
        return False, "No strong signals matched on page"

    except requests.exceptions.RequestException as e:
        return False, f"Request failed: {e}"
    except Exception as e:
         return False, f"Verification error: {e}"

# --- Search API Logic ---
def search_for_parish(title, address, phone):
    """
    Query Brave Search API for the parish.
    Return the best candidate URL.
    """
    if not BRAVE_API_KEY:
        print("ERROR: BRAVE_API_KEY environment variable is not set. Cannot perform search.")
        return None, None
        
    # Build a tight query. Prepend RCCG if not already in the title
    search_title = title if "rccg" in title.lower() else f"RCCG {title}"
    query_parts = [search_title]
    if phone:
        query_parts.append(phone)
    
    # Try to extract city from address (assuming format: City, State Zip)
    city_match = re.search(r'\n(.*?),\s*[A-Za-z]+', address)
    if city_match:
        city = city_match.group(1).strip()
        query_parts.append(city)

    query = " ".join(query_parts)
    print(f"  [Search] Query: {query}")
    
    headers = {
        "Accept": "application/json",
        "Accept-Encoding": "gzip",
        "X-Subscription-Token": BRAVE_API_KEY
    }
    
    try:
        response = session.get(BRAVE_SEARCH_URL, headers=headers, params={"q": query, "count": 5}, timeout=TIMEOUT_SECONDS)
        
        if response.status_code == 429:
             print("  [Search] RATE LIMIT (429) HIT. Backing off.")
             # Raise an exception so the main loop can handle the backoff/exit
             response.raise_for_status()
             
        response.raise_for_status()
        data = response.json()
        
        primary_url = None
        fallback_url = None
        
        if "web" in data and "results" in data["web"]:
            for result in data["web"]["results"]:
                url = result.get("url", "")
                parsed_url = urlparse(url)
                
                # Check for excluded directory domains or fallback rccgna server
                is_excluded = any(domain in parsed_url.netloc for domain in EXCLUDED_DOMAINS)
                is_fallback = "secure.rccgna.org" in parsed_url.netloc or "www.rccgna.org" in parsed_url.netloc
                
                if url and not is_excluded and not is_fallback:
                     if not primary_url:
                          primary_url = url
                if url and is_fallback:
                     if not fallback_url:
                          fallback_url = url
                          
                if primary_url and fallback_url:
                     break
                     
        return primary_url, fallback_url

    except requests.exceptions.RequestException as e:
        print(f"  [Search] API request failed: {e}")
        # Re-raise 429 so the main loop can catch it
        if isinstance(e, requests.exceptions.HTTPError) and e.response.status_code == 429:
             raise e
        return None, None

# --- Main Processor ---
def process_csv(input_file, output_file, delay=1.5):
    """
    Reads input_file, searches/verifies URLs, and writes incrementally to output_file.
    """
    print(f"Starting processing. Input: {input_file}, Output: {output_file}")
    
    # Read existing output to resume progress
    processed_titles = set()
    if os.path.exists(output_file):
        with open(output_file, 'r', encoding='utf-8') as f:
            reader = csv.reader(f)
            next(reader, None)  # Skip header
            for row in reader:
                if len(row) > 0:
                    processed_titles.add(row[0]) # Assuming title is column 0
        print(f"Found {len(processed_titles)} already processed rows in output file. Skipping those.")
    
    # Open input and output files
    with open(input_file, 'r', encoding='utf-8') as infile:
        reader = csv.reader(infile)
        header = next(reader)
        
        # Ensure we have the right header structure (A: Title, C: Address, D: Phone, G: Website)
        # title,description,address,phone,type,origin,website -> 7 columns
        if len(header) < 7:
             header.extend([""] * (7 - len(header)))
             if header[6] != "website":
                  header[6] = "website"
        
        # Open output in append mode
        mode = 'a' if os.path.exists(output_file) else 'w'
        with open(output_file, mode, encoding='utf-8', newline='') as outfile:
            writer = csv.writer(outfile)
            
            if mode == 'w':
                # Write header if new file, plus a verification reason column
                new_header = header[:7] + ["verification_status"]
                writer.writerow(new_header)
            
            row_count = 0
            for row in reader:
                row_count += 1
                
                # Pad row if necessary to match the 7 expected columns
                if len(row) < 7:
                    row.extend([""] * (7 - len(row)))
                    
                title = row[0]
                address = row[2]
                phone = row[3]
                existing_website = row[6]
                
                if title in processed_titles:
                    continue # Skip already processed
                    
                print(f"[{row_count}] Processing: {title}")
                
                # If they already explicitly have a website from the source CSV, we could verify it directly
                # But the prompt asks to find the URL based on A, C, D and place in G.
                verification_reason = "Not Found"
                
                try:
                     primary_url, fallback_url = search_for_parish(title, address, phone)
                     
                     if primary_url:
                          parsed_primary = urlparse(primary_url)
                          is_social = any(domain in parsed_primary.netloc for domain in ["facebook.com", "instagram.com"])
                          
                          if is_social:
                               print(f"  [Found Candidate] {primary_url} - (Bypassing verification logic for Social Media)")
                               row[6] = primary_url
                               verification_reason = "Accepted Social Media Profile"
                          else:
                               print(f"  [Found Candidate] {primary_url} - Verifying...")
                               is_valid, reason = verify_website(primary_url, title, address, phone)
                               
                               if is_valid:
                                    print(f"  [VERIFIED] Reason: {reason}")
                                    row[6] = primary_url
                                    verification_reason = f"Verified: {reason}"
                               else:
                                    print(f"  [NOT VERIFIED] Reason: {reason}")
                                    verification_reason = f"Failed Verification: {reason}"
                     else:
                          print("  [No Candidate found]")
                          
                     # Fallback logic if we haven't verified a regular site
                    #  if row[6] in (existing_website, "") and fallback_url:
                    #       print(f"  [Fallback] Checking RCCGNA directory for email: {fallback_url}")
                    #       time.sleep(delay) # Prevent rate limiting on the target server
                    #       email = extract_email_from_fallback(fallback_url)
                    #       if email:
                    #            print(f"  [Fallback SUCCESS] Found email: {email}")
                    #            row[6] = email
                    #            verification_reason = "Fallback: Found Email on RCCGNA directory"
                    #       else:
                    #            print("  [Fallback FAILED] Could not find email on page.")
                          
                except requests.exceptions.HTTPError as e:
                     if e.response.status_code == 429:
                          print("\n*** RATE LIMIT EXCEEDED ***")
                          print("Saving current progress and exiting cleanly. Re-run script to resume later.")
                          return
                     else:
                          print(f"  [Error] {e}")
                except Exception as e:
                     print(f"  [Unexpected Error] {e}")
                
                # Write individual row to output CSV to checkpoint progress
                out_row = row[:7] + [verification_reason]
                writer.writerow(out_row)
                outfile.flush() # Force write to disk
                
                # Throttle API requests
                time.sleep(delay)
                
    print("Processing complete!")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Enrich RCCG Parish CSV with verified website URLs.")
    parser.add_argument("--input", required=True, help="Path to input CSV file")
    parser.add_argument("--output", required=True, help="Path to output CSV file")
    parser.add_argument("--delay", type=float, default=1.5, help="Delay between searches in seconds (to respect rate limits)")
    
    args = parser.parse_args()
    
    process_csv(args.input, args.output, args.delay)
