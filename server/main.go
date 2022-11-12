package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Moeda struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type MoedaDb struct {
	ID      int `gorm:"primaryKey"`
	Cotacao string
}

const URL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
const REQUEST_MAX_DURATION = 200 * time.Millisecond

func main() {
	http.HandleFunc("/cotacao", fetchCurrency)

	http.ListenAndServe(":8080", nil)
}

func fetchCurrency(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&MoedaDb{})

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, REQUEST_MAX_DURATION)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var moeda Moeda

	json.Unmarshal(body, &moeda)

	db.Create(&MoedaDb{
		Cotacao: moeda.Usdbrl.Bid,
	})

	w.Write(body)
}
