package cstruct

import (
	"bytes"
	"fmt"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestLexer_Structure(t *testing.T) {
	buf := bytes.NewBufferString(text)
	lex, err := NewLexer(buf)
	assert.Equal(t, err, nil)
	cStruct, err := lex.Structure()
	assert.Equal(t, err, nil)

	fmt.Println(cStruct.Name)

	for _, v := range cStruct.FieldSet {
		fmt.Println(v)
	}
}

func TestLexer_AllStructure(t *testing.T) {
	buf := bytes.NewBufferString(nestedText)
	lex, err := NewLexer(buf)
	assert.Equal(t, err, nil)
	list, err := lex.AllStructure()
	assert.Equal(t, err, nil)
	for _, item := range list {
		fmt.Println(item.Name)
		for _, v := range item.FieldSet {
			fmt.Println(v)
		}
	}
}
