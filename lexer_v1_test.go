package cstruct

import (
	"bytes"
	"fmt"
	"github.com/go-playground/assert/v2"
	"testing"
)

func Test_lexer_single_structure(t *testing.T) {
	buf := bytes.NewBufferString(text)
	lex, err := NewLexerV1(buf)
	assert.Equal(t, err, nil)
	cStruct, err := lex.singleStructure()
	assert.Equal(t, err, nil)

	fmt.Println(cStruct.FieldName())

	fields, err := cStruct.Flat()
	assert.Equal(t, err, nil)

	for _, field := range fields {
		fmt.Println(field.FieldName(), field.Type(), field.Index())
	}
}
