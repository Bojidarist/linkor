# Linkor

A self-hosted link shortener with an admin panel for managing short URLs and tracking click statistics.

## Features

- Create short links with custom or auto-generated slugs
- Track total clicks and unique clicks per link
- Admin panel with a modern dark UI
- API and admin panel protected by secret key
- `.env` file support for configuration
- Single binary deployment (all assets embedded)
- SQLite database (no external dependencies)
- Pure Go — no CGO required

## Requirements

- Go 1.22 or later

## Build

```bash
go build -o linkor .
```

This produces a single `linkor` binary with all HTML/CSS/JS assets embedded.

### Cross-compile examples

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o build/linkor .

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -trimpath -o build/linkor .

# Windows
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o build/linkor.exe .
```

## Configuration

Linkor is configured via environment variables. You can set them directly or use
a `.env` file in the working directory.

Copy the example file to get started:

```bash
cp .env.example .env
```

Then edit `.env` with your values.

| Variable           | Required | Default      | Description                        |
| ------------------ | -------- | ------------ | ---------------------------------- |
| `ADMIN_SECRET_KEY` | Yes      | —            | Secret key to access the admin panel and API |
| `PORT`             | No       | `8080`       | HTTP server port                   |
| `DATABASE_PATH`    | No       | `linkor.db`  | Path to the SQLite database file   |

The `.env` file is optional. Environment variables set in the shell take
precedence over values in `.env`.

## Run

```bash
# Using a .env file (recommended)
cp .env.example .env
# Edit .env with your values
./linkor
```

Or with inline environment variables:

```bash
ADMIN_SECRET_KEY=mysecret ./linkor
```

The server starts on the configured port (default `8080`).

## Usage

Open the admin panel in your browser:

```
http://localhost:8080/admin/management?key=your-secret-key-here
```

From there you can create, edit, and delete short links and monitor click statistics.

For detailed usage instructions and API reference, see [docs/usage.md](docs/usage.md).

## License

See [LICENSE](LICENSE) for details.
