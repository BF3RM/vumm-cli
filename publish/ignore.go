package publish

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ignorePattern struct {
	regex  *regexp.Regexp
	negate bool
}

type FileIgnorer struct {
	patterns []*ignorePattern
}

func parseRegexFromLine(line string) (*regexp.Regexp, bool) {
	line = strings.TrimRight(line, "\r")

	// Strip comments [Rule 2]
	if strings.HasPrefix(line, "#") {
		return nil, false
	}

	// Trim string [Rule 3]
	line = strings.Trim(line, " ")

	if line == "" {
		return nil, false
	}

	// Handle [Rule 4] which negates the match for patterns leading with "!"
	negatePattern := false
	if line[0] == '!' {
		negatePattern = true
		line = line[1:]
	}

	// Handle [Rule 2, 4], when # or ! is escaped with a \
	// Handle [Rule 4] once we tag negatePattern, strip the leading ! char
	if regexp.MustCompile(`^(\#|\!)`).MatchString(line) {
		line = line[1:]
	}

	// Prepend with a /
	if regexp.MustCompile(`([^\/+]/.*\*\.)`).MatchString(line) && line[0] != '/' {
		line = "/" + line
	}

	// Escape the "." char
	line = regexp.MustCompile(`\.`).ReplaceAllString(line, `\.`)

	magicStar := "#$~"

	// Handle "/**/" usage
	if strings.HasPrefix(line, "/**/") {
		line = line[1:]
	}
	line = regexp.MustCompile(`/\*\*/`).ReplaceAllString(line, `(/|/.+/)`)
	line = regexp.MustCompile(`\*\*/`).ReplaceAllString(line, `(|.`+magicStar+`/)`)
	line = regexp.MustCompile(`/\*\*`).ReplaceAllString(line, `(|/.`+magicStar+`)`)

	// Handle escaping the "*" char
	line = regexp.MustCompile(`\\\*`).ReplaceAllString(line, `\`+magicStar)
	line = regexp.MustCompile(`\*`).ReplaceAllString(line, `([^/]*)`)

	// Handle escaping the "?" char
	line = strings.Replace(line, "?", `\?`, -1)

	line = strings.Replace(line, magicStar, "*", -1)

	// Temporary regex
	var expr = ""
	if strings.HasSuffix(line, "/") {
		expr = line + "(|.*)$"
	} else {
		expr = line + "(|/.*)$"
	}
	if strings.HasPrefix(expr, "/") {
		expr = "^(|/)" + expr[1:]
	} else {
		expr = "^(|.*/)" + expr
	}
	pattern, _ := regexp.Compile(expr)

	return pattern, negatePattern
}

func CompileFileIgnorer(filePath string) (*FileIgnorer, error) {
	bs, err := ioutil.ReadFile(filepath.Join(filePath, ".vummignore"))
	if err != nil {
		return nil, err
	}

	ignoreFile := &FileIgnorer{}
	lines := strings.Split(string(bs), "\n")

	for _, line := range lines {
		regex, negate := parseRegexFromLine(line)
		if regex != nil {
			pattern := &ignorePattern{regex, negate}
			ignoreFile.patterns = append(ignoreFile.patterns, pattern)
		}
	}

	return ignoreFile, nil
}

func (ignorer FileIgnorer) Matches(filePath string) bool {
	filePath = strings.ReplaceAll(filePath, string(os.PathSeparator), "/")

	matchesPath := false
	for _, pattern := range ignorer.patterns {
		if pattern.regex.MatchString(filePath) {
			if !pattern.negate {
				matchesPath = true
			} else if matchesPath {
				matchesPath = false
			}
		}
	}

	return matchesPath
}
