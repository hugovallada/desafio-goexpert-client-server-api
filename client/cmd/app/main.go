package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hugovallada/client-server-api/client/conversion"
)

var (
	URL       = "http://localhost:8080/cotacao"
	FILE_PATH = "./files/conversao.txt"
)

func main() {
	value, err := getDolarValueInBRL()
	if err != nil {
		log.Fatal("Não foi possível completar o processamento!", err)
	}
	err = saveToFile(value)
	if err != nil {
		log.Fatal("Não foi possível completar o processamento!", err)
	}
}

func saveToFile(value float64) error {
	file, err := os.OpenFile(FILE_PATH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.WriteString(fmt.Sprintf("Dolar:%v\n", value)); err != nil {
		return err
	}
	return nil
}

func getDolarValueInBRL() (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var dolarValue conversion.Value
	err = json.NewDecoder(resp.Body).Decode(&dolarValue)
	if err != nil {
		return 0, err
	}
	return dolarValue.DolarValue, nil
}
