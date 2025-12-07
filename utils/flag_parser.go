package utils

import (
	"flag"
	"fmt"
	"log"
	"os"

	"disgreps/domain"
)

func FlagParser() domain.Config {
	var cfg domain.Config

	flag.IntVar(&cfg.After, "A", 0, "Вывести N строк после")
	flag.IntVar(&cfg.Before, "B", 0, "Вывести N строк перед")
	flag.IntVar(&cfg.Context, "C", 0, "Вывести N строк перед и после подстроки")
	flag.BoolVar(&cfg.CountOnly, "c", false, "Выводить только количество строк совпадающих с шаблоном ")
	flag.BoolVar(&cfg.IgnoreCase, "i", false, "Игнорировать регистр")
	flag.BoolVar(&cfg.Invert, "v", false, "Инвертировать фильтр: выводить строки, не содержащие шаблон")
	flag.BoolVar(&cfg.Fixed, "F", false, "Выполнять точное совпадение подстроки")
	flag.BoolVar(&cfg.LineNum, "n", false, "Выводить номер строки перед каждой найденной строкой")

	flag.BoolVar(&cfg.Mode, "work", false, "Запустить в режиме обработчика (worker). После флага указываются порты, например: ./grep --work 9002")

	flag.Parse()

	args := flag.Args()

	if cfg.Mode {
		if len(args) == 0 {
			log.Fatal("Укажите хотя бы один порт: ./grep --work 9002")
		}
		cfg.Ports = append(cfg.Ports, args...)
		return cfg
	}

	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: grep [flags] pattern [file]  (или ./grep --work <ports>)")
		os.Exit(1)
	}
	cfg.Pattern = args[0]
	if len(args) > 1 {
		cfg.Filename = args[1]
	}
	if cfg.Context > 0 {
		cfg.Before = cfg.Context
		cfg.After = cfg.Context
	}

	inWorkMode := false
	for _, arg := range args {
		if arg == "--work" {
			inWorkMode = true
			continue
		}
		if inWorkMode {
			cfg.Ports = append(cfg.Ports, arg)
		}
	}

	if len(cfg.Ports) == 0 {
		log.Fatal("Недостаточно портов!")
	}

	return cfg
}
