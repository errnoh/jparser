package jparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

var Verbose bool

func Parse(r io.Reader) (root *Token, err error) {
	// TODO: error when unsuccessfull
	var (
		current *Token
		tkn     *Token
		key     string
		t       json.Token
	)

	if Verbose {
		var buf bytes.Buffer
		r = io.TeeReader(r, &buf)
		defer func() {
			fmt.Println(buf.String())
		}()
	}

	dec := json.NewDecoder(r)
	dec.UseNumber() // Prevent converting integers into floats

loop:
	for {
		tkn = nil
		t, err = dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch x := t.(type) {
		case json.Delim:
			switch x {
			case '{':
				tkn = &Token{
					tokenType: OBJECT,
					parent:    current,
					object:    make(map[string]*Token),
				}
			case '}':
				current = current.parent
				continue
			case '[':
				tkn = &Token{
					tokenType: ARRAY,
					parent:    current,
					array:     make([]*Token, 0),
				}
			case ']':
				if len(current.array) > 0 {
					for _, v := range current.array {
						if v.tokenType != current.array[0].tokenType {
							current.arrayType = MIXED
							break
						}
						current.arrayType = current.array[0].tokenType
					}
				}
				current = current.parent
				continue
			}
		case string:
			if current.tokenType == OBJECT && key == "" {
				key = x
				continue loop
			}
			tkn = &Token{
				tokenType: STRING,
				parent:    current,
				value:     x,
			}
		case json.Number:
			tkn = &Token{
				tokenType: NUMBER,
				parent:    current,
				value:     x,
			}
		case float64:
			tkn = &Token{
				tokenType: NUMBER,
				parent:    current,
				value:     x,
			}
		case bool:
			tkn = &Token{
				tokenType: BOOL,
				parent:    current,
				value:     x,
			}
		case nil:
			tkn = &Token{
				tokenType: NULL,
				parent:    current,
			}
		}
		if root == nil {
			root = tkn
			current = tkn
			continue
		}
		if current.tokenType == ARRAY {
			current.array = append(current.array, tkn)
		} else if current.tokenType == OBJECT {
			current.object[key] = tkn
		}
		if current != nil && tkn != nil {
			if current.tokenType == ARRAY {
				tkn.key = strconv.Itoa(len(current.array) - 1)
				//				fmt.Printf("1: %#v\n2: %#v\n", current, tkn)
			}
			if current.tokenType == OBJECT {
				tkn.key = key
				//				fmt.Printf("1: %#v\n2: %#v\n", current, tkn)
			}
		}
		if tkn.tokenType == ARRAY || tkn.tokenType == OBJECT {
			current = tkn
		}
		if key != "" {
			key = ""
		}
	}
	return root, nil
}
