package jparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type TokenType int

func (t TokenType) String() string {
	if int(t) >= 0 && int(t) < len(typestr) {
		return typestr[int(t)]
	}
	return "?"
}

const (
	NULL = TokenType(iota)
	STRING
	NUMBER
	BOOL
	OBJECT
	ARRAY
	MIXED
)

var typestr = []string{"null", "string", "number", "bool", "object", "array", "mixed"}

type Token struct {
	parent    *Token
	tokenType TokenType
	arrayType TokenType
	array     []*Token
	object    map[string]*Token
	value     interface{}
	key       string
}

func (t *Token) String() string {
	var b bytes.Buffer
	if t.array != nil {
		fmt.Fprintf(&b, "(%v)%+v", t.arrayType, t.array)
	} else if t.object != nil {
		fmt.Fprint(&b, "{")
		for k, v := range t.object {
			fmt.Fprintf(&b, " %s:%v", k, v)
		}
		fmt.Fprint(&b, " }")
	} else {
		fmt.Fprint(&b, t.tokenType)
	}
	return b.String()
}

func (t *Token) StringValue() string {
	var b bytes.Buffer
	if t.array != nil {
		fmt.Fprintf(&b, "(%v)[", t.arrayType)
		for _, v := range t.array {
			fmt.Fprintf(&b, "%s ", v.StringValue())
		}
		fmt.Fprint(&b, "]")
	} else if t.object != nil {
		fmt.Fprint(&b, "{")
		for k, v := range t.object {
			fmt.Fprintf(&b, " %s:%v", k, v.StringValue())
		}
		fmt.Fprint(&b, " }")
	} else {
		fmt.Fprint(&b, t.value)
	}
	return b.String()
}

func (t *Token) Get(path string) (str string, ok bool) {
	var (
		current *Token = t
	)

	parts := strings.Split(path, ".")
	for _, part := range parts {
		if path == "." {
			break
		}
		if len(part) > 2 && part[0] == '[' && part[len(part)-1] == ']' {
			if current.tokenType != ARRAY {
				return "", false
			}

			index, err := strconv.Atoi(part[1 : len(part)-1])
			if err != nil {
				return "", false
			}

			if current.array == nil || index < 0 || len(current.array) <= index {
				return "", false
			}
			current = current.array[index]
			continue
		}
		if current.tokenType != OBJECT {
			return "", false
		}
		current, ok = current.object[part]
		if !ok {
			return "", false
		}
	}

	switch val := current.value.(type) {
	case string:
		str = val
		ok = true
	case json.Number:
		str = string(val)
		ok = true
	case float64:
		str = fmt.Sprintf("%f", val)
		ok = true
	case bool:
		str = fmt.Sprint(val)
		ok = true
	case nil:
		str = ""
		ok = true
	default:
		return "", false
	}
	return
}
