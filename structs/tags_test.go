package structs_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gookit/goutil"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/goutil/testutil/assert"
)

func ExampleTagParser_Parse() {
	type User struct {
		Age   int    `json:"age" yaml:"age" default:"23"`
		Name  string `json:"name,omitempty" yaml:"name" default:"inhere"`
		inner string
	}

	u := &User{}
	p := structs.NewTagParser("json", "yaml", "default")
	goutil.MustOK(p.Parse(u))

	tags := p.Tags()
	dump.P(tags)
	/*tags:
	map[string]maputil.SMap { #len=2
	  "Age": maputil.SMap { #len=3
	    "json": string("age"), #len=3
	    "yaml": string("age"), #len=3
	    "default": string("23"), #len=2
	  },
	  "Name": maputil.SMap { #len=3
	    "default": string("inhere"), #len=6
	    "json": string("name,omitempty"), #len=14
	    "yaml": string("name"), #len=4
	  },
	},
	*/

	dump.P(p.Info("name", "json"))
	/*info:
	maputil.SMap { #len=2
	  "name": string("name"), #len=4
	  "omitempty": string("true"), #len=4
	},
	*/

	fmt.Println(
		tags["Age"].Get("json"),
		tags["Age"].Get("default"),
	)

	// Output:
	// age 23
}

func ExampleTagParser_Parse_parseTagValueDefine() {
	// eg: "desc;required;default;shorts"
	type MyCmd struct {
		Name string `flag:"set your name;false;INHERE;n"`
	}

	c := &MyCmd{}
	p := structs.NewTagParser("flag")

	sepStr := ";"
	defines := []string{"desc", "required", "default", "shorts"}
	p.ValueFunc = structs.ParseTagValueDefine(sepStr, defines)

	goutil.MustOK(p.Parse(c))
	// dump.P(p.Tags())
	/*
		map[string]maputil.SMap { #len=1
		  "Name": maputil.SMap { #len=1
		    "flag": string("set your name;false;INHERE;n"), #len=28
		  },
		},
	*/
	fmt.Println("tags:", p.Tags())

	info, _ := p.Info("Name", "flag")
	dump.P(info)
	/*
		maputil.SMap { #len=4
		  "desc": string("set your name"), #len=13
		  "required": string("false"), #len=5
		  "default": string("INHERE"), #len=6
		  "shorts": string("n"), #len=1
		},
	*/

	// Output:
	// tags: map[Name:{flag:set your name;false;INHERE;n}]
}

func TestParseTagValueINI(t *testing.T) {
	mp, err := structs.ParseTagValueNamed("name", "")
	assert.NoErr(t, err)
	assert.Empty(t, mp)

	mp, err = structs.ParseTagValueNamed("name", "default=inhere")
	assert.NoErr(t, err)
	assert.NotEmpty(t, mp)
	assert.Eq(t, "inhere", mp.Str("default"))
}

func TestParseTags(t *testing.T) {
	type user struct {
		Age   int    `json:"age" default:"23"`
		Name  string `json:"name" default:"inhere"`
		inner string
	}

	tags, err := structs.ParseTags(user{}, []string{"json", "default"})
	assert.NoErr(t, err)
	assert.NotEmpty(t, tags)
	assert.NotContains(t, tags, "inner")

	assert.Contains(t, tags, "Age")
	assert.Eq(t, "age", tags["Age"].Str("json"))
	assert.Eq(t, 23, tags["Age"].Int("default"))

	assert.Contains(t, tags, "Name")
	assert.Eq(t, "name", tags["Name"].Str("json"))
	assert.Eq(t, 0, tags["Name"].Int("default"))
}

func TestParseReflectTags(t *testing.T) {
	type user struct {
		Age   int    `json:"age" default:"23"`
		Name  string `json:"name" default:"inhere"`
		inner string
	}

	rt := reflect.TypeOf(user{})
	tags, err := structs.ParseReflectTags(rt, []string{"json", "default"})
	assert.NoErr(t, err)
	assert.NotEmpty(t, tags)
	assert.NotContains(t, tags, "inner")

	assert.Contains(t, tags, "Age")
	assert.Eq(t, "age", tags["Age"].Str("json"))
	assert.Eq(t, 23, tags["Age"].Int("default"))

	assert.Contains(t, tags, "Name")
	assert.Eq(t, "name", tags["Name"].Str("json"))
	assert.Eq(t, 0, tags["Name"].Int("default"))
}
