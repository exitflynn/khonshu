import os
import argparse
import fnmatch

def parse_gitignore(gitignore_path):
    ignore_patterns = []
    try:
        with open(gitignore_path, 'r') as f:
            for line in f:
                line = line.strip()
                if line and not line.startswith('#'):
                    ignore_patterns.append(line)
    except FileNotFoundError:
        pass
    return ignore_patterns

def should_ignore(path, ignore_dirs, ignore_exts, gitignore_patterns):
    file_name = os.path.basename(path)
    if file_name.startswith('.'):
        return True
    
    path_parts = path.split(os.path.sep)
    for ignored_dir in ignore_dirs:
        if ignored_dir in path_parts:
            return True
    
    for ext in ignore_exts:
        if path.endswith(ext):
            return True
    
    for pattern in gitignore_patterns:
        if fnmatch.fnmatch(path, pattern):
            return True
    
    return False

def generate_directory_structure(source_path, output_path, ignore_dirs, ignore_exts):
    if not os.path.exists(source_path):
        print(f"Error: Source path {source_path} does not exist.")
        return
    
    gitignore_path = os.path.join(source_path, '.gitignore')
    gitignore_patterns = parse_gitignore(gitignore_path)
    
    abs_path = os.path.abspath(source_path)
    root_name = os.path.basename(abs_path)
    if root_name == '.':
        root_name = os.path.basename(os.getcwd())
    structure = [root_name]
    
    def build_tree(path, prefix=''):
        items = sorted(os.listdir(path))
        
        dirs = []
        files = []
        
        for item in items:
            full_item_path = os.path.join(path, item)
            if os.path.isdir(full_item_path):
                dirs.append(item)
            elif os.path.isfile(full_item_path):
                files.append(item)
        
        all_items = dirs + files
        
        for i in range(len(all_items)):
            item = all_items[i]
            full_path = os.path.join(path, item)
            
            if should_ignore(full_path, ignore_dirs, ignore_exts, gitignore_patterns):
                continue
            
            is_last = (i == len(all_items) - 1)
            connector = '└── ' if is_last else '├── '
            
            if os.path.isdir(full_path):
                structure.append(f"{prefix}{connector}{item}")
                if is_last:
                    build_tree(full_path, prefix + '    ')
                else:
                    build_tree(full_path, prefix + '│   ')
            else:
                structure.append(f"{prefix}{connector}- {item}")
    
    build_tree(source_path)
    
    try:
        with open(output_path, 'w') as f:
            f.write('\n'.join(structure))
        print(f"Project structure written to {output_path}")
    except IOError as e:
        print(f"Error writing to {output_path}: {e}")

def main():
    parser = argparse.ArgumentParser(description='Generate project directory structure in markdown')
    parser.add_argument('-s', '--source', default='/', 
                        help='Absolute path to the project folder')
    parser.add_argument('-o', '--output', default='project_structure', 
                        help='Path to save the output markdown file')
    parser.add_argument('-id', '--i_directory', default='', 
                        help='Directories to ignore, comma-separated')
    parser.add_argument('-ie', '--i_extension', default='', 
                        help='Extensions to ignore, comma-separated')
    
    args = parser.parse_args()
    
    ignore_dirs = []
    if args.i_directory:
        for d in args.i_directory.split(','):
            ignore_dirs.append(d.strip())
    
    ignore_exts = []
    if args.i_extension:
        for e in args.i_extension.split(','):
            e = e.strip()
            if not e.startswith('.'):
                e = '.' + e
            ignore_exts.append(e)
    
    output_path = args.output
    if not output_path.endswith('.md'):
        output_path += '.md'
    
    generate_directory_structure(
        source_path=args.source, 
        output_path=output_path, 
        ignore_dirs=ignore_dirs, 
        ignore_exts=ignore_exts
    )

if __name__ == '__main__':
    main()