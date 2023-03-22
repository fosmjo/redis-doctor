# redis-doctor

`redis-doctor` is a cli tool for diagnosing redis problems, such as hotkey, bigkey, slowlog, etc.

# Install

```sh
go install github.com/fosmjo/redis-doctor/cmd/redis-doctor@latest
```

# Usage

```sh
$ redis-doctor -h
redis-doctor is a cli tool for diagnosing redis problems, such as hotkey, bigkey, slowlog, etc.

Usage:
  redis-doctor [flags]

Flags:
  -c, --count int        specify the number of returned entries (default 10)
  -n, --db int           redis database (default 0)
  -f, --format string    output format (oneof: csv, json) (default "csv")
  -h, --help             help for redis-doctor
      --host string      redis server host (default "127.0.0.1")
      --pass string      redis password
      --pattern string   keys pattern when using the --bigkeys or --hotkey options (default "*")
  -p, --port int         redis server port (default 6379)
  -s, --symptom string   symptom to diagnose (required, oneof: bigkey, hotkey, slowlog)
  -u, --user string      redis username
```