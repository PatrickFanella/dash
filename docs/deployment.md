# Deployment Guide

## Cloudflare Tunnel Configuration

ALMAZ is exposed at `almaz.subcult.tv` via Cloudflare Tunnel.

### Update Tunnel Config

Edit the tunnel configuration file (typically `~/.cloudflared/config.yml` or managed via the Cloudflare dashboard):

```yaml
ingress:
  - hostname: almaz.subcult.tv
    service: http://almaz:8080
  # ... other services
```

### Reload Tunnel

```bash
# If using cloudflared as a systemd service:
sudo systemctl restart cloudflared

# If using docker:
docker restart cloudflared
```

### Rollback

To revert to Dashy, change the service URL back:

```yaml
ingress:
  - hostname: almaz.subcult.tv
    service: http://dashy:8080
```

Then restart the tunnel.

---

## Authelia Protected Application

ALMAZ is protected by Authelia. The following headers are forwarded to the backend:

| Header | Description |
|--------|-------------|
| `Remote-User` | Authenticated username |
| `Remote-Name` | Display name |
| `Remote-Email` | Email address |
| `Remote-Groups` | Comma-separated group memberships |

### Access Control Rule

Add to your Authelia `configuration.yml` under `access_control.rules`:

```yaml
access_control:
  rules:
    - domain: almaz.subcult.tv
      policy: one_factor
      subject:
        - "group:admins"
```

Adjust the policy (`one_factor` or `two_factor`) and subject groups as needed.

### Reverse Proxy Middleware

If using Traefik, ensure the Authelia middleware is applied:

```yaml
# docker-compose labels on the almaz service
labels:
  - "traefik.http.routers.almaz.middlewares=authelia@docker"
```

If using nginx, add the auth_request directives:

```nginx
location / {
    auth_request /authelia;
    auth_request_set $user $upstream_http_remote_user;
    auth_request_set $name $upstream_http_remote_name;
    auth_request_set $email $upstream_http_remote_email;
    auth_request_set $groups $upstream_http_remote_groups;
    proxy_set_header Remote-User $user;
    proxy_set_header Remote-Name $name;
    proxy_set_header Remote-Email $email;
    proxy_set_header Remote-Groups $groups;
    proxy_pass http://almaz:8080;
}
```

### Verify Headers

After deployment, confirm identity headers reach the backend:

```bash
curl -s https://almaz.subcult.tv/api/v1/whoami
# Should return: {"username":"...","display_name":"...","email":"...","groups":[...]}
```
