package log

func ln(s string) string {
	if n := len(s); n == 0 || s[n-1] != '\n' {
		s += "\n"
	}
	return s
}
