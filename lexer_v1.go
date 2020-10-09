package cstruct

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
)

type LexerV1 interface {
	Fetch() (Token, error)
	Next() (Token, error)
	Index(n int) (Token, error)
	CurrentIndex() int
	AllStruct() ([]BaseStructV1, error)
}

func (l *lexer) AllStruct() ([]BaseStructV1, error) {
	r := make([]BaseStructV1, 0)
	for {
		cStruct, err := l.scanStruct()
		if err != nil {
			return nil, err
		}
		r = append(r, *cStruct)
		l.cache[cStruct.sName] = *cStruct

		_, err = l.Fetch()
		if err == io.EOF {
			break
		}
	}
	return r, nil
}

func (l *lexer) scanStruct() (*BaseStructV1, error) {
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

		fs, err := l.scanStructureInBrace()
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
		obj := BaseStructV1{
			sName:  name,
			fields: fs,
		}

		endToken, err := l.Next()
		if err != nil {
			return nil, errors.Wrap(err, "';' bot found after struct defined")
		}
		if endToken.Type() != TLineEnd {
			return nil, errors.Wrap(err, "';' bot found after struct defined")
		}

		return &obj, nil
	case TStruct:
		nameToken, err := l.Next()
		if err != nil {
			return nil, err
		}

		if nameToken.Type() != TString {
			return nil, errors.New("structure type not found")
		}

		fs, err := l.scanStructureInBrace()
		if err != nil {
			return nil, err
		}
		name := nameToken.Value().(string)
		obj := BaseStructV1{
			sName:  name,
			fields: fs,
		}

		endToken, err := l.Next()
		if err != nil {
			return nil, errors.Wrap(err, "';' bot found after struct defined")
		}
		if endToken.Type() != TLineEnd {
			return nil, errors.Wrap(err, "';' bot found after struct defined")
		}

		return &obj, nil
	}
	return nil, errors.New("struct format error")
}

func (l *lexer) scanStructureInBrace() ([]StatementV1, error) {
	token, err := l.Next()
	if err != nil {
		return nil, err
	}
	if token.Type() != TLeftBrace {
		return nil, errors.New("'{' not found")
	}

	result := make([]StatementV1, 0)

	for {
		state, err := l.scanStatement()
		if err != nil {
			return nil, err
		}
		result = append(result, state)

		token, err := l.Fetch()
		if err != nil {
			return nil, err
		}
		if token.Type() == TRightBrace {
			_, _ = l.Next()
			break
		}
	}

	return result, nil
}

// 读取声明语句
func (l *lexer) scanStatement() (StatementV1, error) {
	dtToken, err := l.Next()
	if err != nil {
		return nil, err
	}

	nameToken, err := l.Next()
	if err != nil {
		return nil, err
	}
	name := nameToken.Value().(string)

	nToken, err := l.Fetch()
	if err != nil {
		return nil, err
	}

	var dType FieldType
	switch dtToken.Type() {
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
	case TByte:
		dType = BYTE
	}

	switch nToken.Type() {
	case TLineEnd:
		// 读取完后面的';'
		_, _ = l.Next()
		base := BaseV1{
			name:  name,
			fType: dType,
			size:  fieldTypeSize(dType),
		}
		switch dtToken.Type() {
		case TByte:
			return &Byte{base}, nil
		case TUint8:
			return &U8{base}, nil
		case TInt8:
			return &S8{base}, nil
		case TUint16:
			return &U16{base}, nil
		case TInt16:
			return &S16{base}, nil
		case TUint32:
			return &U32{base}, nil
		case TInt32:
			return &S32{base}, nil
		case TUint64:
			return &U64{base}, nil
		case TInt64:
			return &S64{base}, nil
		case TFloat32:
			return &F32{base}, nil
		case TFloat64:
			return &F64{base}, nil
		case TString:
			className := dtToken.Value().(string)
			if item, ok := l.cache[className]; ok {
				v := item
				v.fName = name

				fmt.Println("~~~~~~~~", v)

				return &v, nil
			}
			return nil, fmt.Errorf("struct %s not exist", className)
		}

	case TLeftBracket:
		_, _ = l.Next()
		sizeToken, err := l.Next()
		if err != nil {
			return nil, err
		}
		size, ok := sizeToken.Value().(int)
		if !ok {
			return nil, errors.New("array size is not a number")
		}
		rbToken, err := l.Next()
		if err != nil {
			return nil, err
		}
		if rbToken.Type() != TRightBracket {
			return nil, fmt.Errorf("']' not found after %s", name)
		}
		endToken, err := l.Next()
		if err != nil {
			return nil, err
		}
		if endToken.Type() != TLineEnd {
			return nil, fmt.Errorf("';' not found after %s", name)
		}

		base := BaseV1{
			fType: dType,
			size:  fieldTypeSize(dType),
		}
		array := BaseArrayV1{
			name: name,
			size: size,
		}
		switch dtToken.Type() {
		case TByte:
			array.fType = &Byte{base}
		case TUint8:
			array.fType = &U8{base}
		case TInt8:
			array.fType = &S8{base}
		case TUint16:
			array.fType = &U16{base}
		case TInt16:
			array.fType = &S16{base}
		case TUint32:
			array.fType = &U32{base}
		case TInt32:
			array.fType = &S32{base}
		case TUint64:
			array.fType = &U64{base}
		case TInt64:
			array.fType = &S64{base}
		case TFloat32:
			array.fType = &F32{base}
		case TFloat64:
			array.fType = &F64{base}
		case TString:
			className := dtToken.Value().(string)
			if item, ok := l.cache[className]; ok {
				v := item
				v.fName = name
				return &BaseArrayV1{
					name:  name,
					fType: &v,
					size:  size,
				}, nil
			}
			return nil, fmt.Errorf("struct %s not exist", className)
		default:
			return nil, fmt.Errorf("unsupported array type")
		}
		return &array, nil
	}

	return nil, errors.New("unknown error")

}
