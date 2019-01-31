package utility

import (
	"strconv"
	"strings"
)

func ToInt(s string) (int, error) {
	return strconv.Atoi(strings.Replace(s, ",", "", -1))
}
