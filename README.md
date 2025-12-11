# Tunnel

A compact Go utility that establishes a reverse SSH tunnel. It listens on a remote TCP address exposed by your SSH host and forwards traffic to a service running on your local machine.

## Features
- Sets up a password-authenticated SSH session and reverse tunnel in a single binary.
- Pulls configuration from environment variables (optional `.env` file for local development).
- Logs tunnel activity and propagates connection failures to aid debugging.

## Prerequisites
- Go 1.21 or newer
- Network reachability to the SSH bastion/server you plan to use

## Configuration
Set the following variables as environment variables or inside a `.env` file (loaded via [`github.com/joho/godotenv`](https://github.com/joho/godotenv)):

| Variable | Description | Default |
| --- | --- | --- |
| `SSH_USER` | SSH username used to authenticate. | `root` |
| `SSH_PASSWORD` | SSH password. | `123456` |
| `SSH_SERVER_ADDR` | SSH server hostname and port. | `example.com:22` |
| `REMOTE_ADDR` | Remote bind address exposed by the SSH server. | `0.0.0.0:9001` |
| `LOCAL_ADDR` | Local service address to forward traffic to. | `127.0.0.1:9091` |

> ⚠️ Update the defaults before exposing this tunnel in production. The sample password and addresses are placeholders and **not secure**.

Example `.env`:

```env
SSH_USER=alice
SSH_PASSWORD=super-secret
SSH_SERVER_ADDR=ssh.example.com:22
REMOTE_ADDR=0.0.0.0:443
LOCAL_ADDR=127.0.0.1:8080
```

## Usage

```powershell
# Install dependencies (once)
go mod tidy

# Run the tunnel
$env:SSH_USER="alice"; $env:SSH_PASSWORD="super-secret"; go run main.go
```

To run with a `.env` file, place it next to `main.go` and execute `go run main.go`. The loader only affects a local development environment; production deployments should rely on your process manager or secrets store to provide environment variables.

## Troubleshooting
- **Connection refused**: confirm the SSH server is reachable and allows TCP forwarding/reverse tunnels.
- **Authentication failed**: verify the credentials or switch to key-based auth before deploying in a real environment.
- **Bind errors** (`address already in use`): adjust `REMOTE_ADDR` to a port not already claimed by the SSH host.

## Next Steps
- Replace password auth with SSH keys for better security.
- Add health checks or metrics if you will monitor the tunnel in production.
