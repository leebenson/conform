package conform

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	RegExTrim    *regexp.Regexp
	RegExLTrim   *regexp.Regexp
	RegExRTrim   *regexp.Regexp
	RegExLower   *regexp.Regexp
	RegExUpper   *regexp.Regexp
	RegExTitle   *regexp.Regexp
	RegExUCFirst *regexp.Regexp
	RegExEmail   *regexp.Regexp
	RegExCamel   *regexp.Regexp
	RegExSnake   *regexp.Regexp
	RegExSlug    *regexp.Regexp
}

func (t *testSuite) leftPadding() string {
	return strings.Repeat(" ", rand.Intn(100))
}

func (t *testSuite) rightPadding() string {
	return strings.Repeat(" ", rand.Intn(100))
}

func (t *testSuite) padding(s string) string {
	return t.leftPadding() + s + t.rightPadding()
}

func (t *testSuite) SetupTest() {
	t.RegExTrim = regexp.MustCompile("^[^\\s].+[^\\s]$")
	t.RegExLTrim = regexp.MustCompile("^[^\\s]")
	t.RegExRTrim = regexp.MustCompile("[^\\s]$")
	t.RegExLower = regexp.MustCompile("^[a-z]+$")
	t.RegExUpper = regexp.MustCompile("^[A-Z]+$")
	t.RegExTitle = regexp.MustCompile("^[A-Z][a-z\\.]*([\\s][A-Z][a-z\\.]*)+$")
	t.RegExUCFirst = regexp.MustCompile("^[A-Z][a-z]+$")
	t.RegExEmail = regexp.MustCompile("^[^A-Z\\s]+$")
	t.RegExCamel = regexp.MustCompile("[A-Z]([A-Z0-9]*[a-z][a-z0-9]*[A-Z]|[a-z0-9]*[A-Z][A-Z0-9]*[a-z])[A-Za-z0-9]*")
	t.RegExSnake = regexp.MustCompile("^[a-z]+_[a-z]+$")
	t.RegExSlug = regexp.MustCompile("^[a-z]+-[a-z]+$")
}

func (t *testSuite) TestTrim() {
	assert := assert.New(t.T())

	for i := 0; i < 10000; i++ {

		var s struct {
			nothing   string `conform:"trim"`
			FirstName string `conform:"trim"`
			LastName  string `conform:"trim"`
		}

		s.FirstName = t.padding(fake.FirstName())
		s.LastName = t.padding(fake.LastName())

		Strings(&s)
		if ok := assert.Regexp(t.RegExTrim, s.FirstName, "First name not trimmed"); !ok {
			break
		}

		if ok := assert.Regexp(t.RegExTrim, s.LastName, "Last name not trimmed"); !ok {
			break
		}
	}
}

func (t *testSuite) TestLeftTrim() {
	assert := assert.New(t.T())

	for i := 0; i < 10000; i++ {

		var s struct {
			nothing   string `conform:"ltrim"`
			FirstName string `conform:"ltrim"`
			LastName  string `conform:"ltrim"`
		}

		s.FirstName = t.padding(fake.FirstName())
		s.LastName = t.padding(fake.LastName())

		Strings(&s)
		if ok := assert.Regexp(t.RegExLTrim, s.FirstName, "First name should be left trimmed"); !ok {
			break
		}
		if ok := assert.Regexp(t.RegExLTrim, s.LastName, "Last name should be left trimmed"); !ok {
			break
		}
	}
}

func (t *testSuite) TestRightTrim() {
	assert := assert.New(t.T())

	for i := 0; i < 10000; i++ {

		var s struct {
			nothing   string `conform:"rtrim"`
			FirstName string `conform:"rtrim"`
			LastName  string `conform:"rtrim"`
		}

		s.FirstName = t.padding(fake.FirstName())
		s.LastName = t.padding(fake.LastName())

		Strings(&s)
		if ok := assert.Regexp(t.RegExRTrim, s.FirstName, "First name should be right trimmed"); !ok {
			break
		}
		if ok := assert.Regexp(t.RegExRTrim, s.LastName, "Last name should be right trimmed"); !ok {
			break
		}
	}
}

func (t *testSuite) TestNoChanges() {
	assert := assert.New(t.T())

	for i := 0; i < 10; i++ {
		var s struct {
			FirstName string
			LastName  string
		}

		fn := t.padding(fake.FirstName())
		ln := t.padding(fake.LastName())

		s.FirstName = fn
		s.LastName = ln
		Strings(&s)
		if ok := assert.Equal(s.FirstName, fn, "First name shouldn't change"); !ok {
			break
		}
		if ok := assert.Equal(s.LastName, ln, "Last name shouldn't change"); !ok {
			break
		}
	}
}

func (t *testSuite) TestSomeChanges() {
	assert := assert.New(t.T())

	for i := 0; i < 10; i++ {
		var s struct {
			FirstName string
			LastName  string `conform:"trim"`
		}

		fn := t.padding(fake.FirstName())
		ln := t.padding(fake.LastName())

		s.FirstName = fn
		s.LastName = ln
		Strings(&s)
		if ok := assert.Equal(s.FirstName, fn, "First name shouldn't change"); !ok {
			break
		}
		if ok := assert.Regexp(t.RegExTrim, s.LastName, "Last name not trimmed"); !ok {
			break
		}
	}
}

