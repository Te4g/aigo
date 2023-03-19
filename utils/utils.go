package utils

import "strings"

func RemoveFirstCarriageReturn(s string) string {
	return strings.Replace(s, "\n", "", 1)
}
