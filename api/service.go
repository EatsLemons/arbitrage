package api

import (
	"arbitrage/exchange"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Service struct {
	Exchange *exchange.Exchange

	httpServer *http.Server
}

func (s *Service) Run(port int) {
	log.Printf("[INFO] server started at :%d", port)

	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1/").Subrouter()

	v1.HandleFunc("/ping", s.pingHndlr).Methods("GET")
	v1.HandleFunc("/arbitrage", s.arbitrageHndlr).Methods("GET")

	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	err := s.httpServer.ListenAndServe()
	log.Printf("[WARN] http server terminated, %s", err)
}

func (s *Service) arbitrageHndlr(w http.ResponseWriter, r *http.Request) {
	s.Exchange.UpdatePrices()

	s.makeJSONResponse(w, s.Exchange.FindProfitCurrPairs())
}

func (s *Service) pingHndlr(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong"))
}

func (s *Service) makeJSONResponse(w http.ResponseWriter, response interface{}) {
	jsonResponse, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
