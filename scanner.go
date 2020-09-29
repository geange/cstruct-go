package cstruct

import (
	"errors"
	"io"
	"strconv"
)

type Scanner interface {
	Format() error
	Fetch() (byte, error)
	Next() (byte, error)
	Index(n int) (byte, error)
	CurrentIndex() (int, error)
	ScanDigit() (Token, error)
	ScanLetter() (Token, error)
	Scan() ([]Token, error)
}

func NewScanner(text []byte) Scanner {
	return &scanner{
		bs:    text,
		index: 0,
	}
}

func NewScannerString(text string) Scanner {
	return &scanner{
		bs:    []byte(text),
		index: 0,
	}
}

type scanner struct {
	bs    []byte
	index int
}

func (s *scanner) Format() error {
	panic("implement me")
}

func (s *scanner) Fetch() (byte, error) {
	if s.index >= len(s.bs) {
		return 0, io.EOF
	}
	return s.bs[s.index], nil
}

func (s *scanner) Next() (byte, error) {
	if s.index >= len(s.bs) {
		return 0, io.EOF
	}

	r := s.bs[s.index]
	s.index++
	return r, nil
}

func (s *scanner) Index(n int) (byte, error) {
	if n > len(s.bs) {
		return 0, io.EOF
	}
	return s.bs[n], nil
}

func (s *scanner) CurrentIndex() (int, error) {
	return s.index, nil
}

func (s *scanner) ScanDigit() (Token, error) {
	buf := make([]byte, 0, 10)

	numFlag := false
	pointFlag := false

	for {
		word, err := s.Fetch()
		if err != nil {
			return nil, err
		}
		if isDigit(word) {
			if !numFlag {
				numFlag = true
			}
			buf = append(buf, word)
			s.Next()
			continue
		}

		if word == '.' {
			if !pointFlag {
				pointFlag = true
				buf = append(buf, word)
				s.Next()
				continue
			} else {
				return nil, errors.New("float number has not only one point")
			}
		}

		// end of digit
		break
	}

	if pointFlag {
		num, err := strconv.ParseFloat(string(buf), 64)
		if err != nil {
			return nil, err
		}
		return &DigitToken{
			ttype:      TFloatingPoint,
			floatValue: num,
		}, nil
	}

	num, err := strconv.Atoi(string(buf))
	if err != nil {
		return nil, err
	}
	return &DigitToken{
		ttype:    TInteger,
		intValue: num,
	}, nil
}

func (s *scanner) ScanLetter() (Token, error) {
	buf := make([]byte, 0, 5)

	letterFlag := false

	for {
		b, err := s.Fetch()
		if err != nil {
			return nil, err
		}

		if isLetter(b) {
			letterFlag = true
			buf = append(buf, b)
			s.Next()
			continue
		}

		if isDigit(b) {
			if !letterFlag {
				return nil, errors.New("number not allow to be use as the first of the word")
			}
			buf = append(buf, b)
			s.Next()
			continue
		}
		// end of digit
		break
	}

	word := string(buf)
	if IN(word, ReserveWords...) {
		switch word {
		case "s64":
			return &ReservedWordToken{ttype: TInt64}, nil
		case "s32":
			return &ReservedWordToken{ttype: TInt32}, nil
		case "s16":
			return &ReservedWordToken{ttype: TInt16}, nil
		case "s8":
			return &ReservedWordToken{ttype: TInt8}, nil
		case "u64":
			return &ReservedWordToken{ttype: TUint64}, nil
		case "u32":
			return &ReservedWordToken{ttype: TUint32}, nil
		case "u16":
			return &ReservedWordToken{ttype: TUint16}, nil
		case "u8":
			return &ReservedWordToken{ttype: TUint8}, nil
		case "float":
			return &ReservedWordToken{ttype: TFloat32}, nil
		case "double":
			return &ReservedWordToken{ttype: TFloat64}, nil
		case "type":
			return &ReservedWordToken{ttype: TType}, nil
		case "typedef":
			return &ReservedWordToken{ttype: TTypedef}, nil
		case "struct":
			return &ReservedWordToken{ttype: TStruct}, nil
		case "byte":
			return &ReservedWordToken{ttype: TByte}, nil
		}
	}
	return &LetterToken{ttype: TString, value: word}, nil
}

func (s *scanner) Scan() ([]Token, error) {
	tokens := make([]Token, 0, 10)
	for {
		b, err := s.Fetch()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		switch {
		case isLetter(b):
			token, err := s.ScanLetter()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
		case isDigit(b):
			token, err := s.ScanDigit()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
		case b == ' ' || b == '\n' || b == '\t':
			s.Next()
			continue
		case b == ';':
			s.Next()
			token := &DelimiterToken{ttype: TLineEnd}
			tokens = append(tokens, token)
		case b == '[':
			s.Next()
			token := &DelimiterToken{ttype: TLeftBracket}
			tokens = append(tokens, token)
		case b == ']':
			s.Next()
			token := &DelimiterToken{ttype: TRightBracket}
			tokens = append(tokens, token)
		case b == '{':
			s.Next()
			token := &DelimiterToken{ttype: TLeftBrace}
			tokens = append(tokens, token)
		case b == '}':
			s.Next()
			token := &DelimiterToken{ttype: TRightBrace}
			tokens = append(tokens, token)
		}
	}
	return tokens, nil
}
