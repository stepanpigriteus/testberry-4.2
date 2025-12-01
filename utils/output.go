package utils

import "fmt"

func OutRes(results []Line, cfg Сonfig) {
	for _, result := range results {
		if cfg.countOnly {
			fmt.Println(result.text)
			continue
		}
		prefix := ""
		if cfg.lineNum {
			prefix += fmt.Sprintf("%d:", result.lineNum)
		}
		fmt.Printf("%s%s\n", prefix, result.text)
	}
}
