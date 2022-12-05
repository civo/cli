package utility

// TruncateID is a function to truncate the ID
func TruncateID(id string) string {
	if len(id) > 6 {
		return id[:6]
	}
	return id
}
