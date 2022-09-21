package version

import (
	"regexp"
)

var versionRegexp = regexp.MustCompile(`(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)`)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p Parser) Parse(version string) string {
	result := versionRegexp.FindString(version)
	if result != "" {
		return result
	}

	return version
}
