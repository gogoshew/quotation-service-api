package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Server struct {
	ctx    context.Context
	router *mux.Router
	db     *gorm.DB
	http.Server
}

func NewServer(ctx context.Context, router *mux.Router, db *gorm.DB) (s *Server) {
	s = &Server{
		ctx:    ctx,
		router: router,
		db:     db,
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
}

type updateResp struct {
}

func (s *Server) updateQuotation(w http.ResponseWriter, r *http.Request) {
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

	var respBody getRatesResp
	if err := json.Unmarshal(body, &respBody); err != nil {
		log.Fatalf("Error unmarshaling response: %v\n", err)
	}

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
