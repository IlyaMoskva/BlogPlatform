package validation

import (
	"errors"
	"strconv"
)

// ValidateID checks if the given ID is a valid positive integer
func ValidateID(idStr string) (int, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		return 0, errors.New("invalid post ID")
	}
	return id, nil
}
