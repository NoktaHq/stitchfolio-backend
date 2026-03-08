package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/util"
)

type EntityDocumentRepository interface {
	Create(*context.Context, *entities.EntityDocuments) *errs.XError
	Update(*context.Context, *entities.EntityDocuments) *errs.XError
	Get(*context.Context, uint) (*entities.EntityDocuments, *errs.XError)
	Delete(*context.Context, uint) *errs.XError

	GetAllEntityDocuments(*context.Context) ([]entities.EntityDocuments, *errs.XError)
	GetEntityDocumentsByEntity(*context.Context, uint, string, string) ([]entities.EntityDocuments, *errs.XError)
}

type entityDocumentRepository struct {
	GormDAL
}

func ProvideEntityDocumentRepository(dal GormDAL) EntityDocumentRepository {
	return &entityDocumentRepository{GormDAL: dal}
}

func (repo *entityDocumentRepository) Create(ctx *context.Context, entityDocuments *entities.EntityDocuments) *errs.XError {
	res := repo.WithDB(ctx).Create(&entityDocuments)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save entity documents", res.Error)
	}
	return nil
}

func (repo *entityDocumentRepository) Update(ctx *context.Context, entityDocuments *entities.EntityDocuments) *errs.XError {
	return repo.GormDAL.Update(ctx, *entityDocuments)
}

func (repo *entityDocumentRepository) Get(ctx *context.Context, id uint) (*entities.EntityDocuments, *errs.XError) {
	entityDocuments := entities.EntityDocuments{}
	res := repo.WithDB(ctx).
		Scopes(scopes.Channel()).
		Find(&entityDocuments, id)
	if res.Error != nil || res.RowsAffected == 0 {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find entity documents", res.Error)
	}
	return &entityDocuments, nil
}

func (repo *entityDocumentRepository) GetAllEntityDocuments(ctx *context.Context) ([]entities.EntityDocuments, *errs.XError) {
	entityDocuments := new([]entities.EntityDocuments)

	res := repo.WithDB(ctx).
		Scopes(scopes.IsActive()).
		Scopes(scopes.Channel()).
		Scopes(db.Paginate(ctx)).
		Find(entityDocuments)

	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to fetch all entity documents", res.Error)
	}

	return *entityDocuments, nil
}

func (repo *entityDocumentRepository) GetEntityDocumentsByEntity(ctx *context.Context, entityId uint, entityName, typ string) ([]entities.EntityDocuments, *errs.XError) {

	entityDocuments := make([]entities.EntityDocuments, 0)

	query := repo.WithDB(ctx).
		Scopes(scopes.IsActive()).
		Scopes(scopes.Channel()).
		Scopes(db.Paginate(ctx)).
		Where("entity_name = ? AND entity_id = ? ", entityName, entityId)

	// Usually for an entity we load all the documents present there to be displayed
	// In that case we would not need the type
	// only when we are specifically querying for a type of document typ will be used
	if !util.IsNilOrEmptyString(&typ) {
		query = query.Where("type = ?", typ)
	}

	res := query.Find(&entityDocuments)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to fetch entity documents by entity", res.Error)
	}

	return entityDocuments, nil
}

func (repo *entityDocumentRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	entityDocuments := &entities.EntityDocuments{Model: &entities.Model{ID: id, IsActive: false}}
	return repo.GormDAL.Delete(ctx, entityDocuments)
}
