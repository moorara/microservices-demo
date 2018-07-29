package config

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const (
	flagTagName = "flag"
	envTagName  = "env"
	fileTagName = "file"
	sepTagName  = "sep"
)

func parseFieldName(name string) []string {
	tokens := []string{}
	current := string(name[0])
	lastLower := unicode.IsLower(rune(name[0]))

	add := func(slice []string, str string) []string {
		if str == "" {
			return slice
		}
		return append(slice, str)
	}

	for i := 1; i < len(name); i++ {
		r := rune(name[i])

		if unicode.IsUpper(r) && lastLower {
			// The case is changing from lower to upper
			tokens = add(tokens, current)
			current = string(name[i])
		} else if unicode.IsLower(r) && !lastLower {
			// The case is changing from upper to lower
			l := len(current) - 1
			tokens = add(tokens, current[:l])
			current = current[l:] + string(name[i])
		} else {
			// Increment current token
			current += string(name[i])
		}

		lastLower = unicode.IsLower(r)
	}

	tokens = append(tokens, string(current))

	return tokens
}

func getFlagName(name string) string {
	parts := parseFieldName(name)
	result := strings.Join(parts, ".")
	result = strings.ToLower(result)

	return result
}

func getEnvVarName(name string) string {
	parts := parseFieldName(name)
	result := strings.Join(parts, "_")
	result = strings.ToUpper(result)

	return result
}

func getFlagValue(flagName string) string {
	flagNameRegex := regexp.MustCompile("-{1,2}" + flagName)
	flagGeneralRegex := regexp.MustCompile("^-{1,2}[A-Za-z].*")

	for i, arg := range os.Args {
		if flagNameRegex.MatchString(arg) {
			if s := strings.Index(arg, "="); s > 0 {
				return arg[s+1:]
			}

			if i+1 < len(os.Args) {
				val := os.Args[i+1]
				if !flagGeneralRegex.MatchString(val) {
					return val
				}
			}

			return "true"
		}
	}

	return ""
}

func getFieldValue(flag, env, file string) string {
	var value string

	// First, try reading from flag
	value = getFlagValue(flag)

	// Second, try reading from environment variable
	if value == "" {
		value = os.Getenv(env)
	}

	// Third, try reading from file
	if value == "" {
		filepath := os.Getenv(file)
		if content, err := ioutil.ReadFile(filepath); err == nil {
			value = string(content)
		}
	}

	return value
}

func float32Slice(strs []string) []float32 {
	floats := []float32{}
	for _, str := range strs {
		if f, err := strconv.ParseFloat(str, 64); err == nil {
			floats = append(floats, float32(f))
		}
	}
	return floats
}

func float64Slice(strs []string) []float64 {
	floats := []float64{}
	for _, str := range strs {
		if f, err := strconv.ParseFloat(str, 64); err == nil {
			floats = append(floats, f)
		}
	}
	return floats
}

func intSlice(strs []string) []int {
	ints := []int{}
	for _, str := range strs {
		if i, err := strconv.ParseInt(str, 10, 64); err == nil {
			ints = append(ints, int(i))
		}
	}
	return ints
}

func int8Slice(strs []string) []int8 {
	ints := []int8{}
	for _, str := range strs {
		if i, err := strconv.ParseInt(str, 10, 64); err == nil {
			ints = append(ints, int8(i))
		}
	}
	return ints
}

func int16Slice(strs []string) []int16 {
	ints := []int16{}
	for _, str := range strs {
		if i, err := strconv.ParseInt(str, 10, 64); err == nil {
			ints = append(ints, int16(i))
		}
	}
	return ints
}

func int32Slice(strs []string) []int32 {
	ints := []int32{}
	for _, str := range strs {
		if i, err := strconv.ParseInt(str, 10, 64); err == nil {
			ints = append(ints, int32(i))
		}
	}
	return ints
}

func int64Slice(strs []string) []int64 {
	ints := []int64{}
	for _, str := range strs {
		if i, err := strconv.ParseInt(str, 10, 64); err == nil {
			ints = append(ints, i)
		}
	}
	return ints
}

func uintSlice(strs []string) []uint {
	uints := []uint{}
	for _, str := range strs {
		if u, err := strconv.ParseUint(str, 10, 64); err == nil {
			uints = append(uints, uint(u))
		}
	}
	return uints
}

func uint8Slice(strs []string) []uint8 {
	uints := []uint8{}
	for _, str := range strs {
		if u, err := strconv.ParseUint(str, 10, 64); err == nil {
			uints = append(uints, uint8(u))
		}
	}
	return uints
}

