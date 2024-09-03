package common

import "strings"

// ConcatStrings concatenates multiple strings into one string. Using strings.Builder for more efficient memory allocation.
func ConcatStrings(strs ...string) string {
	total := 0
	for i := 0; i < len(strs); i++ {
		total += len(strs[i])
	}

	sb := strings.Builder{}
	sb.Grow(total)
	for _, v := range strs {
		_, _ = sb.WriteString(v)
	}
	return sb.String()
}
