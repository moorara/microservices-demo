package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/moorara/microservices-demo/go-service/service"
)

type (
	// VoteHandler provides http handlers for Vote
	VoteHandler interface {
		PostVote(w http.ResponseWriter, r *http.Request)
		GetVotes(w http.ResponseWriter, r *http.Request)
		GetVote(w http.ResponseWriter, r *http.Request)
		DeleteVote(w http.ResponseWriter, r *http.Request)
	}

	postgresVoteHandler struct {
		vm     service.VoteManager
		logger log.Logger
	}
)

// NewVoteHandler creates a new vote handler
func NewVoteHandler(db service.DB, logger log.Logger) VoteHandler {
	return &postgresVoteHandler{
		vm:     service.NewVoteManager(db, logger),
		logger: logger,
	}
}

func (vh *postgresVoteHandler) PostVote(w http.ResponseWriter, r *http.Request) {
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

	vote, err := vh.vm.Create(context.Background(), req.LinkID, req.Stars)
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

func (vh *postgresVoteHandler) GetVotes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkID := vars["linkId"]

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if linkID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	votes, err := vh.vm.GetAll(context.Background(), linkID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(votes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (vh *postgresVoteHandler) GetVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	voteID := vars["id"]

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if voteID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vote, err := vh.vm.Get(context.Background(), voteID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if vote == nil {
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

func (vh *postgresVoteHandler) DeleteVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	voteID := vars["id"]

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if voteID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := vh.vm.Delete(context.Background(), voteID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
