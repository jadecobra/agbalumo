import csv
import re
import os

def update_csv_urls():
    base_dir = '/Users/johnnyblase/gym/agbalumo/.tester/scripts/enrich_urls'
    website_csv_path = os.path.join(base_dir, 'rccgna_website.csv')
    enriched_csv_path = os.path.join(base_dir, 'rccgna_enriched.csv')
    
    # We will overwrite the original enriched file safely by writing to a temp file first.
    output_csv_path = os.path.join(base_dir, 'rccgna_enriched_updated.csv')

    # Parse URL mapping
    mapping = {}
    
    # Matches lines like:
    # [74]  The Master Builder's Parish
    # "[86]  CITY OF PRAISE, CARROLLTON"
    # Group 1 captures the title
    title_pattern = re.compile(r'^\"?\[\d+\]\s+(.*?)\"?$')
    current_title = None

    print("Parsing website.csv for URLs...")
    with open(website_csv_path, 'r', encoding='utf-8') as f:
        for line in f:
            line_str = line.strip()
            if not line_str:
                continue
            
            m = title_pattern.match(line_str)
            if m:
                # Found a title header block
                current_title = m.group(1).strip()
                continue
                
            # If line starts with "http", it is the website URL for the last parsed title
            if line_str.startswith('http') and current_title:
                # Remove quotes or extra spaces if present
                url = line_str.strip('"\' ')
                mapping[current_title] = url
                # Reset title so we don't accidentally map it to following URLs
                current_title = None

    print(f"Extracted {len(mapping)} title-to-URL mappings.")

    # Update enriched CSV
    print("Updating enriched.csv...")
    updated_count = 0
    with open(enriched_csv_path, 'r', encoding='utf-8') as infile, \
         open(output_csv_path, 'w', encoding='utf-8', newline='') as outfile:
         
        reader = csv.DictReader(infile)
        fieldnames = reader.fieldnames
        writer = csv.DictWriter(outfile, fieldnames=fieldnames)
        writer.writeheader()
        
        for row in reader:
            title = row.get('title', '').strip()
            if title in mapping:
                row['website'] = mapping[title]
                updated_count += 1
            writer.writerow(row)

    # Rename the updated file to replace the old one
    os.replace(output_csv_path, enriched_csv_path)
    print(f"✅ Successfully matched and updated {updated_count} websites in rccgna_enriched.csv.")

if __name__ == '__main__':
    update_csv_urls()
