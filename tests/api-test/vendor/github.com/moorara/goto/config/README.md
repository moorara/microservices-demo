# config

This is a very minimal utility for reading configuration values from either
**command-line flags**, **environment variables**, or **configuration files**.

This library does not use `flag` package for parsing flags, so you can still parse your flags separately.

## Quick Start

```go
package main

import (
  "fmt"
  "github.com/moorara/goto/config"
)

type Spec struct {
  Enabled     bool
  ServicePort int
  LogLevel    string
  DBEndpoints []string
}

func main() {
  spec := Spec{}
  config.Pick(&spec)
  fmt.Printf("%+v\n", spec)
}
```

The precendence of sources for values is as follows:

  1. command-line flags
  2. environment variables
  3. configuration files
  4. default values (set when creating `spec`)

You can pass the configuration values using **flags** using any of the syntaxes below:

```bash
main -enabled -service.port=8080 -log.level=info -db.endpoints=arango1,arango2,arango3
main --enabled --service.port=8080 --log.level=info --db.endpoints=arango1,arango2,arango3
main -enabled -service.port 8080 -log.level info -db.endpoints arango1,arango2,arango3
main --enabled --service.port 8080 --log.level info --db.endpoints arango1,arango2,arango3
```

You can pass the configuration values using **environment variables** as follows:

```bash
export ENABLED=true
export SERVICE_PORT=8080
export LOG_LEVEL=info
export DB_ENDPOINTS=arango1,arango2,arango3
```

You can also write the configuration values in **files**
and set the paths to the files using environment variables:

```bash
export ENABLED_FILE=...
export SERVICE_PORT_FILE=...
export LOG_LEVEL_FILE=...
export DB_ENDPOINTS_FILE=...
```

## Complete Example

```go
package main

import (
  "fmt"
  "github.com/moorara/goto/config"
)

type Spec struct {
  field             string     // Unexported, will be skipped
  FieldString       string    `flag:"fieldString" env:"CONFIG_FIELD_STRING" file:"CONFIG_FILE_FIELD_STRING"`
  FieldBool         bool      `flag:"fieldBool" env:"CONFIG_FIELD_BOOL" file:"CONFIG_FILE_FIELD_BOOL"`
  FieldFloat32      float32   `flag:"fieldFloat32" env:"CONFIG_FIELD_FLOAT32" file:"CONFIG_FILE_FIELD_FLOAT32"`
  FieldFloat64      float64   `flag:"fieldFloat64" env:"CONFIG_FIELD_FLOAT64" file:"CONFIG_FILE_FIELD_FLOAT64"`
  FieldInt          int       `flag:"fieldInt" env:"CONFIG_FIELD_INT" file:"CONFIG_FILE_FIELD_INT"`
  FieldInt8         int8      `flag:"fieldInt8" env:"CONFIG_FIELD_INT8" file:"CONFIG_FILE_FIELD_INT8"`
  FieldInt16        int16     `flag:"fieldInt16" env:"CONFIG_FIELD_INT16" file:"CONFIG_FILE_FIELD_INT16"`
  FieldInt32        int32     `flag:"fieldInt32" env:"CONFIG_FIELD_INT32" file:"CONFIG_FILE_FIELD_INT32"`
  FieldInt64        int64     `flag:"fieldInt64" env:"CONFIG_FIELD_INT64" file:"CONFIG_FILE_FIELD_INT64"`
  FieldUint         uint      `flag:"fieldUint" env:"CONFIG_FIELD_UINT" file:"CONFIG_FILE_FIELD_UINT"`
  FieldUint8        uint8     `flag:"fieldUint8" env:"CONFIG_FIELD_UINT8" file:"CONFIG_FILE_FIELD_UINT8"`
  FieldUint16       uint16    `flag:"fieldUint16" env:"CONFIG_FIELD_UINT16" file:"CONFIG_FILE_FIELD_UINT16"`
  FieldUint32       uint32    `flag:"fieldUint32" env:"CONFIG_FIELD_UINT32" file:"CONFIG_FILE_FIELD_UINT32"`
  FieldUint64       uint64    `flag:"fieldUint64" env:"CONFIG_FIELD_UINT64" file:"CONFIG_FILE_FIELD_UINT64"`
  FieldStringArray  []string  `flag:"field.string.array" env:"FIELD_STRING_ARRAY" file:"FIELD_STRING_ARRAY_FILE" sep:","`
  FieldFloat32Array []float32 `flag:"field.float32.array" env:"FIELD_FLOAT32_ARRAY" file:"FIELD_FLOAT32_ARRAY_FILE" sep:","`
  FieldFloat64Array []float64 `flag:"field.float64.array" env:"FIELD_FLOAT64_ARRAY" file:"FIELD_FLOAT64_ARRAY_FILE" sep:","`
  FieldIntArray     []int     `flag:"field.int.array" env:"FIELD_INT_ARRAY" file:"FIELD_INT_ARRAY_FILE" sep:","`
  FieldInt8Array    []int8    `flag:"field.int8.array" env:"FIELD_INT8_ARRAY" file:"FIELD_INT8_ARRAY_FILE" sep:","`
  FieldInt16Array   []int16   `flag:"field.int16.array" env:"FIELD_INT16_ARRAY" file:"FIELD_INT16_ARRAY_FILE" sep:","`
  FieldInt32Array   []int32   `flag:"field.int32.array" env:"FIELD_INT32_ARRAY" file:"FIELD_INT32_ARRAY_FILE" sep:","`
  FieldInt64Array   []int64   `flag:"field.int64.array" env:"FIELD_INT64_ARRAY" file:"FIELD_INT64_ARRAY_FILE" sep:","`
  FieldUintArray    []uint    `flag:"field.uint.array" env:"FIELD_UINT_ARRAY" file:"FIELD_UINT_ARRAY_FILE" sep:","`
  FieldUint8Array   []uint8   `flag:"field.uint8.array" env:"FIELD_UINT8_ARRAY" file:"FIELD_UINT8_ARRAY_FILE" sep:","`
  FieldUint16Array  []uint16  `flag:"field.uint16.array" env:"FIELD_UINT16_ARRAY" file:"FIELD_UINT16_ARRAY_FILE" sep:","`
  FieldUint32Array  []uint32  `flag:"field.uint32.array" env:"FIELD_UINT32_ARRAY" file:"FIELD_UINT32_ARRAY_FILE" sep:","`
  FieldUint64Array  []uint64  `flag:"field.uint64.array" env:"FIELD_UINT64_ARRAY" file:"FIELD_UINT64_ARRAY_FILE" sep:","`
}

func main() {
  spec := Spec{
    FieldString: "default value",
  }
  config.Pick(&spec)
  fmt.Printf("%+v\n", spec)
}
```
