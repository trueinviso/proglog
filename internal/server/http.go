package server

import (
	"encoding/json"
	"net/http"
)

type MemoryStore struct {
	Log *Log
}

type WriteResponse struct {
	Offset uint64 `json:"offset"`
}

type WriteRequest struct {
	Record Record `json:"record"`
}

type ReadRequest struct {
	Offset uint64 `json:"offset"`
}

type ReadResponse struct {
	Record Record `json:"record"`
}

func NewHTTPServer(address string) *http.Server {
	store := NewStore()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", store.Write)
	mux.HandleFunc("GET /", store.Read)

	return &http.Server{
		Addr:    address,
		Handler: mux,
	}
}

func NewStore() *MemoryStore {
	return &MemoryStore{
		Log: &Log{},
	}
}

func (s *MemoryStore) Write(w http.ResponseWriter, r *http.Request) {
	var wReq WriteRequest
	err := json.NewDecoder(r.Body).Decode(&wReq)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	offset, err := s.Log.Append(wReq.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(WriteResponse{offset})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *MemoryStore) Read(w http.ResponseWriter, r *http.Request) {
	var rReq ReadRequest
	err := json.NewDecoder(r.Body).Decode(&rReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := s.Log.Read(rReq.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(ReadResponse{record})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
