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
  -b, --batch int         the batch size when using the scan command (default 10)
  -c, --cardinality int   the number of elements of a key, used to filter bigkey (default 0)
  -n, --db int            redis database (default 0)
  -f, --format string     output format (oneof: csv, json, xml) (default "csv")
  -h, --help              help for redis-doctor
      --host string       redis server host (default "127.0.0.1")
  -l, --length int        serialized length of a key, used to filter bigkey (default 0)
      --limit int         the number of returned entries (default 10)
      --pass string       redis password
      --pattern string    keys pattern when using the --bigkeys or --hotkey options (default "*")
  -p, --port int          redis server port (default 6379)
  -s, --symptom string    symptom to diagnose (required, oneof: bigkey, hotkey, slowlog)
  -t, --type string       redis data type (oneof: string, list, hash, set, zset)
  -u, --user string       redis username
```