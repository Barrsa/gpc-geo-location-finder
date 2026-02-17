# GPC Geo Location Finder

A Go-based web service that checks ping latency to multiple Google Cloud Platform (GCP) endpoints across different regions. Built with [Gin](https://gin-gonic.com/) web framework, this service provides real-time latency measurements to help identify the closest GCP region.

## Features

- 🌍 **Multi-Region Support**: Checks latency to endpoints across multiple GCP regions worldwide
- ⚡ **Concurrent Pinging**: Performs all ping checks concurrently for fast results
- 📊 **Detailed Summary**: Returns comprehensive statistics including fastest/slowest regions
- 🔧 **Configurable**: Endpoint configuration via JSON file and environment variables
- 🐳 **Docker Support**: Ready-to-use Docker image with health checks
- 🚀 **Production Ready**: Graceful shutdown, timeout handling, and error management

## Requirements

- Go 1.23+ (or Go 1.21+ for Docker builds)
- `endpoints.json` file with endpoint configurations
- Network access to GCP endpoints

## Installation

### Local Development

1. Clone the repository:
```bash
git clone <repository-url>
cd gpc-geo-location-finder
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file (or use environment variables):
```bash
cp .env.example .env
```

4. Edit `.env` and set the path to your endpoints file:
```env
ENDPOINTS_FILE_PATH=/path/to/endpoints.json
PORT=8080
```

5. Build the application:
```bash
go build -o gcping cmd/main.go
```

6. Run the application:
```bash
./gcping
```

Or run directly:
```bash
go run cmd/main.go
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENDPOINTS_FILE_PATH` | Path to the JSON file containing endpoint configurations | **Required** |
| `PORT` | Port number for the HTTP server | `8080` |

### Endpoints JSON Format

The `endpoints.json` file should contain a JSON object mapping region keys to endpoint configurations:

```json
{
  "us-central1": {
    "URL": "https://us-central1-5tkroniexa-uc.a.run.app",
    "Region": "us-central1",
    "RegionName": "Iowa"
  },
  "europe-west1": {
    "URL": "https://europe-west1-5tkroniexa-ew.a.run.app",
    "Region": "europe-west1",
    "RegionName": "Belgium"
  }
}
```

### Command Line Flags

- `-port`: Override the port number (takes precedence over `PORT` env var)
- `-timeout`: Set timeout for ping requests (default: 30s)

Example:
```bash
./gcping -port 9000 -timeout 60s
```

## API Endpoints

### Health Check

**GET** `/health`

Returns the health status of the service.

**Response:**
```json
{
  "status": "ok"
}
```

### Check Ping

**GET** `/api/checkPing`

Checks latency to all configured endpoints and returns a summary.

**Response:**
```json
{
  "timestamp": "2026-02-17T23:14:05Z",
  "totalRegions": 30,
  "successful": 28,
  "failed": 2,
  "fastest": {
    "region": "us-central1",
    "regionName": "Iowa",
    "url": "https://us-central1-5tkroniexa-uc.a.run.app",
    "latencyMs": 45,
    "success": true
  },
  "slowest": {
    "region": "asia-southeast2",
    "regionName": "Jakarta",
    "url": "https://asia-southeast2-5tkroniexa-et.a.run.app",
    "latencyMs": 234,
    "success": true
  },
  "results": [
    {
      "region": "us-central1",
      "regionName": "Iowa",
      "url": "https://us-central1-5tkroniexa-uc.a.run.app",
      "latencyMs": 45,
      "success": true
    },
    ...
  ]
}
```

**Response Fields:**
- `timestamp`: ISO 8601 timestamp of when the check was performed
- `totalRegions`: Total number of regions checked
- `successful`: Number of successful ping checks
- `failed`: Number of failed ping checks
- `fastest`: Details of the fastest responding endpoint (if any)
- `slowest`: Details of the slowest responding endpoint (if any)
- `results`: Array of all ping results

**CORS:** The API includes CORS headers allowing cross-origin requests.

## Docker

### Build Docker Image

```bash
docker build -t gpc-geo-location-finder .
```

### Run Docker Container

```bash
docker run -d \
  -p 8080:8080 \
  -e ENDPOINTS_FILE_PATH=/app/endpoints.json \
  -v /path/to/endpoints.json:/app/endpoints.json \
  gpc-geo-location-finder
```

Or using docker-compose:

```yaml
version: '3.8'
services:
  gpc-geo-location-finder:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENDPOINTS_FILE_PATH=/app/endpoints.json
    volumes:
      - ./endpoints.json:/app/endpoints.json
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
```

## Project Structure

```
gpc-geo-location-finder/
├── cmd/
│   └── main.go              # Application entry point
├── handlers/
│   └── HandleCheckPing.go   # HTTP handlers
├── models/
│   └── AllEndpoints.go       # Endpoint model and loading logic
├── internal/
│   └── CheckPing.go          # Core ping logic
├── endpoints.json            # Endpoint configurations
├── .env                      # Environment variables (not in git)
├── .env.example             # Example environment file
├── Dockerfile                # Docker build configuration
├── go.mod                    # Go module dependencies
├── go.sum                    # Go module checksums
└── README.md                 # This file
```

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
CGO_ENABLED=0 GOOS=linux go build -o gcping cmd/main.go
```

## License

Copyright 2026 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
