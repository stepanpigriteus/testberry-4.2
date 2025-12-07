package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"disgreps/domain"
)

func ReadInput(cfg domain.Config) []domain.Line {
	var scanner *bufio.Scanner
	if cfg.Filename == "" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(cfg.Filename)
		if err != nil {
			if errors.Is(err, os.ErrPermission) {
				fmt.Fprintln(os.Stderr, "Недостаточно прав для открытия файла")
				os.Exit(1)
			}
			if errors.Is(err, os.ErrNotExist) {
				fmt.Fprintln(os.Stderr, "Файл не сущуествует")
				os.Exit(1)
			}
			return nil
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	var lines []domain.Line
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		lines = append(lines, domain.Line{LineNum: lineNum, Text: scanner.Text()})
	}
	return lines
}
