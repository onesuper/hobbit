package hobbit

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Prettify a couple of symbols.
// If `num` <= 0, empty string will be return.
func SeveralSymbols(symbol byte, num int) string {
	if num > 5 {
		return fmt.Sprintf("%c x %d", symbol, num)
	}

	s := ""
	for i := 0; i < num; i++ {
		s += string(symbol)
	}
	return s
}

// Used when casting the fixed-length string
func FixToLength(str string, length int, filler byte) string {
	if len(str) > length {
		return str[0:length]
	} else if len(str) < length {
		return padToLength(str, length, filler)
	} else {
		return str
	}
}

// Padding " " to the left and the right of a string until its length
// meets requirement
func padToLength(str string, length int, filler byte) string {
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

// Util for passing a coordinate from a command
func ParseCoord(command string) (int, int, error) {
	tokens := strings.Split(command, "-")
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
