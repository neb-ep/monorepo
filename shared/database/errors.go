package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type ConstraintMapper map[string]error

type ErrorMapper struct {
	mapper ConstraintMapper
}

func NewErrorDescriber(mapper ConstraintMapper) *ErrorMapper {
	return &ErrorMapper{
		mapper: mapper,
	}
}

func (d *ErrorMapper) Describe(in error) (out error) {
	var pgerr *pgconn.PgError

	if !errors.As(in, &pgerr) {
		return in
	}

	out, found := d.mapper[pgerr.ConstraintName]

	if !found {
		return in
	}

	return out
}
