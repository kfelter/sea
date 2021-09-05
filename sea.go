package sea

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

/*
ex...
myPercent := sea.Load("MY_PERCENT").Must().Int()
this will cause a panic if MY_PERCENT is not set in the environment
this will cause a panic if MY_PERCENT is set, but cannot be converted to an int
this will return an int if MY_PERCENT is set and can be converted to an int

myPercent := sea.LoadWithDefault("MY_PERCENT", "51", "a percent to use in my app", "Int").Int()
if no value is set for MY_PERCENT, the value "51" will be parsed and returned
if MY_PERCENT is set, the value from the environment will be parsed and returned
this function also accepts a description that is not used in the runtime
but is used for document generation using the sea cli, you will also want to provide the type after the description
*/

// V used to hold the name and value of an environment variable
// you should not use this struct directly, instead create it using
// Load or LoadWithDefault then use the methods on V to parse it
type V struct {
	Name  string
	Value string
}

// Load read the os.Getenv value and return the V type
// description is optional if you wish to generate documentation with the sea package a description should be provided
func Load(name string, description ...string) V {
	return V{name, os.Getenv(name)}
}

// Load read the os.Getenv value, if no value is present then use the value passed and return the V type
// description is optional if you wish to generate documentation with the sea package a description should be provided
// as arg 2 and a type should be provided as arg 3
func LoadWithDefault(name, value string, description ...string) V {
	var ev string
	if ev = os.Getenv(name); ev != "" {
		return V{name, ev}
	}
	return V{name, value}
}

// Must if v.Value is not set when this function is called it will cause a panic
func (v V) Must() V {
	if v.Value == "" {
		panic(fmt.Errorf("%s must have a value", v.Name))
	}
	return v
}

// Int convert the environment value to an int type
func (v V) Int() int {
	var (
		value int
		err   error
	)
	if value, err = strconv.Atoi(v.Value); err != nil {
		panic(fmt.Errorf("%s must be an integer: %v", v.Name, err))
	}
	return value
}

// Float convert the environemt value to a float32
func (v V) Float() float64 {
	var (
		value float64
		err   error
	)
	if value, err = strconv.ParseFloat(v.Value, 64); err != nil {
		panic(fmt.Errorf("%s must be a 64 bit float: %v", v.Name, err))
	}
	return value
}

// String returns the raw string value
func (v V) String() string {
	return v.Value
}

// Bool converts the environment value to a boolean
func (v V) Bool() bool {
	var (
		value bool
		err   error
	)
	if value, err = strconv.ParseBool(v.Value); err != nil {
		panic(fmt.Errorf("%s must be a boolean: %v", v.Name, err))
	}
	return value
}

// Base64Decode will attempt to base64 decode the value found in the environment
func (v V) Base64Decode() []byte {
	var (
		decoded []byte
		err     error
	)
	if decoded, err = base64.RawStdEncoding.DecodeString(v.Value); err != nil {
		panic(fmt.Errorf("%s must be base64 encoded text: %v", v.Name, err))
	}
	return decoded
}

// TimeDuration will parse a time duration using time.ParseDuration (https://pkg.go.dev/time#ParseDuration)
func (v V) TimeDuration() time.Duration {
	var (
		dur time.Duration
		err error
	)
	if dur, err = time.ParseDuration(v.Value); err != nil {
		panic(fmt.Errorf(`%s must be a time duration (https://pkg.go.dev/time#ParseDuration): %v`, v.Name, err))
	}
	return dur
}

// Time will parse a time using the format provided
func (v V) Time(format string) time.Time {
	var (
		t   time.Time
		err error
	)
	if t, err = time.Parse(format, v.Value); err != nil {
		panic(fmt.Errorf("%s must be a time in format %s: %v", v.Name, format, err))
	}
	return t
}

// Byte convert env value to bytes
func (v V) Byte() []byte {
	return []byte(v.Value)
}

// Parse allows for any custom parsing not implemented by this package
// the parser function should accept the sea.V type and a pointer
// your parser function should take the value of v.Value
// and perform some operation then set the pointer i to the parsed value
// ex. sea.JSON
func (v V) Parse(parser func(v V, i interface{}), i interface{}) {
	parser(v, i)
}

// JSON attempts to json.Unmarshal v.Value into the pointer i
// utilizes the custom parsing method as an example, it takes the pointer i
// creates a custom parsing function that we pass to v.Parse
// it then runs the function and unmarshals the json into the pointer
// this pattern could be followed for any custom parsing you wish to use
func (v V) JSON(i interface{}) {
	v.Parse(func(v V, ptr interface{}) {
		if err := json.Unmarshal([]byte(v.Value), ptr); err != nil {
			panic(fmt.Errorf("%s must be json: %v", v.Name, err))
		}
	}, i)
}
