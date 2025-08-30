# Frontman Go Web Detecting Proxy

## Overview
This project implements an efficient web proxy server in Go. The proxy server is designed to handle incoming HTTP requests, forward them to the appropriate target server, and return the responses back to the clients, detecting possible attacks.

## Setup Instructions
1. Clone the repository:
   ```
   git clone <repository-url>
   cd go-web-proxy
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Configure the proxy settings in `internal/config/config.go` or set environment variables as needed.

## Usage
To run the web proxy server, execute the following command:
```
go run cmd/main.go
```

## Features
- Handles incoming HTTP requests and forwards them to target server.
- Supports configuration through environment variables.
- Efficient request handling and response forwarding.
- Detects some possible web based attacks.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.