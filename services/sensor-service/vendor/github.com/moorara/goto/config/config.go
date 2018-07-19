package config

import (
	"errors"
	"fmt"
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
)

func getFlagName(name string) string {
	subs := []string{}
	current := string(name[0])
	lastLower := unicode.IsLower(rune(name[0]))

	for i := 1; i < len(name); i++ {
		r := rune(name[i])
		s := string(name[i])

		if unicode.IsUpper(r) && lastLower {
			subs = append(subs, string(current))
			current = ""
		}

		current += s
		lastLower = unicode.IsLower(r)
	}

	subs = append(subs, string(current))
	result := strings.Join(subs, ".")
	result = strings.ToLower(result)

	return result
}

func getEnvVarName(name string) string {
	subs := []string{}
	current := string(name[0])
	lastLower := unicode.IsLower(rune(name[0]))

	for i := 1; i < len(name); i++ {
		r := rune(name[i])
		s := string(name[i])

		if unicode.IsUpper(r) && lastLower {
			subs = append(subs, string(current))
			current = ""
		}

		current += s
		lastLower = unicode.IsLower(r)
	}

	subs = append(subs, string(current))
	result := strings.Join(subs, "_")
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
		kind := vField.Kind()

		flagTag := tField.Tag.Get(flagTagName)
		if flagTag == "" {
			flagTag = getFlagName(name)
		}

		envTag := tField.Tag.Get(envTagName)
		if envTag == "" {
			envTag = getEnvVarName(name)
		}

		fileTag := tField.Tag.Get(fileTagName)
		if fileTag == "" {
			fileTag = envTag + "_FILE"
		}

		str := getFieldValue(flagTag, envTag, fileTag)

		if str != "" {
			switch kind {
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
				} else {
					fmt.Println(err)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if u, err := strconv.ParseUint(str, 10, 64); err == nil {
					vField.SetUint(u)
				}
			}
		}
	}

	return nil
}
