# Vote Application
This application is part of the "ThemaSemester" at the Hogeschool van Amsterdam.
The goal is to implement different cryptographic functions to create a secure vote app.

**A number of assumptions have been made :**
1. The users have pre-voting process received an invitation with a unique code (in the application referenced as userCode)
2. the user has a secure environment where the CLI will be used (no spyware etc).
3. The server IP can't be spoofed; right now, no certificate has been implemented nor is HTTPS enabled

## Voting Process explained (Server Side)
### Setup
1. The server prompts user for a subject of the vote (right now this a yes/no or for or against vote).
2. The server checks for a private key file, if one doesn't exist, a new RSA-Key pair will be created and written to file.
3. A basic database test will be carried out (will write to database and read that data from the database) any errors will be reported.
4. If any error returned, the program will report to stderr and exit,

### Router
1. After the setup process, a router (Gorilla-Mux) is created and four routes are created:
*```$URL & $PORT``` are hardcoded, a flag option could be added, but isn't at this time.*
- ```$URL/```               **-> Index** 
- ```$URL/api/pubkey```     **-> retrieveKey** 
- ```$URL/api/getvote```    **-> getVote** 
- ```$URL/vote/cast```      **-> handleVote** 

### Routes
*Params are passed by the CLI*

#### Index
Returns hello, World. For testing the connection, not used in practice.
##### Methods
No restrictions
##### Params 
None

#### retrieveKey
This route is used to exchange public keys.
Keys will be stored in a Database (Key: userCode, Value: Public Key (*byte array*))
##### Methods
* POST
##### Params
```json
{
  "usercode": <usercode>,
  "publickey":  <publickey>
}
```
###### Potential Security Flaw
Right now any key will be overwritten.
This is an **major** security flaw, however this should be somewhat negated by a long / unique enough usercode.

#### getVote
This route is used to send the vote subject to the CLI-user.
*The vote subject was created during the setup phase.*
##### Methods
* POST
##### Params
* usercode (encrypted with server-public-key)
###### Potential Security Flaw
Right now no signing has been implemented, I plan to do so, but this might not be ready before the deadline

#### handleVote
The bread and butter of the application. Here a vote will be cast, and written to the the database.
##### Methods
* POST
##### Params
An encrypted message (enc with server-public-key):
```json
    "usercode": <usercode>
    "vote_val": <vote_val>
    "hash": <hash>
    "sign": <sign>
```
Any recast of vote will update the vote
##### Database Param
The votes are written to the database using the usercode as key.
After the vote this database should be made public, so one can validate his/her own vote.



## Voting Process explained (client-side (CLI))
### Setup

1. The server checks for a private key file, if one doesn't exist, a new RSA-Key pair will be created and written to file.
   **Good to know, no SSH-key's will be overwritten**
### Vote Process
1. The vote process is retrieved from the server (```$URL/api/getvote```).
2. The user is prompted to vote for or against the subject
3. Sign and hash the vote.
4. Encrypt using server-public-key
5. Print response to the CLI and exit.

## Voting Process explained (server-side)
### Setup
1. First prompts the user what the subject of the vote is
2. A Key-pair is generated (if a keypair doesn't exist) **Good to know, no SSH-key's will be overwritten**
3. The database is checked by writing to and reading from the Redis database. If no connection could be made, the program exits.
4. A router object is created.
### Vote process.
1. A CLI-user will automaticly upon starting send their public key to the server using the URI ```<URL>/api/pubkey```
2. The server respondes with sending his own public key.
3. The key send by the user is written to the database (key is the usercode).
4. the user will then request the vote subject, this is requested by calling the URI ```<URL>/api/getvote```
The response will be encrypted and signed by the server.
5. The user will response if he/she is for or against the votesubject.
6. The data will be stored with the key being the usercode.
7. The server will response with a message that the user has voted.

## Possible weaknesses.
1. The usercode needs to be unique enough otherwise it would be possible to bruteforce. 
2. Right now, no usercheck is carried out. This should be done in a real life environment.
3. There is no 2FA.
4. The server can't be verified.
5. There is no way to check if a user has voted in the past.