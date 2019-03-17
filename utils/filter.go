package utils

// Filter filter out keys in params
func Filter(params map[string]interface{}, filters []string) {
	for _, f := range filters {
		delete(params, f)
	}
}
