package job

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/types"
)

type JobHandler func(ctx context.Context, job database.Job) error

type JobProcessor struct {
	db       *database.Database
	handlers map[string]JobHandler
}

func NewJobProcessor(db *database.Database) *JobProcessor {
	return &JobProcessor{
		db:       db,
		handlers: make(map[string]JobHandler),
	}
}

func (p *JobProcessor) RegisterHandler(jobType string, handler JobHandler) {
	p.handlers[jobType] = handler
}

func (p *JobProcessor) Start(workerCount int) {
	for range workerCount {
		go p.workerLoop()
	}
}

func (p *JobProcessor) workerLoop() {
	for {
		job, err := p.fetchNextJob()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		if job == nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		handler, ok := p.handlers[job.Type]
		if !ok {
			slog.Error("No handler for job type", "type", job.Type)
			p.markFailed(job, fmt.Errorf("no handler"))
			continue
		}

		ctx := context.Background()
		err = handler(ctx, *job)
		if err != nil {
			p.retryOrFail(job, err)
		} else {
			p.markSuccess(job)
		}
	}
}

func (p *JobProcessor) fetchNextJob() (*database.Job, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	job, err := tx.GetNextJob(context.Background())
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return nil, nil
		}

		return nil, err
	}

	err = tx.UpdateJob(context.Background(), job.Id, database.JobChanges{
		Status: database.Change[types.JobStatus]{
			Value:   types.JobStatusRunning,
			Changed: true,
		},
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (p *JobProcessor) retryOrFail(job *database.Job, jobErr error) {
	job.Attempts++
	if job.Attempts >= job.MaxAttempts {
		p.markFailed(job, jobErr)
		return
	}

	backoff := time.Duration(job.Attempts*job.Attempts) * time.Minute

	err := p.db.UpdateJob(context.Background(), job.Id, database.JobChanges{
		Status: database.Change[types.JobStatus]{
			Value:   types.JobStatusQueued,
			Changed: true,
		},
		RunAt: database.Change[int64]{
			Value:   time.Now().Add(backoff).UnixMilli(),
			Changed: true,
		},
		Attempts: database.Change[int]{
			Value:   job.Attempts,
			Changed: true,
		},
		Error: database.Change[sql.NullString]{
			Value: sql.NullString{
				String: jobErr.Error(),
				Valid:  true,
			},
			Changed: true,
		},
	})
	if err != nil {
		slog.Error("failed to update job retryOrFail", "id", job.Id, "err", err)
	}
}

func (p *JobProcessor) markSuccess(job *database.Job) {
	err := p.db.UpdateJob(context.Background(), job.Id, database.JobChanges{
		Status: database.Change[types.JobStatus]{
			Value:   types.JobStatusSuccess,
			Changed: true,
		},
	})
	if err != nil {
		slog.Error("failed to mark job success", "err", err)
	} else {
		slog.Info("job is marked success", "id", job.Id)
	}
}

func (p *JobProcessor) markFailed(job *database.Job, jobErr error) {
	err := p.db.UpdateJob(context.Background(), job.Id, database.JobChanges{
		Status: database.Change[types.JobStatus]{
			Value:   types.JobStatusFailed,
			Changed: true,
		},
		Error: database.Change[sql.NullString]{
			Value: sql.NullString{
				String: jobErr.Error(),
				Valid:  true,
			},
			Changed: true,
		},
	})
	if err != nil {
		slog.Error("failed to mark job failed", "err", err)
	} else {
		slog.Error("job is marked failed", "id", job.Id, "err", jobErr)
	}
}
