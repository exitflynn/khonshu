#!/bin/bash
mkdir -p test_project/src/utils
touch test_project/README.md
touch test_project/src/main.py
touch test_project/src/utils/helper.py
touch test_project/.gitignore
echo "*.log" > test_project/.gitignore
touch test_project/debug.log

python main.py -s test_project -o test_structure.md