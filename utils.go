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
