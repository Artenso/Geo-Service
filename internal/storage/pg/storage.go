package storage

import (
	"context"

	"github.com/Artenso/Geo-Service/internal/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

const (
	table = "users"

	idCol   = "id"
	nameCol = "name"
	passCol = "password"
)

type IStorage interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*model.User, error)
	GetByName(ctx context.Context, user *model.User) ([]*model.User, error)
}

// Storage ...
type storage struct {
	dbConn *pgx.Conn
}

// New creates new repository object
func NewStorage(dbConn *pgx.Conn) IStorage {
	return &storage{
		dbConn: dbConn,
	}
}

// Create adds user
func (r *storage) Create(ctx context.Context, user *model.User) error {
	builder := sq.Insert(table).
		Columns(nameCol, passCol).
		Values(user.Name, user.Pass).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build sql query")
	}

	rows, err := r.dbConn.Query(ctx, query, args...)
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

// GetByID gets users by id
func (r storage) GetByID(ctx context.Context, id int) (*model.User, error) {
	builder := sq.Select("*").
		From(table).
		Where(sq.Eq{
			idCol: id,
		}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build sql query")
	}

	user := new(model.User)

	rows, err := r.dbConn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if err = pgxscan.ScanAll(user, rows); err != nil {
		return nil, err
	}

	return user, nil
}

// List gets users from offset to limit
func (r storage) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	builder := sq.Select("*").
		From(table).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		GroupBy(idCol).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "failed to build sql query")
	}

	rows, err := r.dbConn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var users []*model.User

	if err = pgxscan.ScanAll(&users, rows); err != nil {
		return nil, err
	}

	return users, nil
}

// Update updates users password
func (r storage) Update(ctx context.Context, user *model.User) error {
	builder := sq.Update(table).
		Where(sq.Eq{
			nameCol: user.Name,
		}).
		PlaceholderFormat(sq.Dollar).
		Set(passCol, user.Pass)

	query, args, err := builder.ToSql()

	if err != nil {
		return errors.Wrap(err, "failed to build sql query")
	}

	rows, err := r.dbConn.Query(ctx, query, args...)
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

// Delete deletes user from database
func (r storage) Delete(ctx context.Context, id int) error {
	builder := sq.Delete(table).
		Where(sq.Eq{
			idCol: id,
		}).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build sql query")
	}

	rows, err := r.dbConn.Query(ctx, query, args...)
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

func (r storage) GetByName(ctx context.Context, user *model.User) ([]*model.User, error) {
	builder := sq.Select(nameCol, passCol).
		From(table).
		Where(sq.Eq{
			nameCol: user.Name,
		}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build sql query")
	}

	rows, err := r.dbConn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var users []*model.User

	if err = pgxscan.ScanAll(&users, rows); err != nil {
		return nil, err
	}

	return users, nil
}
