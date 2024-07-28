package regex

import (
	"regexp"
)

var (
	// OnlyDigits is a regex that matches any non-digit character.
	OnlyDigits = regexp.MustCompile(`\D`)

	// SpaceSequence is a regex that matches any sequence of spaces.
	SpaceSequence = regexp.MustCompile(`\s+`)
)
