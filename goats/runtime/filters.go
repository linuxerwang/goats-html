package runtime

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"
)

type BuiltinFilter struct{}

func NewBuiltinFilter() *BuiltinFilter {
	return &BuiltinFilter{}
}

/**
 * Prints an arbitrary variable's detailed information to help template authors
 * to debug. This is supposed to be used for dev mode only and thus it might have
 * a compromised performance.
 */
func (bf *BuiltinFilter) Debug(a interface{}) string {
	var buffer bytes.Buffer
	bf.inspect("data", a, &buffer)
	return buffer.String()
}

func (bf *BuiltinFilter) Length(a []interface{}) int {
	return len(a)
}

func (bf *BuiltinFilter) inspect(name string, a interface{}, writer io.Writer) {
	a = indirect(a)
	t := reflect.TypeOf(a)
	v := reflect.ValueOf(a)

	io.WriteString(writer, "<div style=\"padding-left:20px\">")
	// TODO: Improve it so that all level of composite types can be detailed.
	switch t.Kind() {
	case reflect.Struct:
		io.WriteString(writer, fmt.Sprintf("%s (%s): ", name, t.String()))
		io.WriteString(writer, "{<br>")
		for i := 0; i < v.NumField(); i++ {
			fieldVal := v.Field(i)
			fieldType := t.Field(i)
			bf.inspect(fieldType.Name, fieldVal.Interface(), writer)
		}
		io.WriteString(writer, "}<br>")
	default:
		io.WriteString(writer, fmt.Sprintf("%s (%s): %v<br>", name, t.String(), a))
	}
	io.WriteString(writer, "</div>")
}

func (bf *BuiltinFilter) Title(a interface{}) string {
	a = indirect(a)
	return strings.Title(fmt.Sprint(reflect.ValueOf(a)))
}

func (bf *BuiltinFilter) Capfirst(a interface{}) string {
	a = indirect(a)
	t := reflect.TypeOf(a)
	v := reflect.ValueOf(a)
	switch t.Kind() {
	case reflect.String:
		s := v.String()
		if len(s) > 0 {
			f, size := utf8.DecodeRuneInString(s)
			s = fmt.Sprint(strings.ToUpper(string(f)), s[size:])
		}
		return s
	}
	// TODO: do we need a util to print out the real value instead of type?
	return fmt.Sprint(v)
}

func (bf *BuiltinFilter) Center(a interface{}, width int) string {
	a = indirect(a)
	t := reflect.TypeOf(a)
	v := reflect.ValueOf(a)
	switch t.Kind() {
	case reflect.String:
		s := v.String()
		if len(s) < width {
			extra := width - len(s)
			left := extra / 2
			right := extra - left
			var buffer bytes.Buffer
			buffer.WriteString(generateString(left, " "))
			buffer.WriteString(s)
			buffer.WriteString(generateString(right, " "))
			s = buffer.String()
		}
		return s
	}
	// TODO: do we need a util to print out the real value instead of type?
	return fmt.Sprint(v)
}

func (bf *BuiltinFilter) Ljust(a interface{}, width int) string {
	return just(a, width, true)
}

func (bf *BuiltinFilter) Rjust(a interface{}, width int) string {
	return just(a, width, false)
}

func (bf *BuiltinFilter) Cut(a interface{}, removed string) string {
	a = indirect(a)
	s := a.(string)
	return strings.Replace(s, removed, "", -1)
}

func (bf *BuiltinFilter) Join(a interface{}, separator string) string {
	a = indirect(a)
	switch t := a.(type) {
	case string:
		return t
	case []string:
		return strings.Join(t, separator)
	}
	return fmt.Sprint(a)
}

func (bf *BuiltinFilter) FloatFormat(a interface{}, precision int) string {
	a = indirect(a)
	switch a := a.(type) {
	case float32:
		return bf.floatFormat(float64(a), precision)
	case float64:
		return bf.floatFormat(a, precision)
	default:
		return ""
	}
}

func (bf *BuiltinFilter) floatFormat(a float64, precision int) string {
	r := new(big.Rat)
	r.SetFloat64(a)
	if float64(int(a)) == a {
		if precision < 2 {
			precision = 0
		}
		return r.FloatString(precision)
	}
	if precision < 0 {
		precision = -precision
	}
	return r.FloatString(precision)
}

func (bf *BuiltinFilter) Quote(a interface{}) string {
	// TODO: Allow specify quote types (", ')
	return fmt.Sprintf("\"%s\"", indirect(a))
}

func (bf *BuiltinFilter) Format(fmtStr string, a ...interface{}) string {
	return fmt.Sprintf(fmtStr, a...)
}

func just(a interface{}, width int, left bool) string {
	a = indirect(a)
	s := a.(string)
	if len(s) < width {
		extra := width - len(s)
		var buffer bytes.Buffer
		if left {
			buffer.WriteString(s)
			buffer.WriteString(generateString(extra, " "))
		} else {
			buffer.WriteString(generateString(extra, " "))
			buffer.WriteString(s)
		}
		s = buffer.String()
	}
	return s
}

func generateString(length int, char string) string {
	var buffer bytes.Buffer
	for i := 0; i < length; i++ {
		buffer.WriteString(char)
	}
	return buffer.String()
}

// Example:
//     <div go:content="unixdate('01/02/2006 15:04:05'), 1390637053836182000"></div>
func (bf *BuiltinFilter) UnixDate(format string, a interface{}) string {
	return time.Unix(convertTime(a), 0).Format(format)
}

// Example:
//     <div go:content="unixnanodate('01/02/2006 15:04:05'), 1390637053836182000"></div>
func (bf *BuiltinFilter) UnixNanoDate(format string, a interface{}) string {
	return time.Unix(0, convertTime(a)).Format(format)
}

func convertTime(a interface{}) int64 {
	a = indirect(a)
	switch a := a.(type) {
	case int64:
		return a
	case uint64:
		return int64(a)
	}
	return 0
}
