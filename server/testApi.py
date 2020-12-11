#!/bin/python3

import requests
import json

# with open("./pub_key") as f:
#     data = f.read()
#     requests.post("http://localhost:8000/api/pubkey", data=data)

string = {"usercode": "userco", "username":"user", "voted":False, "publicKey":""}
js = json.dumps(string)

requests.post("http://localhost:8000/api/pubkey", data=js)