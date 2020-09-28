package cstruct

var (
	ReserveWords = []string{
		"auto", "break", "case", "char", "const", "continue",
		"default", "do", "double", "else", "enum", "extern",
		"float", "for", "goto", "if", "int", "long",
		"s64", "s32", "s16", "s8",
		"u64", "u32", "u16", "u8",
		"register", "return", "short", "signed", "sizeof", "static",
		"struct", "switch", "typedef", "union", "unsigned", "void",
		"volatile", "while",
	}

	IntegerWords = []string{
		"s64", "s32", "s16", "s8",
		"u64", "u32", "u16", "u8",
	}

	FloatingPointWords = []string{
		"float", "double",
	}

	TokenTypeMap = map[TokenType]string{
		TFloat32:       "float32",
		TFloat64:       "float64",
		TInt64:         "int64",
		TInt32:         "int32",
		TInt16:         "int16",
		TInt8:          "int8",
		TUint64:        "uint64",
		TUint32:        "uint32",
		TUint16:        "uint16",
		TUint8:         "uint8",
		TLineEnd:       "line_end",
		TStruct:        "struct",
		TType:          "type",
		TTypedef:       "typedef",
		TLeftBrace:     "left_brace",
		TRightBrace:    "right_brace",
		TLeftBracket:   "left_bracket",
		TRightBracket:  "right_bracket",
		TInteger:       "integer",
		TFloatingPoint: "floating_point",
		TString:        "string",
	}
)
