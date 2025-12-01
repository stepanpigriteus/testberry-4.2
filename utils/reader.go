package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

// ReadInput читает из файла или стдин
func ReadInput(cfg Сonfig) []Line {
	var scanner *bufio.Scanner
	if cfg.filename == "" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(cfg.filename)
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

	var lines []Line
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		lines = append(lines, Line{lineNum, scanner.Text()})
	}
	return lines
}
