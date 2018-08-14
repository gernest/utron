package misc

import "strings"

// StripCtlAndExtFromUTF8 removes all Control and Special characters from the given string
func StripCtlAndExtFromUTF8(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r < 127 {
			return r
		}
		return -1
	}, str)
}
