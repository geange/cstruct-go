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

		// 内嵌结构体
		if isStruct {
			fs, err := structure.ToStatement()
			if err != nil {
				return nil, err
			}
			for index := range fs {
				fs[index].Name = fmt.Sprintf("%s_%s", subName, fs[index].Name)
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
		case TByte:
			dType = Hex
		}

		if dType == Hex {
			field := Field{
				Name: subName,
				Type: Hex,
				Size: uint(n),
			}
			result = append(result, field)
		} else {
			for i := 0; i < n; i++ {
				name := fmt.Sprintf("%s_%d", subName, i)
				field := Field{
					Name: name,
					Type: dType,
					Size: getFieldTypeSize(dType),
				}
				result = append(result, field)
			}
		}

		endToken, err := l.Next()
		if err != nil {
			return nil, err
		}
		if endToken.Type() != TLineEnd {
			return nil, errors.New("';' not found")
		}
		return result, nil
	default:
		return nil, errors.New("unsupported type in array")
	}
}
