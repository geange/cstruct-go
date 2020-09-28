package cstruct

type TokenType int

const (
	TFloat32 = TokenType(iota)
	TFloat64
	TInt64
	TInt32
	TInt16
	TInt8
	TUint64
	TUint32
	TUint16
	TUint8
	TChar
	TLineEnd
	TStruct
	TType
	TTypedef
	TLeftBrace
	TRightBrace
	TLeftBracket
	TRightBracket
	TInteger
	TFloatingPoint
	TString
)

//
type Token interface {
	Type() TokenType
	Value() interface{}
}

// 标识符
type LetterToken struct {
	ttype TokenType
	value string
}

func (l *LetterToken) Type() TokenType {
	return l.ttype
}

func (l *LetterToken) Value() interface{} {
	return l.value
}

// 常数
type DigitToken struct {
	ttype      TokenType
	floatValue float64
	intValue   int
}

func (d *DigitToken) Type() TokenType {
	return d.ttype
}

func (d *DigitToken) Value() interface{} {
	switch d.ttype {
	case TInteger:
		return d.intValue
	case TFloatingPoint:
		return d.floatValue
	default:
		return 0
	}
}

// 保留字
type ReservedWordToken struct {
	ttype TokenType
}

func (r *ReservedWordToken) Type() TokenType {
	return r.ttype
}

func (r *ReservedWordToken) Value() interface{} {
	return nil
}

// 界符
type DelimiterToken struct {
	ttype TokenType
}

func (d *DelimiterToken) Type() TokenType {
	return d.ttype
}

func (d *DelimiterToken) Value() interface{} {
	return nil
}

// 运算符
type OperatorToken struct {
}
