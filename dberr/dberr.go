package dberr

import (
	"errors"
	"slices"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
)

// IsPgError checks whether err is a *pgconn.PgError or *pq.Error with one of the given Postgres SQLSTATE codes.
func IsPgError(err error, codes ...string) bool {
	if len(codes) == 0 {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return slices.Contains(codes, string(pgErr.Code))
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return slices.Contains(codes, string(pqErr.Code))
	}

	return false
}

// IsUniqueViolation returns true if the error is a UNIQUE constraint violation (SQLSTATE 23505).
func IsUniqueViolation(err error) bool {
	return IsPgError(err, pgerrcode.UniqueViolation)
}

// IsForeignKeyViolation returns true if the error is a foreign key violation (SQLSTATE 23503).
func IsForeignKeyViolation(err error) bool {
	return IsPgError(err, pgerrcode.ForeignKeyViolation)
}

// IsNotNullViolation returns true if the error is a NOT NULL constraint violation (SQLSTATE 23502).
func IsNotNullViolation(err error) bool {
	return IsPgError(err, pgerrcode.NotNullViolation)
}

// IsCheckViolation returns true if the error is a CHECK constraint violation (SQLSTATE 23514).
func IsCheckViolation(err error) bool {
	return IsPgError(err, pgerrcode.CheckViolation)
}

// IsExclusionViolation returns true if the error is an exclusion constraint violation (SQLSTATE 23P01).
func IsExclusionViolation(err error) bool {
	return IsPgError(err, pgerrcode.ExclusionViolation)
}

// IsSerializationFailure returns true if the error is a transaction serialization failure (SQLSTATE 40001).
func IsSerializationFailure(err error) bool {
	return IsPgError(err, pgerrcode.SerializationFailure)
}

// IsDeadlockDetected returns true if the error is a deadlock detected (SQLSTATE 40P01).
func IsDeadlockDetected(err error) bool {
	return IsPgError(err, pgerrcode.DeadlockDetected)
}

// IsNoRows returns true if the error is pgx.ErrNoRows.
func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
