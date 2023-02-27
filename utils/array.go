package utils

func RemoveEmptyElementsFromStringArray(arr []string) []string {
	var r []string
	for _, str := range arr {
		if str != "" {
			r = append(r, str)
		}
	}
	arr = r
	return r
}
