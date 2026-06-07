package idutil

import (
	"unicode"

	"lineoa-miniapp/pkg/mentalhealthcache"
)

// ClassifyID determines whether a raw ID string is a Thai national ID
// or a passport ID, then returns it in the correct field.
//
// Rules:
//   - Normalize first (strips spaces, hyphens, converts scientific notation)
//   - Exactly 13 digits and all numeric → national_id
//   - Anything else (contains letters, wrong length) → passport_id
func ClassifyID(raw string) (nationalID, passportID string) {
	norm := mentalhealthcache.NormalizeID(raw)
	if norm == "" {
		return "", ""
	}
	if len(norm) == 13 && allDigits(norm) {
		return norm, ""
	}
	return "", norm
}

func allDigits(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
