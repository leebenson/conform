package conform

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/etgryphon/stringUp"
)

type x map[string]string

type sanitizer func(string) string

var sanitizers = map[string]sanitizer{}

var patterns = map[string]*regexp.Regexp{
	"numbers":    regexp.MustCompile("[0-9]"),
	"nonNumbers": regexp.MustCompile("[^0-9]"),
	"alpha":      regexp.MustCompile("[\\pL]"),
	"nonAlpha":   regexp.MustCompile("[^\\pL]"),
	"name":       regexp.MustCompile("[\\p{L}]([\\p{L}|[:space:]|\\-|\\']*[\\p{L}])*"),
}

var truncateRegex = regexp.MustCompile(`^truncate=([0-9]+)$`)

// a valid email will only have one "@", but let's treat the last "@" as the domain part separator
func emailLocalPart(s string) string {
	i := strings.LastIndex(s, "@")
	if i == -1 {
		return s
	}
	return s[0:i]
}

func emailDomainPart(s string) string {
	i := strings.LastIndex(s, "@")
	if i == -1 {
		return ""
	}
	return s[i+1:]
}

func email(s string) string {
	// According to rfc5321, "The local-part of a mailbox MUST BE treated as case sensitive"
	return emailLocalPart(s) + "@" + strings.ToLower(emailDomainPart(s))
}

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

func onlyNumbers(s string) string {
	return patterns["nonNumbers"].ReplaceAllLiteralString(s, "")
}

func stripNumbers(s string) string {
	return patterns["numbers"].ReplaceAllLiteralString(s, "")
}

func onlyAlpha(s string) string {
	return patterns["nonAlpha"].ReplaceAllLiteralString(s, "")
}

func stripAlpha(s string) string {
	return patterns["alpha"].ReplaceAllLiteralString(s, "")
}

func onlyOne(s string, m []x) string {
	for _, v := range m {
		for f, r := range v {
			s = regexp.MustCompile(fmt.Sprintf("%s", f)).ReplaceAllLiteralString(s, r)
		}
	}
	return s
}

func formatName(s string) string {
	first := onlyOne(strings.ToLower(s), []x{
		{"[^\\pL-\\s']": ""}, // cut off everything except [ alpha, hyphen, whitespace, apostrophe]
		{"\\s{2,}": " "},     // trim more than two whitespaces to one
		{"-{2,}": "-"},       // trim more than two hyphens to one
		{"'{2,}": "'"},       // trim more than two apostrophes to one
		{"( )*-( )*": "-"},   // trim enclosing whitespaces around hyphen
	})
	return strings.Title(patterns["name"].FindString(first))
}

func getSliceElemType(t reflect.Type) reflect.Type {
	var elType reflect.Type
	if t.Kind() == reflect.Ptr {
		elType = t.Elem().Elem()
	} else {
		elType = t.Elem()
	}

	return elType
}

func transformValue(tags string, val reflect.Value) reflect.Value {
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return val
	}

	var oldStr string
	if val.Kind() == reflect.Ptr {
		oldStr = val.Elem().String()
	} else {
		oldStr = val.String()
	}

	newStr := transformString(oldStr, tags)

	var newVal reflect.Value
	if val.Kind() == reflect.Ptr {
		newVal = reflect.ValueOf(&newStr)
	} else {
		newVal = reflect.ValueOf(newStr)
	}

	return newVal.Convert(val.Type())
}

