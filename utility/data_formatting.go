package utility

// BoolToYesNo returns Yes or No depending on the value of a boolean
// This wouldn't be as necessary if Go had a ternary operator
func BoolToYesNo(d bool) string {
	ret := "No"
	if d {
		ret = "Yes"
	}
	return ret
}
