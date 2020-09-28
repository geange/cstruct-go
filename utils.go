package cstruct

func IN(f string, fields ...string) bool {
	for _, v := range fields {
		if f == v {
			return true
		}
	}
	return false
}

func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c == '_'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func getFieldTypeSize(t FieldType) uint {
	switch t {
	case UInt8, Int8:
		return LEN1
	case UInt16, Int16:
		return LEN2
	case UInt32, Int32:
		return LEN4
	case UInt64, Int64:
		return LEN8
	default:
		return LEN0
	}
}
