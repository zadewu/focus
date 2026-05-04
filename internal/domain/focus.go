package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const nameSeparator = "__"

type Focus struct {
	Name      string
	CreatedAt time.Time
	Archived  bool
}

type Note struct {
	Timestamp time.Time
	Message   string
}

type File struct {
	Name    string
	Content []byte
}

func ValidateName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if strings.ContainsAny(name, " \t\n/\\") {
		return fmt.Errorf("invalid focus name %q: must not contain spaces, slashes, or newlines", name)
	}
	if strings.Contains(name, nameSeparator) {
		return fmt.Errorf("invalid focus name %q: %q is reserved as a separator", name, nameSeparator)
	}
	if strings.HasPrefix(name, "archive") {
		return fmt.Errorf("invalid focus name %q: cannot start with 'archive'", name)
	}
	if len(name) > 100 {
		return fmt.Errorf("invalid focus name %q: too long (max 100 chars)", name)
	}
	return nil
}

// GenerateFullName prepends a creation timestamp to the short name.
// Format: 2006-01-02-1504__shortName (e.g. 2026-05-03-2125__my-task)
func GenerateFullName(shortName string, t time.Time) string {
	return t.Format("2006-01-02-1504") + nameSeparator + shortName
}

// ExtractShortName strips any timestamp prefix from a branch name, returning
// just the user-supplied short name. Handles three formats:
//   - YYYY-MM-DD-HHmm__name  (current)
//   - YYYY-MM-DD--name       (legacy)
//   - name                   (plain, no prefix)
func ExtractShortName(fullName string) string {
	if _, after, found := strings.Cut(fullName, nameSeparator); found {
		return after
	}
	if isLegacyPrefixed(fullName) {
		return fullName[12:]
	}
	return fullName
}

// isLegacyPrefixed reports whether s begins with a YYYY-MM-DD-- date prefix.
// Validates digit positions to avoid false positives on branches like feat-my-do--work.
func isLegacyPrefixed(s string) bool {
	if len(s) < 12 {
		return false
	}
	return isASCIIDigit(s[0]) && isASCIIDigit(s[1]) && isASCIIDigit(s[2]) && isASCIIDigit(s[3]) &&
		s[4] == '-' &&
		isASCIIDigit(s[5]) && isASCIIDigit(s[6]) &&
		s[7] == '-' &&
		isASCIIDigit(s[8]) && isASCIIDigit(s[9]) &&
		s[10] == '-' && s[11] == '-'
}

func isASCIIDigit(b byte) bool { return b >= '0' && b <= '9' }
