package gomin

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

const (
	LineFeed = '\u000A'
)

var theA, theB, theLookahead, theX, theY rune
var output string
var content []byte

func Js(b []byte) (o string, err error) {
	theLookahead = utf8.RuneError
	theX = utf8.RuneError
	theY = utf8.RuneError
	content = b
	theA = LineFeed
	theB = LineFeed
	output = ""
	err = nil

	action(3)

	for theA != utf8.RuneError {
		switch theA {
		case ' ':
			if isAlphanum(theB) {
				action(1)
			} else {
				action(2)
			}
			break
		case LineFeed:
			switch theB {
			case '{', '[', '(', '+', '-', '!', '~':
				action(1)
				break
			case ' ':
				action(3)
				break
			default:
				if isAlphanum(theB) {
					action(1)
				} else {
					action(2)
				}
			}
			break
		default:
			switch theB {
			case ' ':
				if isAlphanum(theA) {
					action(1)
				} else {
					action(3)
				}
				break
			case LineFeed:
				switch theA {
				case '}', ']', ')', '+', '-', '"', '\'', '`':
					action(1)
					break
				default:
					if isAlphanum(theA) {
						action(1)
					} else {
						action(3)
					}
				}
				break
			default:
				action(1)
				break
			}
		}
	}

	o = output
	// Do we want to trim??
	// o = strings.TrimSpace(output)
	return
}

func getc() rune {
	if len(content) <= 0 {
		return utf8.RuneError
	}

	r, size := utf8.DecodeRune(content)
	content = content[size:]
	return r
}

func putc(r rune) {
	output += string(r)
}

// get -- return the next character from stdin. Watch out for lookahead. If
// the character is a control character, translate it to a space or
// linefeed.
func get() rune {
	var c rune = theLookahead
	theLookahead = utf8.RuneError
	if c == utf8.RuneError {
		c = getc()
	}
	if c >= ' ' || c == LineFeed || c == utf8.RuneError {
		return c
	}
	if c == '\r' {
		return LineFeed
	}
	return ' '
}

// peek -- get the next character without getting it.
func peek() rune {
	theLookahead = get()
	return theLookahead
}

// next -- get the next character, excluding comments. peek() is used to see
// if a '/' is followed by a '/' or '*'.
func next() rune {
	c := get()
	if c == '/' {
		switch peek() {
		case '/':
			for {
				c = get()
				if c <= LineFeed || c == utf8.RuneError {
					break
				}
			}
			break
		case '*':
			get()
			for c != ' ' {
				switch get() {
				case '*':
					if peek() == '/' {
						get()
						c = ' '
					}
					break
				case utf8.RuneError:
					panic("Unterminated comment.")
				}
			}
			break
		}
	}

	theY = theX
	theX = c
	return c
}

func action(d int) {
	switch d {
	case 1:
		putc(theA)
		if (theY == LineFeed || theY == ' ') &&
			(theA == '+' || theA == '-' || theA == '*' || theA == '/') &&
			(theB == '+' || theB == '-' || theB == '*' || theB == '/') {
			putc(theY)
		}
		action(2)
	case 2:
		theA = theB
		if theA == '\'' || theA == '"' || theA == '`' {
			for {
				putc(theA)
				theA = get()
				if theA == theB {
					break
				}
				if theA == '\\' {
					putc(theA)
					theA = get()
				}
				if theA == utf8.RuneError {
					panic("Unterminated string literal.")
				}
			}
		}
		action(3)
	case 3:
		theB = next()
		if theB == '/' && (theA == '(' || theA == ',' || theA == '=' || theA == ':' ||
			theA == '[' || theA == '!' || theA == '&' || theA == '|' ||
			theA == '?' || theA == '+' || theA == '-' || theA == '~' ||
			theA == '*' || theA == '/' || theA == '{' || theA == LineFeed) {
			putc(theA)
			if theA == '/' || theA == '*' {
				putc(' ')
			}
			putc(theB)
			for {
				theA = get()
				if theA == '[' {
					for {
						putc(theA)
						theA = get()
						if theA == ']' {
							break
						}
						if theA == '\\' {
							putc(theA)
							theA = get()
						}
						if theA == utf8.RuneError {
							fmt.Println("Unterminated set in Regular Expression literal.")
						}
					}
				} else if theA == '/' {
					switch peek() {
					case '/', '*':
						fmt.Println("Unterminated set in Regular Expression literal.")
					}
					break
				} else if theA == '\\' {
					putc(theA)
					theA = get()
				}
				if theA == utf8.RuneError {
					fmt.Println("Unterminated Regular Expression literal.")
				}
				putc(theA)
			}
			theB = next()
		}
	}
}

func isAlphanum(r rune) bool {
	return unicode.IsLetter(r) ||
		unicode.IsDigit(r) ||
		r == '_' ||
		r == '$' ||
		r == '\\'
}
