package hexboard

// Basic type for a hex
type Hex struct {
	width  int
	height int
	side   int
}

// Functions for processing the hex template.
// Used when casting the appearance of the hex
func fixToLength(str string, length int) string {
	if len(str) > length {
		return str[0:length]
	} else if len(str) < length {
		return padToLength(str, length)
	} else {
		return str
	}
}

// Padding " " to the left and the right of a string until its length
// meets requirement
func padToLength(str string, length int) string {
	n := length - len(str)
	for n > 0 {
		if n%2 == 0 {
			str = " " + str
		} else {
			str = str + " "
		}
		n--
	}
	return str
}
