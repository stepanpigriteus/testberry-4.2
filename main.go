package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	httpsh "disgreps/internal/http"
	"disgreps/internal/serv/master"
	"disgreps/utils"
)

func main() {
	ctx := context.Background()
	cfg := utils.FlagParser()
	host := "localhost"
	var wg sync.WaitGroup

	switch len(cfg.Ports) {
	case 0:
		log.Fatal("zero ports")
	case 1:
		fmt.Println("режим обработчика")
		workerCfg := cfg
		workerCfg.Mode = true
		server := httpsh.NewServer(cfg.Ports[0], host, false, workerCfg)
		if err := server.RunServer(ctx); err != nil {
			log.Fatal("ошибка запуска worker сервера")
		}
		lines := utils.ReadInput(cfg)
		resultLines := utils.Proccessor(cfg, lines)
		utils.OutRes(resultLines, cfg)

	default:
		fmt.Println("режим мастера")
		server := httpsh.NewServer("8081", host, true, cfg)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := server.RunServer(ctx); err != nil {
				log.Printf("ошибка запуска worker сервера: %v", err)
			}
		}()
		results := master.Master(cfg, host)
		fmt.Println(results)
	}
	wg.Wait()
}
