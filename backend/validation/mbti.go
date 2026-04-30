package validation

import (
	"fmt"
	"strings"
)

// Valid MBTI types (16 types)
var validMBTITypes = map[string]bool{
	// Analysts (NT)
	"INTJ": true,
	"INTP": true,
	"ENTJ": true,
	"ENTP": true,
	// Diplomats (NF)
	"INFJ": true,
	"INFP": true,
	"ENFJ": true,
	"ENFP": true,
	// Sentinels (SJ)
	"ISTJ": true,
	"ISFJ": true,
	"ESTJ": true,
	"ESFJ": true,
	// Explorers (SP)
	"ISTP": true,
	"ISFP": true,
	"ESTP": true,
	"ESFP": true,
}

// Valid MBTI extensions
var validMBTIExtensions = map[string]bool{
	"A": true, // Assertive
	"T": true, // Turbulent
}

// ValidateMBTI validates an MBTI string in the format "TYPE-EXT" (e.g., "INTP-A", "ENFJ-T")
// Returns nil if valid, error otherwise
func ValidateMBTI(mbti string) error {
	if mbti == "" {
		return nil // MBTI is optional
	}

	// Split by dash
	parts := strings.Split(mbti, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid MBTI format: must be TYPE-EXTENSION (e.g., 'INTP-A', 'ENFJ-T'), got '%s'", mbti)
	}

	mbtiType := strings.ToUpper(parts[0])
	mbtiExt := strings.ToUpper(parts[1])

	// Validate type
	if !validMBTITypes[mbtiType] {
		return fmt.Errorf("invalid MBTI type: '%s'. Must be one of the 16 MBTI types (INTJ, INTP, ENTJ, ENTP, INFJ, INFP, ENFJ, ENFP, ISTJ, ISFJ, ESTJ, ESFJ, ISTP, ISFP, ESTP, ESFP)", mbtiType)
	}

	// Validate extension
	if !validMBTIExtensions[mbtiExt] {
		return fmt.Errorf("invalid MBTI extension: '%s'. Must be 'A' (Assertive) or 'T' (Turbulent)", mbtiExt)
	}

	return nil
}
