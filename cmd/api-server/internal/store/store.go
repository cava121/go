package store

import (
	"context"
	"curs/cmd/api-server/internal/models"

	"github.com/jackc/pgx/v5"
)

type Store struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *Store {
	return &Store{conn: conn}
}

func (s *Store) CreateList(context context.Context, list models.List) (error) {
	_, err := s.conn.Exec(context, `INSERT INTO lists(id, name) VALUES($1, $2)`, list.Id, list.Name); {
		return err
	}
}

func (s *Store) GetList(context context.Context, listId string) (models.List, error) {
	var name string;

	if err := s.conn.QueryRow(context, `SELECT name FROM lists WHERE id = $1`, listId).Scan(&name); err != nil {
		return models.List{}, err
	}

	return models.List{
		Id: listId,
		Name: name,
	}, nil
}