// Strings conforms strings based on reflection tags
func Strings(iface interface{}) error {
	ifv := reflect.ValueOf(iface)
	if ifv.Kind() != reflect.Ptr {
		return errors.New("Not a pointer")
	}
	ift := reflect.Indirect(ifv).Type()
	if ift.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < ift.NumField(); i++ {
		v := ift.Field(i)
		el := reflect.Indirect(ifv.Elem().FieldByName(v.Name))
		switch el.Kind() {
		case reflect.Slice:
			if el.CanInterface() {
				elType := getSliceElemType(v.Type)

				// allow strings and string pointers
				str := ""
				if (elType.ConvertibleTo(reflect.TypeOf(str)) && reflect.TypeOf(str).ConvertibleTo(elType)) ||
					(elType.ConvertibleTo(reflect.TypeOf(&str)) && reflect.TypeOf(&str).ConvertibleTo(elType)) {
					tags := v.Tag.Get("conform")
					for i := 0; i < el.Len(); i++ {
						el.Index(i).Set(transformValue(tags, el.Index(i)))
					}
				} else {
					val := reflect.ValueOf(el.Interface())
					for i := 0; i < val.Len(); i++ {
						elVal := val.Index(i)
						if elVal.Kind() != reflect.Ptr {
							elVal = elVal.Addr()
						}
						Strings(elVal.Interface())
					}
				}
			}
		case reflect.Map:
			if el.CanInterface() {
				val := reflect.ValueOf(el.Interface())
				for _, key := range val.MapKeys() {
					mapValue := val.MapIndex(key)
					mapValuePtr := reflect.New(mapValue.Type())
					mapValuePtr.Elem().Set(mapValue)
					if mapValuePtr.Elem().CanAddr() {
						Strings(mapValuePtr.Elem().Addr().Interface())
					}
					val.SetMapIndex(key, reflect.Indirect(mapValuePtr))
				}
			}
		case reflect.Struct:
			if el.CanAddr() && el.Addr().CanInterface() {
				// To handle "sql.NullString" we can assume that tags are added to a field of type struct rather than string
				if tags := v.Tag.Get("conform"); tags != "" && el.CanSet() {
					field := el.FieldByName("String")
					str := field.String()
					field.SetString(transformString(str, tags))
				} else {
					Strings(el.Addr().Interface())
				}
			}
		case reflect.String:
			if el.CanSet() {
				tags := v.Tag.Get("conform")
				input := el.String()
				el.SetString(transformString(input, tags))
			}
		}
	}
	return nil
}

func transformString(input, tags string) string {
	if tags == "" {
		return input
	}
	for _, split := range strings.Split(tags, ",") {
		switch split {
		case "trim":
			input = strings.TrimSpace(input)
		case "ltrim":
			input = strings.TrimLeft(input, " ")
		case "rtrim":
			input = strings.TrimRight(input, " ")
		case "lower":
			input = strings.ToLower(input)
		case "upper":
			input = strings.ToUpper(input)
		case "title":
			input = strings.Title(input)
		case "camel":
			input = stringUp.CamelCase(input)
		case "snake":
			input = camelTo(stringUp.CamelCase(input), "_")
		case "slug":
			input = camelTo(stringUp.CamelCase(input), "-")
		case "ucfirst":
			input = ucFirst(input)
		case "name":
			input = formatName(input)
		case "email":
			input = email(strings.TrimSpace(input))
		case "num":
			input = onlyNumbers(input)
		case "!num":
			input = stripNumbers(input)
		case "alpha":
			input = onlyAlpha(input)
		case "!alpha":
			input = stripAlpha(input)
		case "!html":
			input = template.HTMLEscapeString(input)
		case "!js":
			input = template.JSEscapeString(input)
		default:
			if truncateParam := truncateRegex.FindString(split); truncateParam != "" {
				l, err := strconv.ParseInt(strings.TrimLeft(truncateParam, "truncate="), 10, 32)
				if err == nil && utf8.RuneCountInString(input) >= int(l) {
					r := []rune(input)
					return string(r[:l])
				}
			}
			if s, ok := sanitizers[split]; ok {
				input = s(input)
			}
		}
	}
	return input
}

// AddSanitizer associates a sanitizer with a key, which can be used in a Struct tag
func AddSanitizer(key string, s sanitizer) {
	sanitizers[key] = s
}
