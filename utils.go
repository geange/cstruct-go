package cstruct

var (
	defaultFormat = "%s_%s"
)

func Format(format string) error {
	defaultFormat = format
	return nil
}

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

func ToUnderLine(src string) string {
	result := make([]byte, 0)
	for i := 0; i < len(src); i++ {
		w := src[i]
		if i == 0 && isUpper(w) {
			result = append(result, toLower(w))
			continue
		}

		if isUpper(w) {
			if src[i-1] == '_' {
				result = append(result, toLower(w))
				continue
			}
			result = append(result, '_', toLower(w))
			continue
		}

		result = append(result, w)
	}
	return string(result)
}

func isUpper(w byte) bool {
	return w >= 'A' && w <= 'Z'
}

func toLower(w byte) byte {
	return w - 'A' + 'a'
}
