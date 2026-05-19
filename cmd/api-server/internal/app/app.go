package app

import (
	"context"
	"curs/cmd/api-server/internal/models"

	"github.com/google/uuid"
)

type store interface {
	GetList(context context.Context, listId string) (models.List, error)
	CreateList(context context.Context, list models.List) (error)
}

type App struct {
	store store
}

func New(store store) *App {
	return &App{store: store}
}

func (a *App) GetList(context context.Context, listId string) (models.List, error) {
	return a.GetList(context, listId)
}

func (a *App) CreateList(context context.Context, listName string) (models.List, error) {
	listID := uuid.NewString();

	list := models.List{
		Id: listID,
		Name: listName,
	}

	return list, a.store.CreateList(context, list)
}
