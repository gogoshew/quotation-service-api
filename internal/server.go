package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"qsapi/pkg/pg_db"
	"qsapi/pkg/repo_cron"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Server struct {
	ctx    context.Context
	router *mux.Router
	db     *pg_db.DatabasePg
	ts     repo_cron.ITaskScheduler
	http.Server
}

func NewServer(ctx context.Context, router *mux.Router, db *pg_db.DatabasePg, ts repo_cron.ITaskScheduler) (s *Server) {
	s = &Server{
		ctx:    ctx,
		router: router,
		db:     db,
		ts:     ts,
		Server: http.Server{
			Addr:         ":8080",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			Handler:      router,
		},
	}

	router.HandleFunc("/quotation/update/{currency}", s.updateQuotation).Methods("POST")
	router.HandleFunc("/quotation/{id}", s.getQuotationByID).Methods("GET")
	router.HandleFunc("/quotation/{currency}", s.getQuotationValue).Methods("GET")

	return
}

type getRatesResp struct {
	Date  time.Time          `json:"date"`
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

type updateReq struct {
	CurrencyCode string `json:"currencyCode"`
}

type updateResp struct {
	UpdateID uuid.UUID `json:"updateID"`
}

func (s *Server) updateQuotation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	vars := mux.Vars(r)
	apiURI := "https://api.currencybeacon.com/v1"

	qp := url.Values{}
	qp.Set("api_key", "y3wA4rW34r5oXGaX592nns8JgouvA6Wm")
	qp.Set("symbols", vars["currency"])

	reqURL := fmt.Sprintf("%s?%s", apiURI, qp.Encode())

	resp, err := http.Get(reqURL)
	if err != nil {
		log.Printf("Error requesting currency API: %v\n", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v\n", err)
	}

	var curBeaconResp getRatesResp
	if err := json.Unmarshal(body, &curBeaconResp); err != nil {
		log.Fatalf("Error unmarshaling response: %v\n", err)
	}

	var reqBody updateReq
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Fatalf("Error decoding request body: %v\n", err)
	}

	// Забираю время в которое будет произведено обновление (НУЖЕН ID Таски cron)
	s.ts.GetResolveTime(s.ts.GetMainTaskID())

	// Пишу в таблицу RateUpdates

	// Возвращаю ID обновления в респонсе

}

func (s *Server) getQuotationByID(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) getQuotationValue(w http.ResponseWriter, r *http.Request) {

}

//func ValidateCurrencies() {
//	rates := map[string]bool{
//		"USD": true,
//		"EUR": true,
//		"MXN": true,
//		"GEL": true,
//		"RUB": true,
//	}
//}
