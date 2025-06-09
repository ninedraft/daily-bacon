# daily-bacon

daily-bacon is a CLI tool to fetch air quality and environmental data from various sources and post formatted updates to Telegram groups.

## Table of Contents

- [daily-bacon](#daily-bacon)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Prerequisites](#prerequisites)
  - [Installation / Build](#installation--build)
  - [Configuration](#configuration)
  - [Usage](#usage)
  - [Development](#development)
  - [Project Structure](#project-structure)
  - [License](#license)

## Features

- Fetches latest air quality metrics (PM2.5, PM10, dust, pollen, ozone, etc.) and environmental data via Open-Meteo.
- Formats data into a Telegram-friendly message.
- Posts updates to one or multiple Telegram groups concurrently.
- Configurable location (latitude/longitude) and time parameters.
- Lightweight Go-based CLI.

## Prerequisites

- Go 1.20+ installed  
- A Telegram bot token saved to a file.

## Installation / Build

```bash
git clone https://github.com/ninedraft/daily-bacon.git
cd daily-bacon
go build -o daily-bacon ./cmd/daily-bacon
```

Or install via Go toolchain:

```bash
go install github.com/ninedraft/daily-bacon/cmd/daily-bacon@latest
```

## Configuration

Set the `TELEGRAM_TOKEN_FILE` environment variable pointing to your bot token file:

```bash
export TELEGRAM_TOKEN_FILE=/path/to/token.txt
```

## Usage

```bash
./daily-bacon [flags]
```

**Flags:**

- `--group-id` Telegram group ID to post to (can be set multiple times or separated by comma, space, or '|').
- `--latitude` Location latitude (default: 34.707130).
- `--longitude` Location longitude (default: 33.022617).
- `--timeout` Request timeout (default: 10s).

**Examples:**

```bash
./daily-bacon --group-id 123456789 --group-id 987654321
```

Or with comma-separated IDs:

```bash
./daily-bacon --group-id "123456789,987654321" --latitude 40.7128 --longitude -74.0060
```

## Development

1. Clone the repository.  
2. Run tests:
   ```bash
   go test ./...
   ```
3. Run linter:
   ```bash
   golangci-lint run --config .golangci.yml
   ```
4. Format code:
   ```bash
   go fmt ./...
   ```

## Project Structure

```text
cmd/
  daily-bacon/   CLI application entrypoint -> [`cmd/daily-bacon/main.go`](cmd/daily-bacon/main.go:1)
  openmeteo/     Open-Meteo API client       -> [`cmd/openmeteo/openmeteo.go`](cmd/openmeteo/openmeteo.go:1)
internal/
  client/        HTTP client wrapper         -> [`internal/client/client.go`](internal/client/client.go:1)
  meteo/         Data fetchers and types      -> [`internal/meteo/meteo.go`](internal/meteo/meteo.go:1)
  tg/            Telegram messaging client    -> [`internal/tg/tg.go`](internal/tg/tg.go:1)
  view/          Message formatter           -> [`internal/view/view.go`](internal/view/view.go:1)
  models/        Shared data models          -> [`internal/models/airquality.go`](internal/models/airquality.go:1)
```

## License

This project is licensed under the Apache License V2. See the [`LICENSE`](LICENSE) file for details.