package common

import "regexp"

type RootPatterns []string

func (s RootPatterns) IsRoot(str string) bool {
	for _, pattern := range s {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			// Fallback to exact string matching if pattern isn't valid regex
			if pattern == str {
				return true
			}
			continue
		}

		if regex.MatchString(str) {
			return true
		}
	}
	return false
}
