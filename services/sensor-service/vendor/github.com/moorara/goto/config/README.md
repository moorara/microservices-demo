# config

This is a very minimal utility for reading configuration values from either
**command-line flags**, **environment variables**, or **configuration files**.

This library does not use `flag` package for parsing flags

## Quick Start

```go
package main

import (
  "fmt"
  "github.com/moorara/goto/config"
)

type Spec struct {
  Enabled bool
  ServicePort int
  DatabaseURL string
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
main -enabled -service.port=8080 -database.url=root@localhost
main --enabled --service.port=8080 --database.url=root@localhost
main -enabled -service.port 8080 -database.url root@localhost
main --enabled --service.port 8080 --database.url root@localhost
```

You can pass the configuration values using **environment variables** as follows:

```bash
export ENABLED=true
export SERVICE_PORT=8080
export DATABASE_URL=root@localhost
```

You can also write the configuration values in **files**
and set the paths to the files using environment variables:

```bash
export ENABLED_FILE=...
export SERVICE_PORT_FILE=...
export DATABASE_URL_FILE=...
```

## Complete Example

```go
package main

import (
  "fmt"
  "github.com/moorara/goto/config"
)

type Spec struct {
	field        string   // Unexported, will be skipped
	FieldString  string   `flag:"fieldString" env:"CONFIG_FIELD_STRING" file:"CONFIG_FILE_FIELD_STRING"`
	FieldBool    bool     `flag:"fieldBool" env:"CONFIG_FIELD_BOOL" file:"CONFIG_FILE_FIELD_BOOL"`
	FieldFloat32 float32  `flag:"fieldFloat32" env:"CONFIG_FIELD_FLOAT32" file:"CONFIG_FILE_FIELD_FLOAT32"`
	FieldFloat64 float64  `flag:"fieldFloat64" env:"CONFIG_FIELD_FLOAT64" file:"CONFIG_FILE_FIELD_FLOAT64"`
	FieldInt     int      `flag:"fieldInt" env:"CONFIG_FIELD_INT" file:"CONFIG_FILE_FIELD_INT"`
	FieldInt8    int8     `flag:"fieldInt8" env:"CONFIG_FIELD_INT8" file:"CONFIG_FILE_FIELD_INT8"`
	FieldInt16   int16    `flag:"fieldInt16" env:"CONFIG_FIELD_INT16" file:"CONFIG_FILE_FIELD_INT16"`
	FieldInt32   int32    `flag:"fieldInt32" env:"CONFIG_FIELD_INT32" file:"CONFIG_FILE_FIELD_INT32"`
	FieldInt64   int64    `flag:"fieldInt64" env:"CONFIG_FIELD_INT64" file:"CONFIG_FILE_FIELD_INT64"`
	FieldUint    uint     `flag:"fieldUint" env:"CONFIG_FIELD_UINT" file:"CONFIG_FILE_FIELD_UINT"`
	FieldUint8   uint8    `flag:"fieldUint8" env:"CONFIG_FIELD_UINT8" file:"CONFIG_FILE_FIELD_UINT8"`
	FieldUint16  uint16   `flag:"fieldUint16" env:"CONFIG_FIELD_UINT16" file:"CONFIG_FILE_FIELD_UINT16"`
	FieldUint32  uint32   `flag:"fieldUint32" env:"CONFIG_FIELD_UINT32" file:"CONFIG_FILE_FIELD_UINT32"`
	FieldUint64  uint64   `flag:"fieldUint64" env:"CONFIG_FIELD_UINT64" file:"CONFIG_FILE_FIELD_UINT64"`
}

func main() {
  spec := Spec{
    FieldString: "default value",
  }
  config.Pick(&spec)
  fmt.Printf("%+v\n", spec)
}
```
