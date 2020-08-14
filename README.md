# timetype ![Go](https://github.com/Semior001/timetype/workflows/Go/badge.svg) [![Coverage Status](https://coveralls.io/repos/github/Semior001/timetype/badge.svg?branch=master)](https://coveralls.io/github/Semior001/timetype?branch=master) [![go report card](https://goreportcard.com/badge/github.com/semior001/timetype)](https://goreportcard.com/report/github.com/semior001/timetype)
Package adds some time types for easier work, serialize and deserialize them in requests.

```go
// Duration is a wrapper of time.Duration, that allows to marshal and unmarshal time in RFC3339 format
type Duration time.Duration
``` 

```go
// Clock is a wrapper for time.time to allow parsing datetime stamp with time only in
// ISO 8601 format, like "15:04:05"
type Clock time.Time
```
