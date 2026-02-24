package lineprotocol

import (
	"strconv"
	"strings"
)

// escaper represents a set of characters that can be escaped.
type escaper struct {
	// table maps from byte value to the byte used to escape that value.
	// If an entry is zero, it doesn't need to be escaped.
	table [256]byte

	// revTable holds the inverse of table - it maps
	// from escaped value to the unescaped value.
	revTable [256]byte

	// escapes holds all the characters that need to be escaped.
	escapes string
}

// newEscaper returns an escaper that escapes the
// given characters.
func newEscaper(escapes string) *escaper {
	var e escaper
	for _, b := range escapes {
		// Note that this works because the Go escaping rules
		// for white space are the same as line-protocol's.
		q := strconv.QuoteRune(b)
		q = q[1 : len(q)-1]             // strip single quotes.
		q = strings.TrimPrefix(q, "\\") // remove backslash if present.
		e.table[byte(b)] = q[0]         // use single remaining character.
		e.revTable[q[0]] = byte(b)
	}
	e.escapes = escapes
	return &e
}

// appendEscaped returns the escaped form of s appended to buf.
func (e *escaper) appendEscaped(buf []byte, s string) []byte {
	// Fast path: find first character that needs escaping.
	first := -1
	for i := 0; i < len(s); i++ {
		if e.table[s[i]] != 0 {
			first = i
			break
		}
	}
	if first == -1 {
		return append(buf, s...)
	}
	// Copy prefix that needs no escaping, then single-pass escape the rest.
	buf = append(buf, s[:first]...)
	for i := first; i < len(s); i++ {
		if r := e.table[s[i]]; r != 0 {
			buf = append(buf, '\\', r)
		} else {
			buf = append(buf, s[i])
		}
	}
	return buf
}

// escaped returns the length that s will be after escaping
// and the index of the first character in s that needs escaping.
func (e *escaper) escapedLen(s string) (escLen, startIndex int) {
	startIndex = len(s)
	n := len(s)
	for i := 0; i < len(e.escapes); i++ {
		k := strings.IndexByte(s, e.escapes[i])
		if k == -1 {
			continue
		}
		if k < startIndex {
			startIndex = k
		}
		n += 1 + strings.Count(s[k+1:], e.escapes[i:i+1])
	}
	return n, startIndex
}
