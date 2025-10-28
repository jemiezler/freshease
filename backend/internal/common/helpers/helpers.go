package helpers

// helper: turn sql-nullable string into *string for DTO
func PtrIfNotNil(s *string) *string {
	if s == nil {
		return nil
	}
	return s
}
