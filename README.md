## Simple Environemnt Abstraction

Implements the section of my book on configuration
Main features:
- handles typing of environment variables
- useful error messages
- document generation

## Installation

go install github.com/kfelter/sea/cmd/sea

## Usage

[example](examples/percent)

### Loading environment variables with defaults and types
```go
package main

import (
	"fmt"

	"github.com/kfelter/sea"
)

func main() {
	fmt.Println(sea.LoadWithDefault("MY_PERCENT", "50", "a percent to be used", "Int").Int())
	fmt.Println(sea.LoadWithDefault("USE_MY_PERCENT", "false", "bool if we should use the percent", "Boolean").Bool())
	fmt.Println(sea.Load("MY_NAME", "name does not need to be set", "String").String())
}
```

### Generating documentation for environment variables

`$ sea -root=example/percent`
#### examples/percent/main.go
| NAME | DEFAULT | USAGE | TYPE |
| --- | --- | --- | --- |
| `"MY_PERCENT"` | `"50"` | `"a percent to be used"` | `"Int"` |
| `"USE_MY_PERCENT"` | `"false"` | `"bool if we should use the percent"` | `"Boolean"` |
| `"MY_NAME"` | `nil` | `"name does not need to be set"` | `"String"` |

