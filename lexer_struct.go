package cstruct

import "github.com/pkg/errors"

func (l *lexer) structureInBrace() ([]Field, error) {
	token, err := l.Next()
	if err != nil {
		return nil, err
	}
	if token.Type() != TLeftBrace {
		return nil, errors.New("'{' not found")
	}

	result := make([]Field, 0)

	for {
		fs, err := l.Statement()
		if err != nil {
			return nil, err
		}
		result = append(result, fs...)

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
