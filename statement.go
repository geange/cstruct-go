package cstruct

import (
	"errors"
	"fmt"
)

type Statement struct {
}

func (l *lexer) arrayStatement(dataType, nameToken Token, structure *CStruct, isStruct bool) ([]Field, error) {
	bToken, err := l.Next()
	if err != nil {
		return nil, err
	}
	if bToken.Type() != TLeftBracket {
		return nil, errors.New(fmt.Sprintf("'[' not found in array check, token type is %s", TokenTypeMap[bToken.Type()]))
	}

	result := make([]Field, 0)
	token, err := l.Next()
	if err != nil {
		return nil, err
	}

	switch token.Type() {
	case TInteger:
		n := token.Value().(int)
		token, err = l.Next()
		if err != nil {
			return nil, err
		}
		if token.Type() != TRightBracket {
			return nil, errors.New("right bracket not found")
		}

		subName := nameToken.Value().(string)

		if isStruct {
			fs, err := structure.ToStatement()
			if err != nil {
				return nil, err
			}
			for i := range fs {
				fs[i].Name = subName + "_" + fs[i].Name
			}

			for i := 0; i < n; i++ {
				for index := range fs {
					fs[index].Name = fmt.Sprintf("%s_%d", fs[index].Name, i)
					result = append(result, fs[index])
				}
			}
			return result, nil
		}

		var dType FieldType
		switch dataType.Type() {
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
		}

		for i := 0; i < n; i++ {
			name := fmt.Sprintf("%s_%d", subName, i)
			field := Field{
				Name:        name,
				DisplayName: name,
				Display:     true,
				Type:        dType,
			}
			result = append(result, field)
		}

		token, err := l.Next()
		if err != nil {
			return nil, err
		}
		if token.Type() != TLineEnd {
			return nil, errors.New("semicolon not found")
		}
		return result, nil
	default:
		return nil, errors.New("unsupported type in array")
	}
}
