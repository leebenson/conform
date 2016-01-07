# Santize library for Go structs (golang)
---------------------------------------

Trim, sanitize, and modify struct string fields in place, based on tags.

Turns this...

```
type Person struct {
	FirstName string `sanitize:"name"`
	LastName  string `sanitize:"ucfirst,trim"`
	Email     string `sanitize:"email"`
	CamelCase string `sanitize:"camel"`
	UserName  string `sanitize:"snake"`
	Slug      string `sanitize:"slug"`
	Blurb     string `sanitize:"title"`
	Left      string `sanitize:"ltrim"`
	Right     string `sanitize:"rtrim"`
}

p1 := Person{
	" LEE ",
	"     Benson",
	"   LEE@LEEbenson.com  ",
	"I love new york city",
	"lee benson",
	"LeeBensonWasHere",
	"this is a little bit about me...",
	"    Left trim   ",
	"    Right trim  ",
}

sanitize.Strings(&p)

```

Into this...

```
p2 := p1
sanitize.Strings(&p2)

/*
	p1 (left) vs. p2 (right)

	FirstName: ' LEE ' -> 'Lee'
	LastName: '     Benson' -> 'Benson'
	Email: '   LEE@LEEbenson.com  ' -> 'lee@leebenson.com'
	CamelCase: 'I love new york city' -> 'ILoveNewYorkCity'
	UserName: 'lee benson' -> 'lee_benson'
	Slug: 'LeeBensonWasHere' -> 'lee-benson-was-here'
	Blurb: 'this is a little bit about me...' -> 'This Is A Little Bit About Me...'
	Left: '    Left trim   ' -> 'Left trim   '
	Right: '    Right trim  ' -> '    Right trim'
*/
```

## Why?
---------------------------------------

Sanitize helps you fix and format user strings quickly, without writing functions.

If you do form processing with [Gorilla Schema](http://www.gorillatoolkit.org/pkg/schema) or similar, you probably shuttle user data into structs using tags. Adding a `sanitize` tag to your string field gives you "first pass" clean up against user input.

Use it for names, e-mail addresses, URL slugs, or any other form field where formatting matters.

Sanitize doesn't attempt any kind of validation on your fields. Check out [govalidator](https://github.com/asaskevich/govalidator).

## How to use
---------------------------------------

Grab the package with:

`go get github.com/leebenson/sanitize`

Here's an example that formats e-mail addresses:

```
package main

import (
		"fmt"
		"github.com/leebenson/sanitize"
)

type UserForm struct {
	Email string `sanitize:"email"`
}

func main() {
	input := UserForm{
		Email: "   POORLYFormaTTED@EXAMPlE.COM  "
	}
	sanitize.Strings(&input) // <-- pass in a pointer to your struct
	fmt.Println(input.Email) // prints "poorlyformatted@example.com"
}

```

## Using with Gorilla Schema
---------------------------------------

Just add a `sanitize` tag along with your Gorilla `schema` tags:

```
import (
		"github.com/gorilla/schema"
		"github.com/leebenson/sanitize"
)

// the struct that will be filled from the post request...
type newUserForm struct {
	FirstName 		string	`schema:"firstName",sanitize:"name"`
	Email			string	`schema:"emailAddress",sanitize:"email"`
	Password 		string	`schema:"password"` // <-- no change? no tag
	Dob				string	`schema:"dateOfBirth"` // <-- non-strings ignored by Sanitize
}

// ProcessNewUser attempts to register a new user
func ProcessNewUser(r *http.Request) error {
	form := new(newUserForm)
	schema.NewDecoder().Decode(form, r.PostForm) // <-- Gorilla Schema
	sanitize.Strings(form) <-- Sanitize.  Pass in the same pointer that Schema used
	// ...
}
```

## Tags
---------------------------------------

You can use multiple tags in the format of `sanitize:"tag1,tag2"`

### trim
---------------------------------------
Trims leading and trailing spaces. Example: `"   string   "` -> `"string"`

### ltrim
---------------------------------------
Trims leading spaces only. Example: `"   string   "` -> `"string   "`

### rtrim
---------------------------------------
Trims trailing spaces only. Example: `"   string   "` -> `"   string"`

### lower
---------------------------------------
Converts string to lowercase. Example: `"STRING"` -> `"string"`

### upper
---------------------------------------
Converts string to uppercase. Example: `"string"` -> `"STRING"`

### title
---------------------------------------
Converts string to Title Case, e.g. `"this is a sentence"` -> `"This Is A Sentence"`

### camel
---------------------------------------
Converts to camel case via [stringUp](https://github.com/etgryphon/stringUp), Example provided by library: `this is it => thisIsIt, this\_is\_it => thisIsIt, this-is-it => thisIsIt`

### snake
---------------------------------------
Converts to snake_case. Example: `CamelCase` -> `camel_case`, `regular string` -> `regular_string`
Special thanks to [snaker](https://github.com/serenize/snaker/) for inspiration (credited in license)

### slug
---------------------------------------
Turns strings into slugs.  Example: `CamelCase` -> `camel-case`, `blog title here` -> `blog-title-here`

### ucfirst
---------------------------------------
Uppercases first character.  Example: `all lower` -> `All lower`

### name
---------------------------------------
Trims, uppercases the first character, lowercases the rest. Example: ` JOHN ` -> `John`

### email
---------------------------------------
Trims and lowercases the string.  Example: `UNSIGHTLY-EMAIL@EXamPLE.com ` -> `unsightly-email@example.com`

### LICENSE
---------------------------------------

The MIT License (MIT)

Copyright (c) 2016 Lee Benson

Copyright (c) 2015 Serenize UG (haftungsbeschränkt)
(for code modified from https://github.com/serenize/snaker/)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
