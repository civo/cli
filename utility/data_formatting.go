package utility

import "strings"

// BoolToYesNo returns Yes or No depending on the value of a boolean
// This wouldn't be as necessary if Go had a ternary operator
func BoolToYesNo(d bool) string {
	ret := "No"
	if d {
		ret = "Yes"
	}
	return ret
}

// GetStringMap getStringMap convert string in the format a:1,b:2,c:3
// in map[string]string, this is util for the flags StringArrayVarP
func GetStringMap(s string) map[string]string {
	entries := strings.Split(s, ",")

	m := make(map[string]string)
	for _, e := range entries {
		tokens := strings.Split(e, ":")
		k := strings.TrimSpace(tokens[0])
		v := strings.TrimSpace(tokens[1])
		m[k] = v
	}

	return m
}
