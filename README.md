# InfluxDB Line-Protocol Codec

A high-performance Go codec for the [InfluxDB line protocol](https://docs.influxdata.com/influxdb/latest/reference/syntax/line-protocol/) syntax.

The API is intentionally low level -- it's designed for converting line protocol to concrete types without imposing a specific `Point` type. This makes it suitable for high-throughput ingestion pipelines where allocation counts matter.

API documentation: https://pkg.go.dev/github.com/influxdata/line-protocol/v2/lineprotocol

## Usage

### Decoding

```go
data := []byte(`foo,tag1=val1,tag2=val2 x=1,y="hello" 1625823259000000
bar enabled=true
`)
dec := lineprotocol.NewDecoderWithBytes(data)
for dec.Next() {
    m, err := dec.Measurement()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("measurement %s\n", m)
    for {
        key, val, err := dec.NextTag()
        if err != nil {
            log.Fatal(err)
        }
        if key == nil {
            break
        }
        fmt.Printf("  tag %s=%s\n", key, val)
    }
    for {
        key, val, err := dec.NextField()
        if err != nil {
            log.Fatal(err)
        }
        if key == nil {
            break
        }
        fmt.Printf("  field %s=%v\n", key, val)
    }
    t, err := dec.Time(lineprotocol.Microsecond, time.Time{})
    if err != nil {
        log.Fatal(err)
    }
    if !t.IsZero() {
        fmt.Printf("  timestamp %s\n", t.UTC().Format(time.RFC3339Nano))
    }
}
```

### Encoding

```go
var enc lineprotocol.Encoder
enc.SetPrecision(lineprotocol.Microsecond)
enc.StartLine("foo")
enc.AddTag("tag1", "val1")
enc.AddTag("tag2", "val2")
enc.AddField("x", lineprotocol.MustNewValue(1.0))
enc.AddField("y", lineprotocol.MustNewValue("hello"))
enc.EndLine(time.Unix(0, 1625823259000000000))
if err := enc.Err(); err != nil {
    log.Fatal(err)
}
fmt.Printf("%s", enc.Bytes())
// Output: foo,tag1=val1,tag2=val2 x=1,y="hello" 1625823259000000
```

## Tools

- `cmd/lpverify` -- CLI tool that reads line protocol from stdin and reports errors.
- `cmd/verify-lines` -- CGO shared library for verifying line protocol from other languages (e.g. Python via ctypes).
