package page

import "strings"

// Page is a page of lines read from a lot of lines
type Page struct {
	Lines  []string
	Number int
	Total  int
	Size   int
}

// offset start from 0
func From(text string, offset int, size int) Page {
	lines := strings.Split(text, "\n")
	if offset < 0 || len(lines) == 0 || offset > len(lines) {
		return Page{Number: -1}
	}
	// 1,2,3{offset},4,5{len} => size := 5
	if len(lines) < offset+size {
		return Page{Lines: lines[offset:], Number: offset / size, Total: len(lines), Size: size}
	}
	return Page{Lines: lines[offset : offset+size], Number: offset / size, Total: len(lines), Size: size}
}
