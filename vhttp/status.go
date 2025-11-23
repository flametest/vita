package vhttp

type StatusCode int

func (s StatusCode) Int() int {
	return int(s)
}
