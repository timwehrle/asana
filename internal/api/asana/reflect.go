package asana

import (
	"reflect"
	"strings"
)

// Fields gets all valid JSON fields for a type
func Fields(i interface{}) *Options {
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Struct {
		panic("Invalid type requested")
	}

	result := &Options{}

	gatherFields(t, result)

	return result
}

func gatherFields(t reflect.Type, result *Options) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// Process a tag if present
		if jsonTag := f.Tag.Get("json"); jsonTag != "" {
			name := strings.Split(jsonTag, ",")[0]
			if name != "" && name != "-" {
				result.Fields = append(result.Fields, name)
			}
		}

		if f.Anonymous {
			gatherFields(f.Type, result)
		}
	}
}
