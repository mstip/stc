package internal

import (
	"bufio"
	"encoding/json"
	"log"
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	logFile, err := os.OpenFile("stc.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	// defer logFile.Close()
	logger = slog.New(slog.NewJSONHandler(logFile, nil))
}

func Logger() *slog.Logger {
	return logger
}

func FilterLogs(filterFunc func(map[string]interface{}) bool) ([]map[string]interface{}, error) {
	readFile, err := os.Open("stc.log")
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	buf := make([]byte, 0, 128*1024)
	fileScanner.Buffer(buf, 128*1024)
	fileScanner.Split(bufio.ScanLines)

	filteredLogs := []map[string]interface{}{}

	for fileScanner.Scan() {
		logEntry := map[string]interface{}{}
		if err = json.Unmarshal(fileScanner.Bytes(), &logEntry); err != nil {
			continue
		}

		if !filterFunc(logEntry) {
			continue
		}

		filteredLogs = append([]map[string]interface{}{logEntry}, filteredLogs...)
	}
	if fileScanner.Err() != nil {
		return nil, err
	}
	return filteredLogs, nil
}
