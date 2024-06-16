package validation

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// HttpError wraps an error message and an HTTP status code
type HttpError struct {
	Message string
	Code    int
}

func (e HttpError) Error() string {
	return e.Message
}

// ValidateID checks if the given ID is a valid positive integer
func ValidateID(idStr string) (int, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		return 0, errors.New("invalid post ID")
	}
	return id, nil
}

// ValidateQuery checks the validity of the query parameter
func ValidateQuery(query string) error {
	query = strings.TrimSpace(query)

	if query == "" {
		return HttpError{Message: "query parameter is required", Code: http.StatusBadRequest}
	}
	if len(query) > 100 {
		return HttpError{Message: "query parameter is too long", Code: http.StatusBadRequest}
	}
	if !utf8.ValidString(query) {
		return HttpError{Message: "query parameter is not valid UTF-8", Code: http.StatusBadRequest}
	}

	isValid := regexp.MustCompile(`^[a-zA-Z0-9\s]+$`).MatchString
	if !isValid(query) {
		return HttpError{Message: "query parameter contains invalid characters", Code: http.StatusBadRequest}
	}

	return nil
}
