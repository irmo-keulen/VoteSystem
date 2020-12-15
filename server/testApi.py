#!/bin/python3

import requests
import json

# with open("./pub_key") as f:
#     data = f.read()
#     requests.post("http://localhost:8000/api/pubkey", data=data)

pkey = """-----BEGIN PUBLIC KEY----- 
MIICCgKCAgEAwrb8SKaS6qonRGs6YhI3VE0EHR09/GOdFfY+D19a8YUYSm8/ygOT
EqOBhAh+g8NYwS/y55Td+moItMoUhEWulsiyPNFji6xiY2mZ9/7cZxKW1RsvWB7s
nNyz2LeH7/L3Jo/P77wr6NMZjBJKEdu1K9/Dj4/rucJWMz90wSuZgq9lbs6/KkwI
IplB7y/2SyOnYyV0ytq0zlLh0nrctge8UbU4qqo2nQtyr4lUAblm1b9qGiH1YzI3
M+3uHf7vA+5j5RfHHIC9Z5QiXzPw90m1Dks0pRoWEfQZUaT2yCPxpGpfoax9vdrO
I7LLuSBj0SoZAeFq+dCImJ3t9YS9WJ0DlRU6SEg1dnNNq2w547AENqumTKheIAhW
7njjz2MSwaRajJB70muVlSkyh1F5P6BPkzeHp/O6fXVzaMh9L768CEcWjrqWphpI
YsGtt8ojRDx2q3vx5SNtBa05cq4Zj5ZV/6yBA1WbDjEdehQSJX+wJj4w1Azfyhhg
6VN5KI1BzEtewMCeMcfk9p14+UzEGDoJSUXBXQW7gVmq8vRKpTMsYGkDYtTsVkX9
aMP8BCq7C4/6Ifujo+QsCSiiooLyqSiMNvbNVlb4+4H5zEwAiDcYJKqNJH5nKB6H
J7z627h5kky7jNtGVaLIob74xLEgOS6pa7gMCanGtC0Fb4y5DZTA+1sCAwEAAQ==
-----END PUBLIC KEY-----"""

string = {
    "usercode": "VeryUniqueUserCode", 
    "publickey":pkey
    }
js = json.dumps(string)

requests.post("http://localhost:8000/api/pubkey", data=js)
