package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type Job struct {
	RowId int `db:"rowid"`

	Id   string `db:"id"`
	Type string `db:"type"`

	Status   types.JobStatus `db:"status"`
	Priority int             `db:"status"`
	RunAt    int64           `db:"run_at"`

	Attempts    int `db:"attempts"`
	MaxAttempts int `db:"max_attempts"`

	Payload string         `db:"payload"`
	Error   sql.NullString `db:"error"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
func JobQuery() *goqu.SelectDataset {
	query := dialect.From("jobs").
		Select(
			"jobs.rowid",

			"jobs.id",
			"jobs.type",

			"jobs.status",
			"jobs.status",
			"jobs.run_at",

			"jobs.attempts",
			"jobs.max_attempts",

			"jobs.payload",
			"jobs.error",

			"jobs.created",
			"jobs.updated",
		)

	return query
}

func (db DB) GetAllJobs(ctx context.Context) ([]Job, error) {
	query := JobQuery()
	return ember.Multiple[Job](db.db, ctx, query)
}

func (db DB) GetJobById(ctx context.Context, id string) (Job, error) {
	query := JobQuery().
		Where(goqu.I("jobs.id").Eq(id))

	return ember.Single[Job](db.db, ctx, query)
}

func (db DB) GetNextJob(ctx context.Context) (Job, error) {
    // row := tx.QueryRow(`
    //     SELECT id, type, payload, status, priority, run_at, attempts, max_attempts, error
    //     FROM jobs
    //     WHERE status = 'queued' AND run_at <= CURRENT_TIMESTAMP
    //     ORDER BY priority DESC, run_at ASC, created_at ASC
    //     LIMIT 1
    //     `)
	query := JobQuery().
		Where(
			goqu.I("jobs.status").Eq(string(types.JobStatusQueued)),
			goqu.I("jobs.run_at").Lte(time.Now().UnixMilli()),
		).
		Order(
			goqu.I("jobs.priority").Desc(),
			goqu.I("jobs.run_at").Asc(),
			goqu.I("jobs.created").Asc(),
		)

	return ember.Single[Job](db.db, ctx, query)
}

type CreateJobParams struct {
	Id   string
	Type string

	Status   types.JobStatus
	Priority int
	RunAt    int64

	Attempts    int
	MaxAttempts int

	Payload string
	Error   sql.NullString

	Created int64
	Updated int64
}

func (db DB) CreateJob(ctx context.Context, params CreateJobParams) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = utils.CreateJobId()
	}

	if params.Type == "" {
		params.Type = "unknown"
	}

	if params.Status == "" {
		params.Status = types.JobStatusQueued
	}

	query := dialect.Insert("jobs").Rows(goqu.Record{
		"id":   params.Id,
		"type": params.Type,

		"status":   params.Status,
		"priority": params.Priority,
		"run_at":   params.RunAt,

		"attempts":     params.Attempts,
		"max_attempts": params.MaxAttempts,

		"payload": params.Payload,
		"error":   params.Error,

		"created": params.Created,
		"updated": params.Updated,
	}).
		Returning("id")

	return ember.Single[string](db.db, ctx, query)
}

type JobChanges struct {
	Type Change[string]

	Status   Change[types.JobStatus]
	Priority Change[int]
	RunAt    Change[int64]

	Attempts    Change[int]
	MaxAttempts Change[int]

	Payload Change[string]
	Error   Change[sql.NullString]

	Created Change[int64]
}

func (db DB) UpdateJob(ctx context.Context, id string, changes JobChanges) error {
	record := goqu.Record{}

	addToRecord(record, "type", changes.Type)

	addToRecord(record, "status", changes.Status)
	addToRecord(record, "priority", changes.Priority)
	addToRecord(record, "run_at", changes.RunAt)

	addToRecord(record, "attempts", changes.Attempts)
	addToRecord(record, "max_attempts", changes.MaxAttempts)

	addToRecord(record, "payload", changes.Payload)
	addToRecord(record, "error", changes.Error)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("jobs").
		Set(record).
		Where(goqu.I("jobs.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) RemoveJob(ctx context.Context, id string) error {
	query := dialect.Delete("jobs").
		Where(goqu.I("jobs.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