func uint16Slice(strs []string) []uint16 {
	uints := []uint16{}
	for _, str := range strs {
		if u, err := strconv.ParseUint(str, 10, 64); err == nil {
			uints = append(uints, uint16(u))
		}
	}
	return uints
}

func uint32Slice(strs []string) []uint32 {
	uints := []uint32{}
	for _, str := range strs {
		if u, err := strconv.ParseUint(str, 10, 64); err == nil {
			uints = append(uints, uint32(u))
		}
	}
	return uints
}

func uint64Slice(strs []string) []uint64 {
	uints := []uint64{}
	for _, str := range strs {
		if u, err := strconv.ParseUint(str, 10, 64); err == nil {
			uints = append(uints, u)
		}
	}
	return uints
}

// Pick reads values for a struct of specifications
func Pick(spec interface{}) error {
	v := reflect.ValueOf(spec).Elem() // reflect.Value --> v.Type(), v.Kind(), v.NumField()
	t := reflect.TypeOf(spec).Elem()  // reflect.Type --> t.Name(), t.Kind(), t.NumField()

	if t.Kind() != reflect.Struct {
		return errors.New("spec should be a struct type")
	}

	// Iterate over struct fields
	for i := 0; i < v.NumField(); i++ {
		vField := v.Field(i) // reflect.Value --> vField.Kind(), vField.Type().Name(), vField.Type().Kind(), vField.Interface()
		tField := t.Field(i) // reflect.StructField --> tField.Name, tField.Type.Name(), tField.Type.Kind(), tField.Tag.Get(tag)

		// Skip unexported fields
		if !vField.CanSet() {
			continue
		}

		name := tField.Name

		// `flag:""`
		flagTag := tField.Tag.Get(flagTagName)
		if flagTag == "" {
			flagTag = getFlagName(name)
		}

		// `env:""`
		envTag := tField.Tag.Get(envTagName)
		if envTag == "" {
			envTag = getEnvVarName(name)
		}

		// `file:""`
		fileTag := tField.Tag.Get(fileTagName)
		if fileTag == "" {
			fileTag = envTag + "_FILE"
		}

		// `sep:""`
		sepTag := tField.Tag.Get(sepTagName)
		if sepTag == "" {
			sepTag = ","
		}

		str := getFieldValue(flagTag, envTag, fileTag)
		if str == "" {
			continue
		}

		switch vField.Kind() {
		case reflect.String:
			vField.SetString(str)

		case reflect.Bool:
			if b, err := strconv.ParseBool(str); err == nil {
				vField.SetBool(b)
			}

		case reflect.Float32, reflect.Float64:
			if f, err := strconv.ParseFloat(str, 64); err == nil {
				vField.SetFloat(f)
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if i, err := strconv.ParseInt(str, 10, 64); err == nil {
				vField.SetInt(i)
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if u, err := strconv.ParseUint(str, 10, 64); err == nil {
				vField.SetUint(u)
			}

		case reflect.Slice:
			iSlice := vField.Interface()
			tSlice := reflect.TypeOf(iSlice).Elem()
			strs := strings.Split(str, sepTag)

			switch tSlice.Kind() {
			case reflect.String:
				vField.Set(reflect.ValueOf(strs))

			case reflect.Float32:
				floats := float32Slice(strs)
				vField.Set(reflect.ValueOf(floats))

			case reflect.Float64:
				floats := float64Slice(strs)
				vField.Set(reflect.ValueOf(floats))

			case reflect.Int:
				ints := intSlice(strs)
				vField.Set(reflect.ValueOf(ints))

			case reflect.Int8:
				ints := int8Slice(strs)
				vField.Set(reflect.ValueOf(ints))

			case reflect.Int16:
				ints := int16Slice(strs)
				vField.Set(reflect.ValueOf(ints))

			case reflect.Int32:
				ints := int32Slice(strs)
				vField.Set(reflect.ValueOf(ints))

			case reflect.Int64:
				ints := int64Slice(strs)
				vField.Set(reflect.ValueOf(ints))

			case reflect.Uint:
				uints := uintSlice(strs)
				vField.Set(reflect.ValueOf(uints))

			case reflect.Uint8:
				uints := uint8Slice(strs)
				vField.Set(reflect.ValueOf(uints))

			case reflect.Uint16:
				uints := uint16Slice(strs)
				vField.Set(reflect.ValueOf(uints))

			case reflect.Uint32:
				uints := uint32Slice(strs)
				vField.Set(reflect.ValueOf(uints))

			case reflect.Uint64:
				uints := uint64Slice(strs)
				vField.Set(reflect.ValueOf(uints))
			}
		}
	}

	return nil
}
