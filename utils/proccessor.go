package utils

import (
	"fmt"
	"os"
	"regexp"

	"disgreps/domain"
)

func Proccessor(cfg domain.Config, lines []domain.Line) []domain.Line {
	var result []domain.Line
	matched := make(map[int]bool)
	count := 0

	before := cfg.Before
	after := cfg.After
	if cfg.Mode {
		before = 0
		after = 0
	}

	pattern := cfg.Pattern
	if cfg.Fixed {
		pattern = regexp.QuoteMeta(pattern)
	}
	if cfg.IgnoreCase {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Кривой паттерн: %v\n", err)
		os.Exit(1)
	}

	for i, line := range lines {
		isMatch := re.MatchString(line.Text)
		if cfg.Invert {
			isMatch = !isMatch
		}

		if isMatch {
			count++

			if !cfg.CountOnly {

				for j := i - before; j < i; j++ {
					if j >= 0 && !matched[j] {
						result = append(result, lines[j])
						matched[j] = true
					}
				}
				if !matched[i] {
					result = append(result, line)
					matched[i] = true
				}
				for j := i + 1; j <= i+after && j < len(lines); j++ {
					if !matched[j] {
						result = append(result, lines[j])
						matched[j] = true
					}
				}
			}
		}
	}

	if cfg.CountOnly {
		return []domain.Line{
			{
				LineNum: 0,
				Text:    fmt.Sprintf("%d", count),
			},
		}
	}
	return result
}
