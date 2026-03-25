import csv
import re

def clean_address(address):
    # 1. Add comma and space before "USA" or "United States"
    address = re.sub(r'(?<!,\s)(USA|United States)', r', \1', address)
    
    # 2. Add comma and space before State ZIP if smashed
    # Look for [Lowercase][OptionalSpace][StateAbbr][Space][ZIP]
    address = re.sub(r'([a-z])\s*([A-Z]{2})\s*(\d{5})', r'\1, \2 \3', address)
    
    # 3. Handle smashed "RoadSuite", "DriveSuite" etc.
    address = re.sub(r'(Road|Drive|Street|Hwy|Apt\.|Suite|Unit)\s*(\d+|[A-Z])', r'\1 \2', address)
    
    # 4. Add comma before identifiers like Suite, Unit, etc. if following a lowercase letter
    address = re.sub(r'([a-z])\s*(Suite|Unit|Apt\.|Courtyard|Neighborhood|Marriott|Budgeted|Quality|Millbrook)', r'\1, \2', address)

    # 5. Clean up any " , " to ", "
    address = re.sub(r'\s+,\s+', r', ', address)
    address = re.sub(r'\s+,', r',', address)
    address = re.sub(r',([^\s])', r', \1', address)

    return address.strip()

def parse_christ_embassy(input_file, output_file):
    with open(input_file, 'r') as f:
        lines = f.readlines()

    records = []
    current_record = None

    for line in lines:
        line = line.strip()
        if not line:
            continue

        # Check for start of a new record
        if line.upper().startswith("CHRIST EMBASSY"):
            if current_record:
                records.append(current_record)
            
            # Split title and address at the first digit
            match = re.search(r'\d', line)
            if match:
                title = line[:match.start()].strip()
                # Special case: title ends with location which might be smashed with street number
                address = line[match.start():].strip()
            else:
                title = line
                address = ""
            
            current_record = {
                "title": title.title(), # Title case for better look
                "address": clean_address(address),
                "phone": "",
                "email": ""
            }
        elif line.startswith("Phone:"):
            # Extract phone and email
            # Sometimes Phone: is followed by number then Email:
            phone_match = re.search(r'Phone:\s*([\d\s+-]+?)(?=Email:|$)', line)
            if phone_match:
                current_record["phone"] = phone_match.group(1).strip()
            
            email_match = re.search(r'Email:\s*([^\s]+)', line)
            if email_match:
                current_record["email"] = email_match.group(1).strip()
        elif "miDirections" in line or "More info" in line:
            continue
        elif current_record and not current_record["address"] and not line.startswith("Phone"):
             current_record["address"] = clean_address(line)

    if current_record:
        records.append(current_record)

    # Write to CSV
    with open(output_file, 'w', newline='') as f:
        writer = csv.DictWriter(f, fieldnames=["title", "address", "phone", "email"])
        writer.writeheader()
        writer.writerows(records)

if __name__ == "__main__":
    parse_christ_embassy('.tester/scripts/enrich_urls/christ_embassy.csv', '.tester/scripts/enrich_urls/christ_embassy_organized.csv')
