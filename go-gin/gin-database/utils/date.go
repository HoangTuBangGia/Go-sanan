package utils

import "time"

func ParseBirthday(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}

	t, err := time.Parse("2006-01-02", *s)

	if err != nil {
		return nil
	}

	return &t
}
