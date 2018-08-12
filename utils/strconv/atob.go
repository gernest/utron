// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strconv

import (
	"fmt"
)

// ParseBool returns the boolean value represented by the string.
// It accepts built in standard value such as 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False,
// This method also accepts y, Y, yes, Yes, YES, e, E, Enabled, ENABLED for TRUE (1) and n, N, no, No, NO, d, D, diabled, Disabled, DISABLED for false (0)
// Any other value returns an error.
func ParseBool(str string) (bool, error) {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "y", "Y", "yes", "Yes", "YES", "e", "E", "enabled", "Enabled", "ENABLED":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "n", "N", "no", "No", "NO", "d", "D", "diabled", "Disabled", "DISABLED":
		return false, nil
	}
	return false, fmt.Errorf("ParseBool: cannot convert %s to boolean", str)
}

const (
	// FormatBoolTrueFalse used to define out to convert bool to string
	FormatBoolTrueFalse = 0
	// FormatBoolTrueFalse used to define out to convert bool to string
	FormatBoolYesNo = 1
	// FormatBoolTrueFalse used to define out to convert bool to string
	FormatBoolEnabledDisabled = 2
)

// FormatBool expands on the default version and allows returns of "True/False", "Yes/No" and "Enabled/Disabled" according to the value of b
func FormatBool(b bool, t int) string {
	result := ""

	switch t {
	case FormatBoolTrueFalse:
		if b {
			result = "True"
		} else {
			result = "False"
		}
	case FormatBoolYesNo:
		if b {
			result = "Yes"
		} else {
			result = "No"
		}
	case FormatBoolEnabledDisabled:
		if b {
			result = "Enabled"
		} else {
			result = "Disabled"
		}
	}

	return result
}

// AppendBool appends "true" or "false", according to the value of b,
// to dst and returns the extended buffer.
func AppendBool(dst []byte, b bool) []byte {
	if b {
		return append(dst, "true"...)
	}
	return append(dst, "false"...)
}
