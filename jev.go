package jevgo

import "fmt"

var examplePEG = `
Jevko ::= Subjevko* Suffix
Subjevko ::= Prefix Opener Jevko Closer

Opener ::= "["
Closer ::= "]"` +
	"Escaper ::= \"`\"" + `
Delimiter ::= Escaper | Opener | Closer

Suffix ::= Text
Prefix ::= Text

Text ::= Symbol*
Symbol ::= Digraph | Character
Digraph ::= Escaper Delimiter
Character ::= #x0-#x5a | #x5c | #x5e-#x5f | #x61-#x10ffff
`
var Opener = "["
var Closer = "]"
var Delimiter = "`"

type Jevko struct {
	Branches []Subjevko
	Suffix   string
}

type Subjevko struct {
	Prefix string
	Branch Jevko
}

const (
	Prefix = iota
	InJev
	Suffix
)

func ParseJevko(data string, escaped bool) (int, *Jevko, error) {
	//Parse subjevkos until we can't then put the rest in the suffix
	jev := Jevko{}
	length := len(data)
	ii := 0
	ps := 0
	for ii < length {
		nparsed, sj, err := ParseSubJevko(data[ii:], escaped)
		if err != nil {
			return ii + nparsed, &jev, err
		}
		if nparsed > 0 {
			jev.Branches = append(jev.Branches, *sj)
			ii += nparsed
			ps = nparsed
		} else if string(data[ii]) == Closer && !escaped {
			jev.Suffix = data[ps:ii]
			return ii, &jev, nil
		}
		if string(data[ii]) == Delimiter {
			escaped = !escaped
		} else {
			escaped = false
		}
		ii++
	}
	return ii, &jev, nil
}

func ParseSubJevko(data string, escaped bool) (int, *Subjevko, error) {
	// Parse into suffix until we hit an opener
	length := len(data)
	sj := Subjevko{}
	for ii := 0; ii < length; ii++ {
		cstr := string(data[ii])

		if cstr == Opener && !escaped {
			if ii > 0 && string(data[ii-1]) == Delimiter {
				continue
			}
			sj.Prefix = data[:ii]
			nparsed, jv, err := ParseJevko(data[ii+1:], escaped)
			if err != nil {
				return ii, &sj, err
			}
			sj.Branch = *jv
			if ii+nparsed+1 >= length {
				return ii + nparsed + 1, &sj, fmt.Errorf("unexpected end of input, expected closing delimiter ] ")
			}
			closer := string(data[ii+nparsed+1])
			if closer != Closer {
				return 0, nil, fmt.Errorf("expected closing delimiter ], got %s", closer)
			} else {
				return ii + nparsed + 1, &sj, nil
			}
		} else if cstr == Closer && !escaped {
			return 0, nil, nil
		}
		if cstr == Delimiter {
			escaped = !escaped
		} else {
			escaped = false
		}
	}
	return 0, &sj, nil
}
