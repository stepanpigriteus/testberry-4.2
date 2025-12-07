package master

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"disgreps/domain"
	"disgreps/utils"
)

func Master(cfg domain.Config, host string) []domain.Line {
	onlineWorkers := WorkerChecker(cfg, host)
	if len(onlineWorkers) == 0 {
		log.Fatal("нет живых воркеров")
	}
	fmt.Println("онлайн воркеры:", onlineWorkers)

	text := utils.ReadInput(cfg)
	fmt.Printf("ВCего строк для обработки: %d\n", len(text))
	fmt.Printf("Первые 5 строк:\n")
	for i := 0; i < 5 && i < len(text); i++ {
		fmt.Printf("%d: %s\n", text[i].LineNum, text[i].Text[:min(50, len(text[i].Text))])
	}

	workersNum := len(onlineWorkers)
	length := len(text)
	partSize := (length + workersNum - 1) / workersNum

	var result []domain.Line
	var wg sync.WaitGroup

	workerIndex := 0

	quorum := workersNum/2 + 1
	ready := make(chan []domain.Line, workersNum)

	for i := 0; i < length; i += partSize {
		end := i + partSize
		if end > length {
			end = length
		}

		chunk := text[i:end]
		port := onlineWorkers[workerIndex]
		workerIndex++

		wg.Add(1)
		go func(chunk []domain.Line, port string) {
			defer wg.Done()

			taskCfg := cfg
			taskCfg.Mode = true

			res := WorkDespenser(port, host, chunk, taskCfg)
			ready <- res
		}(chunk, port)
	}

	count := 0
	for res := range ready {
		result = append(result, res...)
		count++

		if count >= quorum {
			fmt.Println("⚡ Кворум достигнут, прекращаем ожидание остальных")
			break
		}
	}

	wg.Wait()
	close(ready)

	fmt.Printf("Итого найдено Совпадений: %d\n", len(result))
	return result
}

func CheckWorkerStatus(host string, port string) (string, error) {
	url := fmt.Sprintf("http://%s:%s/on", host, port)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to connect to port %s: %w", port, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("port %s returned status %d", port, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response from port %s: %w", port, err)
	}

	return string(body), nil
}

func WorkDespenser(port, host string, chunk []domain.Line, cfg domain.Config) []domain.Line {
	if len(chunk) == 0 {
		return []domain.Line{}
	}

	url := "http://" + host + ":" + port + "/load"
	body := struct {
		Chunk []domain.Line `json:"chunk"`
		Cfg   domain.Config `json:"cfg"`
	}{
		Chunk: chunk,
		Cfg:   cfg,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		log.Printf("воркер %s: ошибка json.Marshal: %v", port, err)
		return nil
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		log.Printf("воркер %s: ошибка Cоздания запроCа: %v", port, err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("воркер %s: ошибка отправки запроCа: %v", port, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("воркер %s вернул %d: %s", port, resp.StatusCode, string(bodyBytes))
		return nil
	}

	var result []domain.Line
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("воркер %s: ошибка декодирования JSON: %v", port, err)
		return nil
	}

	return result
}

func WorkerChecker(cfg domain.Config, host string) []string {
	var result []string
	workers := len(cfg.Ports)
	resultChan := make(chan string, workers)
	var wg sync.WaitGroup
	for _, r := range cfg.Ports {
		wg.Add(1)
		go func(port string) {
			defer wg.Done()
			res, err := CheckWorkerStatus(host, port)
			if err != nil {
				return
			}
			if res == "on-ok" {
				resultChan <- port
			}
		}(r)
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for res := range resultChan {
		result = append(result, res)
	}
	return result
}
