package experience

import "math/rand"

func RandInt(min int, max int) int {
	return rand.Intn(max + 1 - min) + min
}
