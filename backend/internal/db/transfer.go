package db

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/mlemb/wallet/backend/internal/model"
)

var ErrTransferNotFound = errors.New("transfer not found")

type TransfersParams struct {
	Type       string
	From       string
	To         string
	AmountFrom float64
	AmountTo   float64
	TimeFrom   time.Time
	TimeTo     time.Time
	Page       int
	PageSize   int
}

func CountTransfers(ctx context.Context, params *TransfersParams) (int64, error) {
	builder := sqlbuilder.NewSelectBuilder().
		Select(`COUNT("transfers"."id")`).
		From(`"transfers"`)

	buildSelectTransferQuery(builder, params)

	query, args := builder.Build()
	row := db.QueryRowContext(ctx, query, args...)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func QueryTransfers(ctx context.Context, params *TransfersParams) ([]*model.Transfer, error) {
	builder := sqlbuilder.NewSelectBuilder().
		Select(`"transfers"."id"`, `"transfers"."type"`, `"transfers"."from"`, `"transfers"."to"`, `"transfers"."amount"`, `"transfers"."time"`).
		From(`"transfers"`).
		OrderBy(`"transfers"."time" DESC`)

	buildSelectTransferQuery(builder, params)

	query, args := builder.Build()
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []*model.Transfer
	for rows.Next() {
		var transfer model.Transfer
		err := rows.Scan(&transfer.ID, &transfer.Type, &transfer.From, &transfer.To, &transfer.Amount, &transfer.Time)
		if err != nil {
			return nil, err
		}
		transfers = append(transfers, &transfer)
	}

	return transfers, nil
}

func buildSelectTransferQuery(builder *sqlbuilder.SelectBuilder, params *TransfersParams) {
	if params.Type != "" {
		builder.Where(builder.Equal(`"transfers"."type"`, params.Type))
	}

	if params.From != "" {
		builder.Where(builder.Like(`"transfers"."from"`, "%"+strings.ToLower(params.From)+"%"))
	}

	if params.To != "" {
		builder.Where(builder.Like(`"transfers"."to"`, "%"+strings.ToLower(params.To)+"%"))
	}

	if params.AmountFrom != 0 {
		builder.Where(builder.GreaterEqualThan(`"transfers"."amount"`, params.AmountFrom))
	}

	if params.AmountTo != 0 {
		builder.Where(builder.LessEqualThan(`"transfers"."amount"`, params.AmountTo))
	}

	if !params.TimeFrom.IsZero() {
		builder.Where(builder.GreaterEqualThan(`"transfers"."time"`, params.TimeFrom))
	}

	if !params.TimeTo.IsZero() {
		builder.Where(builder.LessEqualThan(`"transfers"."time"`, params.TimeTo))
	}

	if params.Page > 1 {
		builder.Offset((params.Page - 1) * params.PageSize)
	}

	if params.PageSize != 0 {
		builder.Limit(params.PageSize)
	}
}

func QueryTransferByID(ctx context.Context, id string) (*model.Transfer, error) {
	builder := sqlbuilder.NewSelectBuilder().
		Select(`"transfers"."id"`, `"transfers"."type"`, `"transfers"."from"`, `"transfers"."to"`, `"transfers"."amount"`, `"transfers"."time"`).
		From(`"transfers"`)

	builder.Where(builder.Equal(`"transfers"."id"`, id))

	query, args := builder.Build()
	row := db.QueryRowContext(ctx, query, args...)

	var transfer model.Transfer
	err := row.Scan(&transfer.ID, &transfer.Type, &transfer.From, &transfer.To, &transfer.Amount, &transfer.Time)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTransferNotFound
	}
	if err != nil {
		return nil, err
	}

	return &transfer, nil
}

func CreateTransfer(ctx context.Context, transfer *model.Transfer) error {
	builder := sqlbuilder.NewInsertBuilder().
		InsertInto(`"transfers"`).
		Cols(`"type"`, `"from"`, `"to"`, `"amount"`, `"time"`).
		Values(transfer.Type, transfer.From, transfer.To, transfer.Amount, transfer.Time)

	query, args := builder.Build()
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	transfer.ID = id

	return nil
}

func UpdateTransfer(ctx context.Context, transfer *model.Transfer) error {
	builder := sqlbuilder.NewUpdateBuilder().
		Update(`"transfers"`)

	builder.Set(
		builder.Assign(`"type"`, transfer.Type),
		builder.Assign(`"from"`, transfer.From),
		builder.Assign(`"to"`, transfer.To),
		builder.Assign(`"amount"`, transfer.Amount),
		builder.Assign(`"time"`, transfer.Time),
	)

	builder.Where(builder.Equal(`"transfers"."id"`, transfer.ID))

	query, args := builder.Build()
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrTransferNotFound
	}

	return nil
}

func DeleteTransferByID(ctx context.Context, id int64) error {
	builder := sqlbuilder.NewDeleteBuilder().
		DeleteFrom(`"transfers"`)

	builder.Where(builder.Equal(`"transfers"."id"`, id))

	query, args := builder.Build()
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrTransferNotFound
	}

	return nil
}
