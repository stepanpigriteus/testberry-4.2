package utils

import (
	"fmt"
	"os"
	"regexp"
)

//Proccessor обрабатывает массив строк
func Proccessor(cfg Сonfig, lines []Line) []Line {
	var result []Line
	matched := make(map[int]bool)
	count := 0

	pattern := cfg.pattern
	if cfg.fixed {
		pattern = regexp.QuoteMeta(cfg.pattern)
	}
	if cfg.ignoreCase {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Кривой паттерн: %v\n", err)
		os.Exit(1)
	}

	for i, line := range lines {
		isMatch := re.MatchString(line.text)
		if cfg.invert {
			isMatch = !isMatch
		}

		if isMatch {
			count++
			if !cfg.countOnly {
				for j := i - cfg.before; j < i; j++ {
					if j >= 0 && !matched[j] {
						result = append(result, lines[j])
						matched[j] = true
					}
				}
				if !matched[i] {
					result = append(result, line)
					matched[i] = true
				}
				for j := i + 1; j <= i+cfg.after && j < len(lines); j++ {
					if !matched[j] {
						result = append(result, lines[j])
						matched[j] = true
					}
				}
			}
		}
	}

	if cfg.countOnly {
		return []Line{{0, fmt.Sprintf("%d", count)}}
	}
	return result
}
