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
	AllStructureV1() ([]CStructureV1, error)
	ExistStructureV1(name string) (CStructureV1, bool)

	singleStructure() (*CStructV1, error)
	statementV1() (CStructureV1, error)
	structureInBraceV1() ([]CStructureV1, error)
}

func (l *lexer) singleStructure() (*CStructV1, error) {
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

		fs, err := l.structureInBraceV1()
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
		cStruct := CStructV1{
			gName:  name,
			fields: fs,
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

		fs, err := l.structureInBraceV1()
		if err != nil {
			return nil, err
		}
		name := nameToken.Value().(string)
		cStruct := CStructV1{
			gName:  name,
			fields: fs,
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

func (l *lexer) statementV1() (CStructureV1, error) {
	dataTypeToken, err := l.Next()
	if err != nil {
		return nil, err
	}

	// 获取字段名
	nameToken, err := l.Next()
	if err != nil {
		return nil, err
	}
	name := nameToken.Value().(string)

	tToken, err := l.Fetch()
	if err != nil {
		return nil, err
	}

	isStruct := false

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
	case TByte:
		dType = Hex
	default:
		isStruct = true
	}

	dataTypeName := ""
	if isStruct {
		dataTypeName = dataTypeToken.Value().(string)
	}

	// 非数组
	if tToken.Type() == TLineEnd {
		// 获取终止的分号
		_, err = l.Next()
		if err != nil {
			return nil, err
		}

		// 结构体
		if isStruct {
			// 内嵌子集
			item, ok := l.ExistStructureV1(dataTypeName)
			if !ok {
				return nil, errors.New(fmt.Sprintf("struct %s not defined", dataTypeName))
			}
			item.SetFieldName(name)
			return item, nil
		}

		// 普通类型
		return &BaseStruct{
			name:  name,
			fType: dType,
			size:  getFieldTypeSize(dType),
		}, nil
	}

	if tToken.Type() != TLeftBracket {
		return nil, errors.Errorf("token error, value[%v]", tToken.Value())
	}

	_, _ = l.Next()
	// 获取数组大小
	sizeToken, err := l.Next()
	if err != nil {
		return nil, err
	}
	size, ok := sizeToken.Value().(int)
	if !ok {
		return nil, errors.Errorf("size value is error, value[%v]", sizeToken.Value())
	}

	// 检查是否有右括号
	rToken, err := l.Next()
	if err != nil {
		return nil, err
	}
	if rToken.Type() != TRightBracket {
		return nil, errors.New("right bracket not found")
	}
	endToken, err := l.Next()
	if err != nil {
		return nil, err
	}
	if endToken.Type() != TLineEnd {
		return nil, errors.New("';' not found")
	}

	// 结构体数组
	if isStruct {
		item, ok := l.ExistStructureV1(dataTypeName)
		if !ok {
			return nil, errors.New(fmt.Sprintf("struct %s not defined", dataTypeName))
		}

		return &ArrayV1{
			name:  name,
			size:  uint(size),
			field: item,
		}, nil
	}

	// 普通类型数组
	array := ArrayV1{
		name: name,
		size: uint(size),
		field: &BaseStruct{
			name:  name,
			fType: dType,
			size:  getFieldTypeSize(dType),
		},
	}
	return &array, nil
}

func (l *lexer) structureInBraceV1() ([]CStructureV1, error) {
	token, err := l.Next()
	if err != nil {
		return nil, err
	}
	if token.Type() != TLeftBrace {
		return nil, errors.New("'{' not found")
	}

	result := make([]CStructureV1, 0)

	for {
		state, err := l.statementV1()
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

func (l *lexer) AllStructureV1() ([]CStructureV1, error) {
	result := make([]CStructureV1, 0)
	for {
		cStruct, err := l.singleStructure()
		if err != nil {
			return nil, err
		}
		result = append(result, cStruct)
		l.structMap[cStruct.gName] = *cStruct

		_, err = l.Fetch()
		if err == io.EOF {
			break
		}
	}
	return result, nil
}

func (l *lexer) ExistStructureV1(name string) (CStructureV1, bool) {
	item, ok := l.structMap[name]
	if ok {
		r := item.Copy()
		return r, ok
	}
	return nil, false
}
