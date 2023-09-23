package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/raja-dettex/dap/dap"
)

type Handler struct {
	dap    *dap.Dap
	router *mux.Router
}

func NewHandler(dbName string, r *mux.Router) (*Handler, error) {
	if dbName == "" {
		return nil, errors.New("invalid dbName")
	}
	dapDB, err := dap.New(dbName)
	if err != nil {
		return nil, err
	}
	return &Handler{dap: dapDB, router: r}, nil
}

func (h *Handler) handleInsert(w http.ResponseWriter, r *http.Request) {
	var response = make(map[string]string)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response["error"] = "method except post not allowed"
		json.NewEncoder(w).Encode(response)
		return
	}
	collection := mux.Vars(r)["collection"]
	key := mux.Vars(r)["key"]
	var inputData map[string]any
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&inputData)
	err := h.dap.Insert(collection, key, inputData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response["error"] = fmt.Sprintf("can not write %d", err)
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response["status"] = "write successful"
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleSelect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method except Get not allowed"))
	}
	collection := mux.Vars(r)["collection"]
	key := mux.Vars(r)["key"]
	data, err := h.dap.Select(collection, key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("can not write %d", err)))
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