func (t *testSuite) TestLower() {
	assert := assert.New(t.T())

	for i := 0; i < 10000; i++ {
		var s struct {
			FirstName string `conform:"lower"`
			LastName  string `conform:"lower"`
		}
		s.FirstName = strings.ToUpper(fake.FirstName())
		s.LastName = strings.ToUpper(fake.LastName())
		Strings(&s)
		if ok := assert.Regexp(t.RegExLower, s.FirstName, "First name should be lowercase"); !ok {
			break
		}
		if ok := assert.Regexp(t.RegExLower, s.LastName, "Last name should be lowercase"); !ok {
			break
		}
	}
}

func (t *testSuite) TestUpper() {
	assert := assert.New(t.T())

	for i := 0; i < 10000; i++ {
		var s struct {
			FirstName string `conform:"upper"`
			LastName  string `conform:"upper"`
		}
		s.FirstName = strings.ToLower(fake.FirstName())
		s.LastName = strings.ToLower(fake.LastName())
		Strings(&s)
		if ok := assert.Regexp(t.RegExUpper, s.FirstName, "First name should be uppercase"); !ok {
			break
		}
		if ok := assert.Regexp(t.RegExUpper, s.LastName, "Last name should be uppercase"); !ok {
			break
		}
	}
}

func (t *testSuite) TestCamel() {
	assert := assert.New(t.T())

	for i := 0; i < 10000; i++ {
		var s struct {
			Dashes      string `conform:"camel"`
			Underscores string `conform:"camel"`
			Spaces      string `conform:"camel"`
		}
		s.Dashes = fmt.Sprintf("%s-%s", fake.FirstName(), fake.LastName())
		s.Underscores = fmt.Sprintf("%s_%s", fake.FirstName(), fake.LastName())
		s.Spaces = fmt.Sprintf("%s %s", fake.FirstName(), fake.LastName())
		Strings(&s)
		if ok := assert.Regexp(t.RegExCamel, s.Dashes, "Dashes should be CamelCased"); !ok {
			break
		}
		if ok := assert.Regexp(t.RegExCamel, s.Underscores, "Underscores should be CamelCased"); !ok {
			break
		}
		if ok := assert.Regexp(t.RegExCamel, s.Spaces, "Spaces should be CamelCased"); !ok {
			break
		}
	}

}

func (t *testSuite) TestSnake() {
	assert := assert.New(t.T())

	for i := 0; i < 10000; i++ {
		var s struct {
			Camel  string `conform:"snake"`
			Spaces string `conform:"snake"`
		}
		s.Camel = fmt.Sprintf("%s%s", fake.FirstName(), fake.LastName())
		s.Spaces = fmt.Sprintf("%s %s", fake.FirstName(), fake.LastName())
		Strings(&s)
		if ok := assert.Regexp(t.RegExSnake, s.Camel, "CamelCase should be snake_case"); !ok {
			break
		}
		if ok := assert.Regexp(t.RegExSnake, s.Spaces, "Spaces should be snake_case"); !ok {
			break
		}
	}

}

func (t *testSuite) TestSlug() {
	assert := assert.New(t.T())

	for i := 0; i < 10000; i++ {
		var s struct {
			Camel  string `conform:"slug"`
			Spaces string `conform:"slug"`
		}
		s.Camel = fmt.Sprintf("%s%s", fake.FirstName(), fake.LastName())
		s.Spaces = fmt.Sprintf("%s %s", fake.FirstName(), fake.LastName())
		Strings(&s)
		if ok := assert.Regexp(t.RegExSlug, s.Camel, "CamelCase should be slug-case"); !ok {
			break
		}
		if ok := assert.Regexp(t.RegExSlug, s.Spaces, "Spaces should be slug-case"); !ok {
			break
		}
	}

}

func (t *testSuite) TestTitle() {
	assert := assert.New(t.T())

	for i := 0; i < 10000; i++ {
		var s struct {
			FullName string `conform:"title"`
		}
		s.FullName = strings.ToLower(fake.FullName())
		Strings(&s)
		if ok := assert.Regexp(t.RegExTitle, s.FullName, "Full name should be Title Cased"); !ok {
			break
		}
	}
}

func (t *testSuite) TestUpperFirst() {
	assert := assert.New(t.T())

	for i := 0; i < 10000; i++ {
		var s struct {
			FirstName string `conform:"ucfirst"`
			LastName  string `conform:"ucfirst"`
		}
		s.FirstName = strings.ToLower(fake.FirstName())
		s.LastName = strings.ToLower(fake.LastName())
		Strings(&s)
		if ok := assert.Regexp(t.RegExUCFirst, s.FirstName, "First name should be uppercase first, lower rest"); !ok {
			break
		}
		if ok := assert.Regexp(t.RegExUCFirst, s.LastName, "Last name should be uppercase first, lower rest"); !ok {
			break
		}
	}
}

func (t *testSuite) TestMixed() {
	assert := assert.New(t.T())

	for i := 0; i < 10; i++ {
		var s struct {
			Email     string `conform:"trim,lower"`
			FirstName string `conform:"trim"`
			LastName  string `conform:"trim"`
			Age       int
			Truth     bool
		}

		s.Email = t.padding(fake.EmailAddress())
		s.FirstName = t.padding(fake.FirstName())
		s.LastName = t.padding(fake.LastName())
		Strings(&s)
		if ok := assert.Regexp(t.RegExEmail, s.Email, "E-mail should be lower and trimmed"); !ok {
			break
		}
		if ok := assert.Regexp(t.RegExTrim, s.FirstName, "First name should be trimmed"); !ok {
			break
		}
		if ok := assert.Regexp(t.RegExTrim, s.LastName, "First name should be trimmed"); !ok {
			break
		}
	}
}

func TestStrings(t *testing.T) {
	suite.Run(t, new(testSuite))
}
