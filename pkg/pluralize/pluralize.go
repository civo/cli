package pluralize

import (
	"fmt"
	"strings"
)

// Pluralize returns a plural suffix if the provided size is not 1.
// By default, use 's' as the suffix:
// - If value is 0, vote{{ value|pluralize }} display "votes".
// - If value is 1, vote{{ value|pluralize }} display "vote".
// - If value is 2, vote{{ value|pluralize }} display "votes".

// If an argument is provided, use that string instead:
// - If value is 0, class{{ value|pluralize:"es" }} display "classes".
// - If value is 1, class{{ value|pluralize:"es" }} display "class".
// - If value is 2, class{{ value|pluralize:"es" }} display "classes".

// - If the provided argument contains a comma, use the text before the comma for the singular case and the text after the comma for the plural case:
// - If value is 0, cand{{ value|pluralize:"y,ies" }} display "candies".
// - If value is 1, cand{{ value|pluralize:"y,ies" }} display "candy".
// - If value is 2, cand{{ value|pluralize:"y,ies" }} display "candies".
func Pluralize(size int, word string, arg ...string) string {
	switch size {
	case 1:
		if len(arg) == 0 {
			return word
		} else {
			bits := splitArg(arg[0])
			if len(bits) > 1 {
				return fmt.Sprintf("%s%s", word, bits[0])
			} else {
				return word
			}
		}
	default:
		if len(arg) == 0 {
			return fmt.Sprintf("%ss", word)
		} else {
			bits := splitArg(arg[0])
			if len(bits) > 1 {
				return fmt.Sprintf("%s%s", word, bits[1])
			} else {
				return fmt.Sprintf("%s%s", word, arg[0])
			}
		}
	}
}

func splitArg(arg string) (result []string) {
	value := strings.Split(arg, ",")
	if len(value) > 1 {
		return value
	}
	return []string{}
}

// Has returns the either the singular or plural version of the verb "to have" based on the `size` input.
func Has(size int) string {
	return Pluralize(size, "ha", "s,ve")
}
