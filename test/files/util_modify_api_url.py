import json

file_path = 'test/mock_driver/clock/api.json'

with open(file_path, 'r', encoding='utf-8') as file:
    data = json.load(file)

if 'servers' in data and isinstance(data['servers'], list):
    for server in data['servers']:
        if 'url' in server:
            server['url'] = 'http://172.17.0.1'
else:
    print("The 'servers' key is not found or not a list.")

with open(file_path, 'w', encoding='utf-8') as file:
    json.dump(data, file, ensure_ascii=False, indent=4)

print("URL updated successfully.")
