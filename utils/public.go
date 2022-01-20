package utils

import "regexp"

var SemverRegex = regexp.MustCompile("[0-9]+(\\.[0-9]+)+")
