# structtag [![](https://github.com/fatih/structtag/workflows/build/badge.svg)](https://github.com/fatih/structtag/actions)

structtag provides an easy way of parsing and manipulating struct tag fields.

# Install

```bash
go get github.com/fatih/structtag
```

# Example

```go
package main

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/fatih/structtag"
)

func main() {
	type t struct {
		t string `json:"foo,omitempty,string" xml:"foo"`
	}

	// get field tag
	tag := reflect.TypeOf(t{}).Field(0).Tag

	// ... and start using structtag by parsing the tag
	tags, err := structtag.Parse(string(tag))
	if err != nil {
		panic(err)
	}

	// iterate over all tags
	for _, t := range tags.Tags() {
		fmt.Printf("tag: %+v\n", t)
	}

	// get a single tag
	jsonTag, err := tags.Get("json")
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonTag)         // Output: json:"foo,omitempty,string"
	fmt.Println(jsonTag.Key)     // Output: json
	fmt.Println(jsonTag.Name)    // Output: foo
	fmt.Println(jsonTag.Options) // Output: [omitempty string]

	// change existing tag
	jsonTag.Name = "foo_bar"
	jsonTag.Options = nil
	tags.Set(jsonTag)

	// add new tag
	tags.Set(&structtag.Tag{
		Key:     "hcl",
		Name:    "foo",
		Options: []string{"squash"},
	})

	// print the tags
	fmt.Println(tags) // Output: json:"foo_bar" xml:"foo" hcl:"foo,squash"

	// sort tags according to keys
	sort.Sort(tags)
	fmt.Println(tags) // Output: hcl:"foo,squash" json:"foo_bar" xml:"foo"
}
```
