# Closi

`config.yml`:
```yaml
app:
  name: "Closi"

log:
  level: "debug"
  format: "console"

http:
  host: "127.0.0.1"
  port: 8080
  read_timeout: 10s
  write_timeout: 30s
  idle_timeout: 60s

mongo:
  uri: ""
  database: ""

auth:
  password:
    salt: ""

  tokens:
    access_token:
      signing_key: ""
      ttl: 15m
    refresh_token:
      length: 32
      ttl: 720h
```