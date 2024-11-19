# Closi

`config.yml`:
```yaml
app:
  name: "Closi"

localizer:
  files:
    - ""
    - ""

  default_language: ""

log:
  level: "debug"
  format: "console" # json

http:
  host: "127.0.0.1"
  port: 8080
  read_timeout: 10s
  write_timeout: 30s
  idle_timeout: 60s

mongo:
  uri: ""
  database: ""

redis:
  host: ""
  port: 0
  password: ""
  database: 0

auth:
  confirmation_link_format: ""

  password:
    salt: ""

  tokens:
    access_token:
      signing_key: ""
      ttl: 15m
    refresh_token:
      length: 32
      ttl: 5s # 720h

imgbb:
  api_key: ""
  timeout: 60s

smtp:
  host: "smtp.gmail.com"
  port: 465
  username: ""
  password: ""

```