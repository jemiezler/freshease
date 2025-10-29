package helpers

import "time"

// helper: turn sql-nullable string into *string for DTO
func PtrIfNotNil(s *string) *string {
	if s == nil {
		return nil
	}
	return s
}

// helper: convert time to ISO string for DTO
func TimeToISOString(t *time.Time) *string {
	if t == nil {
		return nil
	}
	isoString := t.Format(time.RFC3339)
	return &isoString
}
