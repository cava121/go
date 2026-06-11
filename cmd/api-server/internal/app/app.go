package app

import (
	"context"
	"curs/cmd/api-server/internal/models"

	"github.com/google/uuid"
)

type store interface {
	GetListById(context context.Context, listId string) (models.List, error)
	CreateList(context context.Context, list models.List) error
	GetLists(context context.Context) ([]models.List, error)
}

type App struct {
	store store
}

func New(store store) *App {
	return &App{store: store}
}

func (a *App) GetLists(context context.Context) ([]models.List, error) {
	return a.store.GetLists(context)
}

func (a *App) GetListById(context context.Context, listId string) (models.List, error) {
	return a.store.GetListById(context, listId)
}

func (a *App) CreateList(context context.Context, listName string) (models.List, error) {
	listID := uuid.NewString();

	list := models.List{
		Id: listID,
		Name: listName,
	}

	return list, a.store.CreateList(context, list)
}
