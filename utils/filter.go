package utils

func Filter(params map[string]interface{}, filters []string) {
	for _, f := range filters {
		delete(params, f)
	}
}
