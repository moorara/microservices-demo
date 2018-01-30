package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	voteTTLInMin = 60
)

type (
	// Vote represents a vote for a link
	Vote struct {
		ID     string `json:"id"`
		LinkID string `json:"linkId"`
		Stars  int    `json:"stars"`
	}

	// VoteManager abstracts CRUD operations for Vote
	VoteManager interface {
		Create(ctx context.Context, linkID string, stars int) (*Vote, error)
		Get(ctx context.Context, id string) (*Vote, error)
	}

	redisVoteManager struct {
		persister Persister
		voteTTL   time.Duration
	}
)

// NewVoteManager creates a new vote manager
func NewVoteManager(persister Persister) VoteManager {
	return &redisVoteManager{
		persister: persister,
		voteTTL:   voteTTLInMin * time.Minute,
	}
}

func (sm *redisVoteManager) Create(ctx context.Context, linkID string, stars int) (*Vote, error) {
	vote := &Vote{
		ID:     uuid.New().String(),
		LinkID: linkID,
		Stars:  stars,
	}

	data, err := json.Marshal(vote)
	if err != nil {
		return nil, err
	}

	chErr := make(chan error, 1)
	go func() {
		chErr <- sm.persister.Save(vote.ID, data, sm.voteTTL)
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-chErr:
	}

	if err != nil {
		return nil, errors.Wrap(err, "Vote creation failed")
	}
	return vote, nil
}

func (sm *redisVoteManager) Get(ctx context.Context, id string) (*Vote, error) {
	var err error
	vote := new(Vote)

	chErr := make(chan error, 1)
	go func() {
		data, err := sm.persister.Load(id)
		if err != nil {
			chErr <- err
			return
		}
		chErr <- json.Unmarshal(data, vote)
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-chErr:
	}

	if err != nil {
		return nil, errors.Wrap(err, "Vote retrieval failed")
	}
	return vote, nil
}
