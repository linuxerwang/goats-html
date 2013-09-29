package runtime

import (
	"fmt"
	"html"
	"reflect"
	"strconv"
	"strings"
)

var MERGEABLE_ATTRS = map[string]bool{
	"class": true,
	"style": true,
}

func IsMergeable(attrName string) bool {
	_, found := MERGEABLE_ATTRS[attrName]
	return found
}

// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

func EscapeAttr(attr string) string {
	if strings.Index(attr, "\"") > -1 {
		return strings.Replace(attr, "\"", "\\\"", -1)
	}
	return attr
}

func EscapeContent(a interface{}) string {
	if a == nil {
		return ""
	}
	a = indirect(a)
	switch val := a.(type) {
	case int:
		return strconv.Itoa(val)
	case int8:
		return strconv.Itoa(int(val))
	case int16:
		return strconv.Itoa(int(val))
	case int32:
		return strconv.Itoa(int(val))
	case int64:
		return strconv.Itoa(int(val))
	case uint:
		return strconv.Itoa(int(val))
	case uint8:
		return strconv.Itoa(int(val))
	case uint16:
		return strconv.Itoa(int(val))
	case uint32:
		return strconv.Itoa(int(val))
	case uint64:
		return strconv.Itoa(int(val))
	case string:
		return html.EscapeString(val)
	default:
		return html.EscapeString(fmt.Sprintf("%v", val))
	}
}
