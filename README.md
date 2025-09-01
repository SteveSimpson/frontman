# Frontman Go Web Detecting Proxy

## Overview
This project implements an efficient web proxy server in Go. 
The proxy server is designed to handle incoming HTTP requests, 
forward them to the appropriate target server,
and return the responses back to the clients, detecting possible attacks.

## Development Notes

- Currently everything just logs to stdout - this would need to be configured to go to something like
Datadog, ELK or Prometheus/Grafana.
- Unit Tests are currently done during build, would be better to move them to CI pipeline.
- Still need to implement functional testing.

## Notes on Detections

I have introduced the concept of scoring with the idea that more likely threats have a higher score.
However, this concept will need a lot of adjustment to be meaningful. For now, it is more of a placeholder.

### (Password) Brute Force

Brute forcing logins can be done in many ways:

- many attempts with same username
- many attempts with same password
- many attempts from same source IP

Since this would typically be done in a distributed fasion, a linked cache will have to be used to track the attempts.

There are a couple of strategies to clearing, we could look at the results (succcessful vs unsuccessesful),
but since we are really talking about high volume requests and keeping this as generic as possible, 
probably more effecienmct to ignore the response and only track the attempts.

The cache entries should be time limited so that it cleans itself out.

### Hostile Sources

Admittedly, this is not the best implementation. I have not added a script to add the source IPs, yet. 

### SQL Injection

This should find serveral attempted SQL injection attacks, but I think the scoring will need a lot of tuning.

I'm also not sure that this is the best implememntation.

## Detectors Future Work

### XSS detector

A simple version could be implemented using `sql_injection.go` as the base and look for things like:

```
<script>, <img src=… onerror=…>, <svg onload=…>
onmouseover=, onclick=
%3Cscript%3Ealert(1)%3C/script%3E
```

### OS / Path Interaction

A simple version could be implemented using `sql_injection.go` as the base and look for things like:

```
- ;, &&, ||, | (chaining operators)
- $(command) or backticks `command`
- cat /etc/passwd
- whoami
- ls
- ping
- curl
- wget
- ../
- ..%2F
- ..\\
- /etc/passwd
- boot.ini
- id_rsa
- page=../../etc/passwd
- page=http://evil.com/shell.txt
```

Not sure how many false positives you would get, so a better stragegy might be to acutally parse the input, 
but that would have the trade off of getting into heavier processing/

## Setup Instructions
1. Clone the repository:
   ```
   git clone <repository-url>
   cd frontman
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Set environment variables as needed in dev.env.

## Usage
To run the web proxy server and dependencies, execute the following commands:
```
docker compose build api
docker compose up -d
```

The default dev setup uses wordpress as the backend web server and sets up a base infrastructure.

Pieces can be accessed directly:

- WordPress via proxy: http://localhost:8888
- Dozzle Log Viewer: http://localhost:8080
- PHP My Admin: http://localhost:8036 (to access WP's database if needed)

To shut it down:
```
docker compose down
```

## Features
- Handles incoming HTTP requests and forwards them to target server.
- Supports configuration through environment variables.
- Efficient request handling and response forwarding.
- Detects some possible web based attacks.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.
