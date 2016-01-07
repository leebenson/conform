package conform

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/etgryphon/stringUp"
)

func camelTo(s, sep string) string {
	var result string
	var words []string
	var lastPos int
	rs := []rune(s)

	for i := 0; i < len(rs); i++ {
		if i > 0 && unicode.IsUpper(rs[i]) {
			if initialism := startsWithInitialism(s[lastPos:]); initialism != "" {
				words = append(words, initialism)

				i += len(initialism) - 1
				lastPos = i
				continue
			}

			words = append(words, s[lastPos:i])
			lastPos = i
		}
	}

	// append the last word
	if s[lastPos:] != "" {
		words = append(words, s[lastPos:])
	}

	for k, word := range words {
		if k > 0 {
			result += sep
		}

		result += strings.ToLower(word)
	}

	return result
}

// startsWithInitialism returns the initialism if the given string begins with it
func startsWithInitialism(s string) string {
	var initialism string
	// the longest initialism is 5 char, the shortest 2
	for i := 1; i <= 5; i++ {
		if len(s) > i-1 && commonInitialisms[s[:i]] {
			initialism = s[:i]
		}
	}
	return initialism
}

// commonInitialisms, taken from
// https://github.com/golang/lint/blob/3d26dc39376c307203d3a221bada26816b3073cf/lint.go#L482
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
}

func ucFirst(s string) string {
	if s == "" {
		return s
	}
	toRune, size := utf8.DecodeRuneInString(s)
	if !unicode.IsLower(toRune) {
		return s
	}
	buf := &bytes.Buffer{}
	buf.WriteRune(unicode.ToUpper(toRune))
	buf.WriteString(s[size:])
	return buf.String()
}

// Strings conforms strings based on reflection tags
func Strings(s interface{}) error {
	v := reflect.ValueOf(s)

	// Must be a pointer
	if v.Kind().String() != "ptr" {
		return errors.New("Not a pointer")
	}

	// Grab the type that the pointer points to
	r := reflect.Indirect(v).Type()

	// Range over the struct fields
	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)

		// Need a `conform:""` Tag
		t := f.Tag.Get("conform")
		if t == "" {
			continue
		}

		// Get the field by name
		n := v.Elem().FieldByName(f.Name)
		a := n.Addr()

		// Must be an exported field
		if a.CanInterface() {
			switch n.Interface().(type) {
			case string:
				// Get the current data
				d := n.String()

				// Range over tags, and perform changes
				for _, split := range strings.Split(t, ",") {
					switch split {
					case "trim":
						d = strings.TrimSpace(d)
					case "ltrim":
						d = strings.TrimLeft(d, " ")
					case "rtrim":
						d = strings.TrimRight(d, " ")
					case "lower":
						d = strings.ToLower(d)
					case "upper":
						d = strings.ToUpper(d)
					case "title":
						d = strings.Title(d)
					case "camel":
						d = stringUp.CamelCase(d)
					case "snake":
						d = camelTo(stringUp.CamelCase(d), "_")
					case "slug":
						d = camelTo(stringUp.CamelCase(d), "-")
					case "ucfirst":
						d = ucFirst(d)
					case "name":
						d = ucFirst(strings.ToLower(strings.TrimSpace(d)))
					case "email":
						d = strings.ToLower(strings.TrimSpace(d))
					}

				}
				n.SetString(d)
			}
		}
	}
	return nil
}
