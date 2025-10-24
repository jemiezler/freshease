package helpers

// helper: turn sql-nullable string into *string for DTO
func PtrIfNotNil(s *string) *string { return s }
