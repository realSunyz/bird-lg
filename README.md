# BIRD Looking Glass

A web-based looking glass for BIRD routers, backed by a central server and remote probe clients.
Run Ping and Traceroute from multiple POPs, and optionally query BIRD via SSO.

## Features

- Central server that hosts the UI and proxies requests to probe nodes
- Remote client nodes that run `ping`, `traceroute`, and BIRD commands
- Streaming output for Ping and Traceroute via SSE
- (Optional) Cloudflare Turnstile CAPTCHA gate for tool access
- (Optional) Logto OIDC SSO for BIRD queries
- (Optional) Trace IP enrichment via IPInfo Lite

## Quick Start

Server:

```bash
docker pull ghcr.io/realsunyz/bird-lg:frontend
```

Client:

```bash
docker pull ghcr.io/realsunyz/bird-lg:client
```

## Configuration

The server reads `config.yaml` by default, or `CONFIG_FILE` if set. Environment variables override the file-based config.

| Variable               | Default                         | Description                                     |
| ---------------------- | ------------------------------- | ----------------------------------------------- |
| `CONFIG_FILE`          | `config.yaml`                   | Config file path                                |
| `LISTEN_ADDR`          | `:3000`                         | Overrides `listen`                              |
| `STATIC_DIR`           | `./static`                      | Overrides `static_dir`                          |
| `TURNSTILE_SITE_KEY`   | -                               | Turnstile site key                              |
| `TURNSTILE_SECRET_KEY` | -                               | Turnstile secret key                            |
| `JWT_SECRET`           | auto-generated if not specified | JWT signing secret                              |
| `HMAC_SECRET`          | -                               | Shared secret for server-to-client requests     |
| `LOGTO_ENDPOINT`       | - (optional)                    | Logto OIDC endpoint                             |
| `LOGTO_APP_ID`         | - (optional)                    | Logto application ID                            |
| `SERVERS`              | use `config.yaml`               | Full server list (YAML format)                  |
| `IPINFO_TOKEN`         | - (optional)                    | IPInfo DB auto-download for trace IP enrichment |

## Contributing

Issues and Pull Requests are definitely welcome!

Please make sure you have tested your code locally before submitting a PR.

## License

Source code is released under the [MIT License](https://github.com/realsunyz/bird-lg/blob/main/LICENSE).
