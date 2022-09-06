package strings

func StringSliceToInterfaces(strings ...string) []interface{} {
	result := make([]interface{}, len(strings))
	for i := range strings {
		result[i] = strings[i]
	}

	return result
}
