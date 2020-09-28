package cstruct

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
)

type Lexer interface {
	Fetch() (Token, error)
	Next() (Token, error)
	Index(n int) (Token, error)
	CurrentIndex() int
	Statement() ([]Field, error)
	Structure() (*CStruct, error)
	AllStructure() ([]*CStruct, error)
	ExistStructure(name string) (*CStruct, bool)
}

func NewLexer(reader io.Reader) (Lexer, error) {
	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	scanner := NewScanner(bs)
	tokens, err := scanner.Scan()
	if err != nil {
		return nil, err
	}

	return &lexer{
		tokens:     tokens,
		index:      0,
		structures: make(map[string]CStruct),
	}, nil
}

type lexer struct {
	tokens     []Token
	index      int
	structures map[string]CStruct
}

func (l *lexer) Fetch() (Token, error) {
	if l.index >= len(l.tokens) {
		return nil, io.EOF
	}

	token := l.tokens[l.index]
	return token, nil
}

func (l *lexer) Next() (Token, error) {
	if l.index >= len(l.tokens) {
		return nil, io.EOF
	}
	token := l.tokens[l.index]
	l.index++
	return token, nil
}

func (l *lexer) Index(n int) (Token, error) {
	if n >= len(l.tokens) {
		return nil, io.EOF
	}

	token := l.tokens[n]
	return token, nil
}

func (l *lexer) CurrentIndex() int {
	return l.index
}

func (l *lexer) Statement() ([]Field, error) {
	result := make([]Field, 0)

	dataTypeToken, err := l.Next()
	if err != nil {
		return nil, err
	}

	nameToken, err := l.Next()
	if err != nil {
		return nil, err
	}

	tToken, err := l.Fetch()
	if err != nil {
		return nil, err
	}

	var dType FieldType
	switch dataTypeToken.Type() {
	case TInt64:
		dType = Int64
	case TInt32:
		dType = Int32
	case TInt16:
		dType = Int16
	case TInt8:
		dType = Int8
	case TUint64:
		dType = UInt64
	case TUint32:
		dType = UInt32
	case TUint16:
		dType = UInt16
	case TUint8:
		dType = UInt8
	case TFloat32:
		dType = Float32
	case TFloat64:
		dType = Float64
	case TChar:
		dType = Hex
	}

	// 基础类型
	if tToken.Type() == TLineEnd {
		name := nameToken.Value().(string)
		switch dataTypeToken.Type() {
		case TString:
			// 内嵌子集
			structName := dataTypeToken.Value().(string)
			cStruct, ok := l.ExistStructure(structName)
			if !ok {
				return nil, errors.New(fmt.Sprintf("%s struct not defined", structName))
			}
			fs, err := cStruct.ToStatement()
			if err != nil {
				return nil, errors.Wrap(err, "get statement from sub struct")
			}
			for i := range fs {
				fs[i].Name = fmt.Sprintf("%s_%s", name, fs[i].Name)
			}
			result = append(result, fs...)
		default:
			field := Field{
				Name: name,
				Type: dType,
				Size: getFieldTypeSize(dType),
			}
			result = append(result, field)
		}
		_, _ = l.Next()
		return result, nil
	}

	switch dataTypeToken.Type() {
	case TFloat64, TInt64, TInt32, TInt16, TInt8, TUint64, TUint32, TUint16, TUint8, TChar:
		if tToken.Type() == TLeftBracket {
			fs, err := l.arrayStatement(dataTypeToken, nameToken, nil, false)
			if err != nil {
				return nil, err
			}
			result = append(result, fs...)
		}
	case TString:
		structName := dataTypeToken.Value().(string)
		if tToken.Type() == TLeftBracket {
			if cStruct, ok := l.ExistStructure(structName); ok {
				fs, err := l.arrayStatement(dataTypeToken, nameToken, cStruct, true)
				if err != nil {
					return nil, err
				}
				result = append(result, fs...)
			}
		}
		return nil, errors.New(fmt.Sprintf("%s struct not defined", structName))
	}
	return result, nil
}

func (l *lexer) Structure() (*CStruct, error) {
	token, err := l.Next()
	if err != nil {
		return nil, err
	}

	switch token.Type() {
	case TTypedef:
		stToken, err := l.Next()
		if err != nil {
			return nil, err
		}

		if stToken.Type() != TStruct {
			return nil, errors.New("'struct' not found after typedef")
		}

		bToken, err := l.Fetch()
		if err != nil {
			return nil, err
		}

		if bToken.Type() != TLeftBrace {
			return nil, errors.New("'{' not found after typedef struct")
		}

		fs, err := l.structureInBrace()
		if err != nil {
			return nil, err
		}

		nameToken, err := l.Next()
		if err != nil {
			return nil, err
		}

		if nameToken.Type() != TString {
			return nil, errors.New("structure type not found")
		}

		name := nameToken.Value().(string)
		cStruct := CStruct{
			Name:     name,
			FieldSet: FieldSet(fs),
		}

		endToken, err := l.Next()
		if err != nil {
			return nil, errors.Wrap(err, "';' bot found after struct defined")
		}
		if endToken.Type() != TLineEnd {
			return nil, errors.Wrap(err, "';' bot found after struct defined")
		}

		return &cStruct, nil
	case TStruct:
		nameToken, err := l.Next()
		if err != nil {
			return nil, err
		}

		if nameToken.Type() != TString {
			return nil, errors.New("structure type not found")
		}

		fs, err := l.structureInBrace()
		if err != nil {
			return nil, err
		}
		name := nameToken.Value().(string)
		cStruct := CStruct{
			Name:     name,
			FieldSet: FieldSet(fs),
		}

		endToken, err := l.Next()
		if err != nil {
			return nil, errors.Wrap(err, "';' bot found after struct defined")
		}
		if endToken.Type() != TLineEnd {
			return nil, errors.Wrap(err, "';' bot found after struct defined")
		}

		return &cStruct, nil
	}
	return nil, errors.New("struct format error")
}

func (l *lexer) AllStructure() ([]*CStruct, error) {
	result := make([]*CStruct, 0)
	for {
		cStruct, err := l.Structure()
		if err != nil {
			return nil, err
		}
		result = append(result, cStruct)
		l.structures[cStruct.Name] = *cStruct

		_, err = l.Fetch()
		if err == io.EOF {
			break
		}
	}
	return result, nil
}

func (l *lexer) ExistStructure(name string) (*CStruct, bool) {
	c, ok := l.structures[name]
	return &c, ok
}
