package vstring

// Deprecated: use github.com/samber/lo.Contains instead.
func StrInList(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
