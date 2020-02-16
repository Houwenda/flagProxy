# flagProxy

A tool to detect cheating in CTF competitions based on TCP proxy. (Just a proof of concept.)

# concept

FlagProxy client works as a proxy server between users and the CTF challenge, inspecting responses from the challenge server. 
Users are recognized by the ports they are visiting. 
When the initial flag is found in a response, flagProxy client requests the flagProxy server for real flag of the participant. 
Then the client replaces the initial flag with real flag.

# assumptions

- All payloads and flags should be transmitted through TCP connections. 
- Flags should not be placed in more than TWO segments. (Test is needed if encoding is used in traffic.)
- Flags can be matched using regular expressions. (In other words, encrypted traffic is currently not supported by flagProxy.)

# client config

An example of flagProxy client.

```yaml
challenge:
  address: "192.168.0.108:8000"  # address of the challenge
  flag_regex: "flag{.+}"  #  the regular expression to detect initial flag
  threads: 10  # controls the throughput of the challenge
server:
  url: "http://localhost:8080"  # server api address
  challenge_id : "testchallengeid"  # a random string for api server to recognize challenge
  key: "testkeyforchallenge0"  # a random key to access real flags
log:
  path: "/tmp/flagProxy.log"
```

# server api

## database

FlagProxy server maintains tables which saves the flag & other data of challenges. 

Example:

challenge table: 

| challenge_id     |     challenge_key |
| :--------------- | ----------------: |
| testchallengeid0 | testchallengekey0 |
| testchallengeid1 | testchallengekey1 |
| testchallengeidx | testchallengekeyx |

flag table: 

| challenge_id | user_id | port | flag |
| :----------- | ------- | --------- | ---: |
| testchallengeid0 | user0 | 10000 | testflag0 |
| testchallengeid0 | user1 | 10001 | testflag1 |

## RESTful API

Here are the RESTful APIs the flagProxy server should implement.

### 1. connection test

Test the connection between client and server.

Request: `GET /`

Response: `flagProxy`

### 2. proxy ports

Get the ports to listen.

Request: `GET /ports/{challengeId}/{key}`

Response example: 
```json
{
    "msg": "success",
    "ports": [
        10001,
        10002,
        10003,
        10004
    ]
}
```
or 
```json
{
    "msg": "auth error",
    "ports": []
}
```

Note: Content-Type of the response should be "application/json".

### 3. real flag

Request: `/flagByPort/{challengeId}/{key}/{port}`

Response example:
```json
{
    "msg": "success",
    "flag": "flag{test_flag_for_flag_proxy}"
}
```
or 
```json
{
    "msg": "auth error",
    "flag": ""
}
```

Note: Content-Type of the response should be "application/json".

# usage

- First, Implement the server APIs.
- Second, compile the flagProxy client.
- Third, configure the client by editing the config file. 
- Fourth, run the compiled client after flagProxy server is brought up by executing `./compiledClient /path/to/config/file`

# warning

This tool is designed to be compatible with both web and binary challenges. But as a proof of concept, it has not been fully tested in security and functions, thus should not be used in production environment.
