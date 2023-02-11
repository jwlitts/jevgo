package jevgo

import (
	"fmt"
	"testing"
)

var example = `
first name [John]
last name [Smith]
is alive [true]
age [27]
address [
street address [21 2nd Street]
city [New York]
state [NY]
postal code [10021-3100]
]
phone numbers [
[
	type [home]
	number [212 555-1234]
]
[
	type [office]
	number [646 555-4567]
]
]
children [seq]
spouse [nil]
`

func TestEasyParse(t *testing.T) {
	oneline := "first name [John]"
	_, jev, err := ParseJevko(oneline, false)
	if err != nil {
		fmt.Printf("%v\n", err)
		t.Fail()
	}
	fmt.Printf("%v\n", jev)
	if jev.Branches[0].Prefix != "first name " {
		t.FailNow()
	}
	if jev.Branches[0].Branch.Suffix != "John" {
		t.FailNow()
	}
}

func TestMediumParse(t *testing.T) {
	oneline := `
	first name [John]
last name [Smith]
is alive [true]
age [27]
address [
	street address [21 2nd Street]
	city [New York]
	state [NY]
	postal code [10021-3100]
]
	`
	_, jev, err := ParseJevko(oneline, false)
	if err != nil {
		fmt.Printf("%v\n", err)
		t.Fail()
	}
	fmt.Printf("%v\n", jev)
	t.Fail()
}

func TestEasyEscape(t *testing.T) {
	{
		oneline := "first name [Jo`[hn]"
		_, jev, err := ParseJevko(oneline, false)
		if err != nil {
			fmt.Printf("%v\n", err)
			t.Fail()
		}
		fmt.Printf("%v\n", jev)
		if jev.Branches[0].Prefix != "first name " {
			t.FailNow()
		}
		if jev.Branches[0].Branch.Suffix != "Jo`[hn" {
			t.FailNow()
		}
	}
	oneline := "first name [Jo```[hn]"
	_, jev, err := ParseJevko(oneline, false)
	if err != nil {
		fmt.Printf("%v\n", err)
		t.Fail()
	}
	fmt.Printf("%v\n", jev)
	if jev.Branches[0].Prefix != "first name " {
		t.FailNow()
	}
	if jev.Branches[0].Branch.Suffix != "Jo```[hn" {
		t.FailNow()
	}
}
