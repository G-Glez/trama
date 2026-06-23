import base64, json, os, urllib.request

url = os.environ["ACTIONS_ID_TOKEN_REQUEST_URL"] + "&audience=sts.amazonaws.com"
token = os.environ["ACTIONS_ID_TOKEN_REQUEST_TOKEN"]
req = urllib.request.Request(url, headers={"Authorization": f"bearer {token}"})
resp = json.loads(urllib.request.urlopen(req).read())
payload = resp["value"].split(".")[1]
pad = 4 - len(payload) % 4 if len(payload) % 4 else 0
payload += "=" * pad
decoded = json.loads(base64.urlsafe_b64decode(payload))
print(json.dumps(decoded, indent=2))
