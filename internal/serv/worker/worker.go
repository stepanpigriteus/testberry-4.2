package worker

import (
	"fmt"

	"disgreps/domain"
	"disgreps/utils"
)

func Worker(cfg domain.Config, chunk []domain.Line) []domain.Line {
	fmt.Printf("Воркер обрабатывает %d строк\n", len(chunk))
	res := utils.Proccessor(cfg, chunk)
	fmt.Printf("Воркер нашел %d совпадений\n", len(res))

	for i, line := range res {
		fmt.Printf("  результат %d: LineNum=%d, Text(первые 50 символов)=%s\n",
			i, line.LineNum, line.Text[:min(50, len(line.Text))])
	}

	return res
}
