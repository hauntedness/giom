package runtime

import (
	"runtime"
	"strconv"
	"strings"
)

func Source() string {
	_, file, line, _ := runtime.Caller(1)
	buf := &strings.Builder{}
	writeTrimmedPath(buf, file)
	buf.WriteByte(':')
	buf.WriteString(strconv.Itoa(line))
	return buf.String()
}

func SourceWithSkip(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	buf := &strings.Builder{}
	writeTrimmedPath(buf, file)
	buf.WriteByte(':')
	buf.WriteString(strconv.Itoa(line))
	return buf.String()
}

func writeTrimmedPath(buf *strings.Builder, path string) {
	idx := strings.LastIndexByte(path, '/')
	if idx == -1 {
		buf.WriteString(path)
		return
	}
	// Find the penultimate separator.
	idx = strings.LastIndexByte(path[:idx], '/')
	if idx == -1 {
		buf.WriteString(path)
		return
	}
	// Find the penultimate separator.
	idx2 := strings.LastIndexByte(path[:idx], '/')
	if idx == -1 {
		buf.WriteString(path[idx+1:])
		return
	}
	buf.WriteString(path[idx2+1:])
}
