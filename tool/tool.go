package tool

import "time"

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

func ConvertUTCToISO(ts time.Time) string {
	if ts.IsZero() {
		return ""
	}
	return ts.UTC().Format(time.RFC3339)
}
