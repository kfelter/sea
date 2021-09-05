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
