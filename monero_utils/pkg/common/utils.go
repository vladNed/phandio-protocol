package common

// Reverse returns a copy of the slice with the bytes in reverse order
func Reverse(s []byte) []byte {
	l := len(s)
	rs := make([]byte, l)
	for i := 0; i < l; i++ {
		rs[i] = s[l-i-1]
	}
	return rs
}
