package tool

func ToString(val interface{}, fallback string) string {
	if val == nil {
		return fallback
	}
	s, ok := val.(string)
	if !ok {
		return fallback
	}
	return s
}
