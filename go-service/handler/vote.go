package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/moorara/toys/microservices/go-service/service"
)

// VoteHandler provides http handlers for Vote
type VoteHandler interface {
	PostVote(w http.ResponseWriter, r *http.Request)
	GetVote(w http.ResponseWriter, r *http.Request)
}

// redisVoteHandler implements VoteHandler interface
type redisVoteHandler struct {
	sm     service.VoteManager
	logger log.Logger
}

// NewVoteHandler creates a new vote handler
func NewVoteHandler(persister service.Persister, logger log.Logger) VoteHandler {
	return &redisVoteHandler{
		sm:     service.NewVoteManager(persister),
		logger: logger,
	}
}

func (sh *redisVoteHandler) PostVote(w http.ResponseWriter, r *http.Request) {
	req := struct {
		LinkID string `json:"linkId"`
		Stars  int    `json:"stars"`
	}{}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.LinkID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vote, err := sh.sm.Create(context.Background(), req.LinkID, req.Stars)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(vote)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (sh *redisVoteHandler) GetVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	voteID := vars["id"]

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if voteID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vote, err := sh.sm.Get(context.Background(), voteID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(vote)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
