package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/hugovallada/client-server-api/server/conversion"
	_ "github.com/mattn/go-sqlite3"
)

var (
	URL                    = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	INTERNAL_ERROR_MESSAGE = "Erro interno do servidor"
)

func init() {
	db, err := sql.Open("sqlite3", "./data/price.db")
	if err != nil {
		panic(err)
	}
	_, tableCheck := db.Query("select * from prices")
	if tableCheck != nil {
		log.Println("Criando a tabela 'prices'")
		db.Exec(`create table prices(
			id integer primary key autoincrement,
			  code text,
			  codein text,
			  name text,
			  high decimal(10,5),
			  low decimal(10,5),
			  varBid decimal(10,5),
			  pctChange decimal(10,5),
			  bid decimal(10,5),
			  ask decimal(10,5),
			  timestamp integer,
			  create_date text
		);`)
	}

}

func main() {
	http.HandleFunc("/cotacao", GetPrice)
	http.ListenAndServe(":8080", nil)
}

func GetPrice(w http.ResponseWriter, r *http.Request) {
	price, err := findUSDBRLPrice()
	if err != nil {
		log.Println(err)
		http.Error(w, INTERNAL_ERROR_MESSAGE, http.StatusInternalServerError)
		return
	}
	err = saveToDB(price)
	if err != nil {
		log.Println(err)
		http.Error(w, INTERNAL_ERROR_MESSAGE, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	dolarValue := conversion.ConversionResponse{DolarValue: price.USDBRL.Bid}
	json.NewEncoder(w).Encode(dolarValue)
}

func saveToDB(price conversion.USDBRL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	db, err := sql.Open("sqlite3", "./data/price.db")
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(
		"insert into prices (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date)" +
			"values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, price.USDBRL.Code, price.USDBRL.CodeIn,
		price.USDBRL.Name, price.USDBRL.High, price.USDBRL.Low,
		price.USDBRL.VarBid, price.USDBRL.PctChange, price.USDBRL.Bid,
		price.USDBRL.Ask, price.USDBRL.Timestamp, price.USDBRL.CreateDate)
	return err
}

func findUSDBRLPrice() (conversion.USDBRL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		return conversion.USDBRL{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return conversion.USDBRL{}, err
	}
	defer resp.Body.Close()
	var price conversion.USDBRL
	err = json.NewDecoder(resp.Body).Decode(&price)
	if err != nil {
		return price, err
	}
	return price, nil
}
