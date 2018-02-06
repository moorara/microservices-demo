package service

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
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
		GetAll(ctx context.Context, linkID string) ([]Vote, error)
		Get(ctx context.Context, id string) (*Vote, error)
		Delete(ctx context.Context, id string) error
	}

	postgresVoteManager struct {
		db     DB
		logger log.Logger
	}
)

// NewVoteManager creates a new vote manager
func NewVoteManager(db DB, logger log.Logger) VoteManager {
	return &postgresVoteManager{
		db:     db,
		logger: logger,
	}
}

func (vm *postgresVoteManager) Create(ctx context.Context, linkID string, stars int) (*Vote, error) {
	vote := &Vote{
		ID:     uuid.New().String(),
		LinkID: linkID,
		Stars:  stars,
	}

	query := `INSERT INTO votes (id, link_id, stars) VALUES ($1, $2, $3)`
	_, err := vm.db.ExecContext(ctx, query, vote.ID, vote.LinkID, vote.Stars)
	if err != nil {
		level.Error(vm.logger).Log("message", err.Error())
		return nil, err
	}

	return vote, nil
}

func (vm *postgresVoteManager) GetAll(ctx context.Context, linkID string) ([]Vote, error) {
	votes := make([]Vote, 0)

	query := `SELECT * FROM votes WHERE link_id=$1`
	rows, err := vm.db.QueryContext(ctx, query, linkID)
	if err != nil {
		level.Error(vm.logger).Log("message", err.Error())
		return nil, err
	}

	for rows.Next() {
		vote := Vote{}
		err = rows.Scan(&vote.ID, &vote.LinkID, &vote.Stars)
		if err == nil {
			votes = append(votes, vote)
		}
	}

	return votes, nil
}

func (vm *postgresVoteManager) Get(ctx context.Context, id string) (*Vote, error) {
	vote := new(Vote)

	query := `SELECT * FROM votes WHERE id=$1`
	row := vm.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&vote.ID, &vote.LinkID, &vote.Stars)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		level.Error(vm.logger).Log("message", err.Error())
		return nil, err
	}

	return vote, nil
}

func (vm *postgresVoteManager) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM votes WHERE id=$1`
	_, err := vm.db.ExecContext(ctx, query, id)
	if err != nil {
		level.Error(vm.logger).Log("message", err.Error())
		return err
	}

	return nil
}
