
import sys
import re

def parse_logs(file_path):
    with open(file_path, 'r') as f:
        content = f.read()
    
    # Extract the clone groups section
    matches = re.findall(r'found (\d+) clones:\n((?:\s+.*\d+,\d+\n?)+)', content)
    
    group_stats = []
    
    for count, instances_str in matches:
        count = int(count)
        instances = instances_str.strip().split('\n')
        
        # Parse the range from the first instance (assuming they are same length)
        first_instance = instances[0].strip()
        m = re.search(r':(\d+),(\d+)$', first_instance)
        if m:
            start, end = int(m.group(1)), int(m.group(2))
            length = end - start + 1
            total_cloned_lines = length * count
            
            group_stats.append({
                'count': count,
                'length': length,
                'total_cloned_lines': total_cloned_lines,
                'instances': instances
            })
            
    # Sort by total cloned lines descending
    group_stats.sort(key=lambda x: x['total_cloned_lines'], reverse=True)
    
    print(f"Top 5 clone groups by total lines (length * count):")
    for i, g in enumerate(group_stats[:5]):
        print(f"{i+1}. Total: {g['total_cloned_lines']} lines ({g['length']} lines across {g['count']} clones)")
        for inst in g['instances']:
            print(f"   - {inst.strip()}")
        print()

    # Aggregate by file
    file_stats = {}
    for g in group_stats:
        for inst in g['instances']:
            # Extract file path
            m = re.match(r'^\s*([^:]+)', inst)
            if m:
                path = m.group(1).split(':')[0].strip() # remove line numbers if any
                file_stats[path] = file_stats.get(path, 0) + g['length']
    
    sorted_files = sorted(file_stats.items(), key=lambda x: x[1], reverse=True)
    print("Top 10 files by total cloned lines:")
    for path, lines in sorted_files[:10]:
        print(f"- {path}: {lines} lines")

if __name__ == "__main__":
    parse_logs('/Users/johnnyblase/gym/agbalumo/ci_full_log.txt')
