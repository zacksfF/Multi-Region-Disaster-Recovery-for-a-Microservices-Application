package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	domain "github.com/zacksfF/TempoScale-Distributed-Infrastructure-for-Time-Series-Data/internal/domain/summation"
)

func (r *ObservationSummationRepoImpl) updateBy(ctx context.Context, k *sq.And, m *domain.ObservationSummation) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	psql := sq.StatementBuilder.RunWith(r.dbCache).PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.
		Update("observation_summations").
		Set("start", m.Start).
		Set("finish", m.Finish).
		Set("day", m.Day).
		Set("week", m.Week).
		Set("month", m.Month).
		Set("year", m.Year).
		Set("frequency", m.Frequency).
		Set("result", m.Result).
		Where(k).
		ToSql()

	stmt, err := r.db.PrepareContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	return err
}

func (r *ObservationSummationRepoImpl) UpdateByPrimaryKey(ctx context.Context, oc *domain.ObservationSummation) error {
	k := &sq.And{
		sq.Eq{"entity_id": oc.EntityID},
		sq.Eq{"frequency": oc.Frequency},
		sq.Eq{"start": oc.Start},
		sq.Eq{"finish": oc.Finish},
	}
	return r.updateBy(ctx, k, oc)
}
