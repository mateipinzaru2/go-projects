import json
import sys
import os

def find_matching_text_values(d, key='text', substring='{"Success":true,"Message":null,"MessageTitle":null,"Data":{"StartIndex"'):
    """Recursive function to find all values for the "text" key that start with '{"Success":true,"Message":null,"MessageTitle":null,"Data":{"StartIndex"'."""
    if isinstance(d, dict):
        for k, v in d.items():
            if k == key and isinstance(v, str) and v.startswith(substring):
                yield v
            elif isinstance(v, (dict, list)):
                for result in find_matching_text_values(v, key, substring):
                    yield result
    elif isinstance(d, list):
        for item in d:
            for result in find_matching_text_values(item, key, substring):
                yield result

def extract_text_values(file_path, output_dir):
    # Create the output directory if it doesn't exist
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    # Load the JSON data from the file
    with open(file_path, 'r') as file:
        data = json.load(file)

    # Extract the matching text values
    matching_text_values = list(find_matching_text_values(data))

    # Define the output file path
    output_file_path = os.path.join(output_dir, 'extracted_text_values.txt')

    # Save the extracted values to the output file
    with open(output_file_path, 'w') as file:
        for item in matching_text_values:
            file.write(item + '\n')

    return output_file_path

if __name__ == '__main__':
    if len(sys.argv) != 3:
        print("Usage: python3 main.py <input_json_file> <output_directory>")
        sys.exit(1)

    # Command-line arguments
    input_json_file = sys.argv[1]
    output_directory = sys.argv[2]

    # Run the function with the provided arguments
    output_file_path = extract_text_values(input_json_file, output_directory)
    print(f"Extracted values are saved to: {output_file_path}")