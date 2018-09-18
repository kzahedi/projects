package util

// ListContains checks if entry is in list
func ListContains(list *[]string, entry string) bool {
	for _, v := range *list {
		if v == entry {
			return true
		}
	}
	return false
}
