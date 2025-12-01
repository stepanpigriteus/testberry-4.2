package utils

import (
	"flag"
	"fmt"
	"os"
)

// FlagParser парсим флаги
func FlagParser() Сonfig {
	var cfg Сonfig
	flag.IntVar(&cfg.after, "A", 0, "Вывести N строк после")
	flag.IntVar(&cfg.before, "B", 0, "Вывести N строк перед")
	flag.IntVar(&cfg.context, "C", 0, "Вывести N строк перед и после подстроки")
	flag.BoolVar(&cfg.countOnly, "c", false, "Выводить только количество строк совпадающих с шаблоном ")
	flag.BoolVar(&cfg.ignoreCase, "i", false, "Игнорировать регистр")
	flag.BoolVar(&cfg.invert, "v", false, "Инвертировать фильтр: выводить строки, не содержащие шаблон")
	flag.BoolVar(&cfg.fixed, "F", false, "Выполнять точное совпадение подстроки")
	flag.BoolVar(&cfg.lineNum, "n", false, "Выводить номер строки перед каждой найденной строкой")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: grep [flags] pattern [file]")
		os.Exit(1)
	}
	cfg.pattern = args[0]
	if len(args) > 1 {
		cfg.filename = args[1]
	}
	if cfg.context > 0 {
		cfg.before = cfg.context
		cfg.after = cfg.context
	}

	return cfg
}
