// Copyright © 2024 Yangxinxin. All right reserved.
// Confidential and Proprietary
package strkit

import "regexp"

var (
	CHART = regexp.MustCompile(`^[a-z]([a-z]|[0-9]|\-|\.){0,60}([a-z]|[0-9])$`)
	XSS   = regexp.MustCompile(`[~{}"'<>?]`)
)

func IsValidChart(name string) bool {
	return CHART.MatchString(name)
}

func ValidateXss(name string) bool {
	return !XSS.MatchString(name)
}
