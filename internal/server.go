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
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	eur = "EUR"
	mxn = "MXN"
	gel = "GEL"
	//base   = "USD"
	apiURI = "https://api.currencybeacon.com/v1/latest"
	apiKey = "y3wA4rW34r5oXGaX592nns8JgouvA6Wm"
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
	router.HandleFunc("/quotation/id", s.getQuotationByID).Methods("GET")
	router.HandleFunc("/quotation/latest", s.getLatest).Methods("GET")

	return
}

type getRatesResp struct {
	Meta     Meta               `json:"meta"`
	Response Response           `json:"response"`
	Date     time.Time          `json:"date"`
	Base     string             `json:"base"`
	Rates    map[string]float64 `json:"rates"`
}

type Meta struct {
	Code       int    `json:"code"`
	Disclaimer string `json:"disclaimer"`
}

type Response struct {
	Date  time.Time          `json:"date"`
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

type updateResp struct {
	UpdateID string `json:"updateID"`
}

func (s *Server) updateQuotation(w http.ResponseWriter, r *http.Request) {
	cur := strings.ToTitle(r.URL.Query().Get("currency"))

	if !ValidateCurrency(cur) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No such currency"))
		return
	}

	rateValue, err := s.getRate(cur)
	if err != nil {
		log.Printf("Error getting rate from currency API: %v\n", err)
		return
	}

	row, err := s.db.GetRowByCurBuffer(cur)
	if err != nil {
		log.Printf("Error getting row from rate buffer table: %v\n", err)
		return
	}
	updateID := row.UpdateID

	if row.UpdateFlag == false {
		updateID = uuid.New().String()
		rateInfo := pg_db.BufferRate{
			UpdateID:   updateID,
			Value:      rateValue,
			UpdateFlag: true,
		}

		if err := s.db.UpdateBuffer(cur, rateInfo); err != nil {
			log.Printf("Error inserting row rate buffer table: %v\n", err)
			return
		}
	}

	resp := updateResp{
		UpdateID: updateID,
	}
	res, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling response: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

type QuotationByIDResp struct {
	Value    float64   `json:"value"`
	UpdateAt time.Time `json:"updateAt"`
}

func (s *Server) getQuotationByID(w http.ResponseWriter, r *http.Request) {
	updateID := r.URL.Query().Get("id")

	row, err := s.db.GetRowByIDBuffer(updateID)
	if err != nil {
		log.Printf("Error getting quotation by ID: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if row.Value == 0 && row.UpdatedAt.IsZero() {
		log.Printf("ID %v doesn't exist", updateID)
		w.Write([]byte(fmt.Sprintf("ID %v doesn't exist", updateID)))
		return
	}

	resp := QuotationByIDResp{
		Value:    row.Value,
		UpdateAt: row.UpdatedAt,
	}

	res, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling response: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

type LatestResp struct {
	Value    float64   `json:"value"`
	UpdateAt time.Time `json:"updateAt"`
}

func (s *Server) getLatest(w http.ResponseWriter, r *http.Request) {
	cur := strings.ToTitle(r.URL.Query().Get("currency"))

	if !ValidateCurrency(cur) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No such currency"))
		return
	}

	row, err := s.db.GetLatest(cur)
	if err != nil {
		log.Printf("Error getting quotation by ID: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if row.Value == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Quotation value didn't set yet..."))
		return
	}

	resp := LatestResp{
		Value:    row.Value,
		UpdateAt: row.UpdatedAt,
	}

	res, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling response: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (s *Server) getRate(cur string) (float64, error) {
	//apiURI := "https://api.currencybeacon.com/v1/latest"

	qp := url.Values{}
	qp.Set("api_key", apiKey)
	qp.Set("symbols", cur)

	reqURL := fmt.Sprintf("%s?%s", apiURI, qp.Encode())

	resp, err := http.Get(reqURL)
	if err != nil {
		log.Printf("Error requesting currency API: %v\n", err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return 0, err
	}

	var curBeaconResp getRatesResp
	if err := json.Unmarshal(body, &curBeaconResp); err != nil {
		log.Printf("Error unmarshaling response: %v\n", err)
		return 0, err
	}
	rate := curBeaconResp.Rates[cur]
	if rate == 0 {
		err = errors.New("quotation rate mustn't be zero")
		log.Println(err)
		return 0, err
	}
	return rate, nil
}

func ValidateCurrency(currency string) bool {
	c := map[string]struct{}{
		eur: {},
		mxn: {},
		gel: {},
	}
	_, ok := c[currency]
	return ok
}
