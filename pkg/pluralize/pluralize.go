package pluralize

import (
	"fmt"
	"strings"
)

func Pluralize(size int, word string, arg ...string) string {

	/*
	   Return a plural suffix if the value is not 1, '1', or an object of
	   length 1. By default, use 's' as the suffix:
	   * If value is 0, vote{{ value|pluralize }} display "votes".
	   * If value is 1, vote{{ value|pluralize }} display "vote".
	   * If value is 2, vote{{ value|pluralize }} display "votes".

	   If an argument is provided, use that string instead:
	   * If value is 0, class{{ value|pluralize:"es" }} display "classes".
	   * If value is 1, class{{ value|pluralize:"es" }} display "class".
	   * If value is 2, class{{ value|pluralize:"es" }} display "classes".

	   If the provided argument contains a comma, use the text before the comma
	   for the singular case and the text after the comma for the plural case:
	   * If value is 0, cand{{ value|pluralize:"y,ies" }} display "candies".
	   * If value is 1, cand{{ value|pluralize:"y,ies" }} display "candy".
	   * If value is 2, cand{{ value|pluralize:"y,ies" }} display "candies".
	*/

	switch {
	case size == 0:
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
	case size == 1:
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
	case size > 1:
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

	return ""
}

func splitArg(arg string) (result []string) {
	value := strings.Split(arg, ",")
	if len(value) > 1 {
		return value
	}
	return []string{}
}
