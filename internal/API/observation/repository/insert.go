package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"

	domain "github.com/zacksfF/TempoScale-Distributed-Infrastructure-for-Time-Series-Data/internal/domain/observation"
)

func (r *ObservationRepoImpl) Insert(ctx context.Context, m *domain.Observation) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	psql := sq.StatementBuilder.RunWith(r.dbCache).PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.
		Insert("observations").
		Columns(
			"entity_id",
			"meta",
			"timestamp",
			"value",
		).
		Values(
			m.EntityID,
			m.Meta,
			m.Timestamp,
			m.Value,
		).
		ToSql()
	if err != nil {
		r.logger.Error().Err(err).Caller().
			Uint64("entity_id", m.EntityID).
			Str("meta", m.Meta).
			Time("timestamp", m.Timestamp).
			Float64("value", m.Value).
			Msg("failed creating query")
		return err
	}

	stmt, err := r.db.PrepareContext(ctx, sql)
	if err != nil {
		r.logger.Error().Err(err).Caller().
			Uint64("entity_id", m.EntityID).
			Str("meta", m.Meta).
			Time("timestamp", m.Timestamp).
			Float64("value", m.Value).
			Msg("prepare context error")
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		r.logger.Error().Err(err).Caller().
			Uint64("entity_id", m.EntityID).
			Str("meta", m.Meta).
			Time("timestamp", m.Timestamp).
			Float64("value", m.Value).
			Msg("exec context error")
	}
	return err
}
