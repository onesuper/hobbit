package hobbit

import (
	"errors"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
)

// Hash a string to a 32-bit int.
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// Prettify a couple of symbols. If `num` <= 0, empty string will be return.
func PrettySymbols(symbol byte, num int) string {
	if num > 5 {
		return fmt.Sprintf("%c x %d", symbol, num)
	}

	s := ""
	for i := 0; i < num; i++ {
		s += string(symbol)
	}
	return s
}

// Casting a fixed-length string, padded with `filler` as character.
func FixToLength(str string, length int, filler byte) string {
	if len(str) > length {
		return str[0:length]
	} else if len(str) < length {
		return PadToLength(str, length, filler)
	} else {
		return str
	}
}

// Padding " " to the left and the right of a string until its length
// meets requirement.
func PadToLength(str string, length int, filler byte) string {
	n := length - len(str)
	for n > 0 {
		if n%2 == 0 {
			str = string(filler) + str
		} else {
			str = str + string(filler)
		}
		n--
	}
	return str
}

// Passing a coordinate like 'x-y' from a string.
func ParseCoord(str string) (int, int, error) {
	tokens := strings.Split(str, "-")
	if len(tokens) != 2 {
		return -1, -1, errors.New("wrong axis count")
	}
	x, err1 := strconv.Atoi(tokens[0])
	if err1 != nil {
		return -1, -1, err1
	}
	y, err2 := strconv.Atoi(tokens[1])
	if err2 != nil {
		return -1, -1, err2
	}
	return x, y, nil
}
