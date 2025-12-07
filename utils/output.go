package utils

import (
	"fmt"

	"disgreps/domain"
)

func OutRes(results []domain.Line, cfg domain.Config) {
	for _, result := range results {
		if cfg.CountOnly {
			fmt.Println(result.Text)
			continue
		}
		prefix := ""
		if cfg.LineNum {
			prefix += fmt.Sprintf("%d:", result.LineNum)
		}
		fmt.Printf("%s%s\n", prefix, result.Text)
	}
}
