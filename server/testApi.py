#!/bin/python3

import requests

with open("./pub_key") as f:
    data = f.read()
    requests.post("http://localhost:8000/api/pubkey", data=data)