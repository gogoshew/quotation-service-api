package internal

import (
	"context"
	"encoding/json"
	"errors"
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

	router.HandleFunc("/quotation/update", s.updateQuotation).Methods("POST")
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
	Currency string `json:"currency"`
}

type updateResp struct {
	UpdateID uuid.UUID `json:"updateID"`
}

func (s *Server) updateQuotation(w http.ResponseWriter, r *http.Request) {
	var updateID string

	var reqBody updateReq
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Fatalf("Error decoding request body: %v\n", err)
	}

	rateValue, err := s.getRate(reqBody.Currency)
	if err != nil {
		log.Fatalf("Error getting rate from currency API: %v\n", err)
	}

	row, err := s.db.GetRowByCurrency(reqBody.Currency)
	if err != nil {
		log.Fatalf("Error getting row from rate buffer table: %v\n", err)
	}
	updateID = row.UpdateID

	if row.UpdateFlag == false {
		updateID = uuid.New().String()
		rateInfo := &pg_db.RateBuffer{
			UpdateID:   updateID,
			Currency:   reqBody.Currency,
			Value:      rateValue,
			Base:       "USD",
			UpdateFlag: true,
		}

		if err := s.db.CreateRowBuffer(rateInfo); err != nil {
			log.Fatalf("Error inserting row rate buffer table: %v\n", err)
		}
	}

	// Возвращаю ID обновления в респонсе

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(updateID))
}

type QuotationByIDReq struct {
	UpdateID string
}

type QuotationByIDResp struct {
	Value    float64
	UpdateAt time.Time
}

func (s *Server) getQuotationByID(w http.ResponseWriter, r *http.Request) {
	var req QuotationByIDReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Fatalf("Error decoding request body: %v\n", err)
	}

	//var row pg_db.RateBuffer
	row, err := s.db.GetRowByUpdateID(req.UpdateID)
	if err != nil {
		log.Fatalf("Error getting quotation by ID: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var resp QuotationByIDResp

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(updateID))
}

func (s *Server) getQuotationValue(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) getRate(cur string) (float64, error) {
	apiURI := "https://api.currencybeacon.com/v1"

	qp := url.Values{}
	qp.Set("api_key", "y3wA4rW34r5oXGaX592nns8JgouvA6Wm")
	qp.Set("symbols", cur)

	reqURL := fmt.Sprintf("%s?%s", apiURI, qp.Encode())

	resp, err := http.Get(reqURL)
	if err != nil {
		log.Fatalf("Error requesting currency API: %v\n", err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v\n", err)
		return 0, err
	}

	var curBeaconResp getRatesResp
	if err := json.Unmarshal(body, &curBeaconResp); err != nil {
		log.Fatalf("Error unmarshaling response: %v\n", err)
		return 0, err
	}
	rate := curBeaconResp.Rates[cur]
	if rate == 0 {
		err = errors.New("quotation rate mustn't be zero")
		log.Fatalf("%v\n", err)
		return 0, err
	}
	return rate, nil
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
