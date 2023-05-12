import requests

DALLE_URL = "https://api.openai.com/v1/images/generations"

headers = {"Content-Type": "application/json", "Authorization": "Bearer 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"}
request_params = {"prompt": "", "n": 1, "size": "1024x1024"}

result = requests.get(DALLE_URL, json = request_params, headers = headers)
print(result.json())
