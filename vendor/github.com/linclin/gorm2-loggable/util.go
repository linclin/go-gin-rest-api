package loggable

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

func isEqual(item1, item2 interface{}, except ...string) bool {
	except = StringMap(except, ToSnakeCase)
	m1, m2 := somethingToMapStringInterface(item1), somethingToMapStringInterface(item2)
	if len(m1) != len(m2) {
		return false
	}
	for k, v := range m1 {
		if isInStringSlice(ToSnakeCase(k), except) {
			continue
		}
		v2, ok := m2[k]
		if !ok || !reflect.DeepEqual(v, v2) {
			return false
		}
	}
	return true
}

func somethingToMapStringInterface(item interface{}) map[string]interface{} {
	if item == nil {
		return nil
	}
	switch raw := item.(type) {
	case string:
		return somethingToMapStringInterface([]byte(raw))
	case []byte:
		var m map[string]interface{}
		err := json.Unmarshal(raw, &m)
		if err != nil {
			return nil
		}
		return m
	default:
		data, err := json.Marshal(item)
		if err != nil {
			return nil
		}
		return somethingToMapStringInterface(data)
	}
	return nil
}

var ToSnakeCase = toSomeCase("_")

func toSomeCase(sep string) func(string) string {
	return func(s string) string {
		for i := range s {
			if unicode.IsUpper(rune(s[i])) {
				if i != 0 {
					s = strings.Join([]string{s[:i], ToLowerFirst(s[i:])}, sep)
				} else {
					s = ToLowerFirst(s)
				}
			}
		}
		return s
	}
}

func ToLowerFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

func StringMap(strs []string, fn func(string) string) []string {
	res := make([]string, len(strs))
	for i := range strs {
		res[i] = fn(strs[i])
	}
	return res
}

func isInStringSlice(what string, where []string) bool {
	for i := range where {
		if what == where[i] {
			return true
		}
	}
	return false
}

func getLoggableFieldNames(value interface{}) []string {
	var names []string

	t := reflect.TypeOf(value).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value, ok := field.Tag.Lookup(loggableTag)
		if ok && value == "disable" {
			continue
		}

		names = append(names, field.Name)
	}

	return names
}

var matchFirstCap = regexp.MustCompile("([A-Z])([A-Z][a-z])")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase converts the provided string to snake_case.
// Based on https://gist.github.com/stoewer/fbe273b711e6a06315d19552dd4d33e6
func ToSnakeCaseRegEx(input string) string {
	output := matchFirstCap.ReplaceAllString(input, "${1}_${2}")
	output = matchAllCap.ReplaceAllString(output, "${1}_${2}")
	output = strings.ReplaceAll(output, "-", "_")
	return strings.ToLower(output)
}
