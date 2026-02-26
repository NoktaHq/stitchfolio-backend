package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
)

type MasterConfigRepository interface {
	Create(*context.Context, *entities.MasterConfig) *errs.XError
	Update(*context.Context, *entities.MasterConfig) *errs.XError
	Get(*context.Context, uint) (*entities.MasterConfig, *errs.XError)
	GetValue(*context.Context, string, string) (*entities.MasterConfig, *errs.XError)
	LoadAll(*context.Context) ([]entities.MasterConfig, *errs.XError)
	GetForBrowse(*context.Context, string) ([]entities.MasterConfig, *errs.XError)
}

type masterConfigRepository struct {
	GormDAL
}

func ProvideMasterConfigRepository(customDB GormDAL) MasterConfigRepository {
	return &masterConfigRepository{GormDAL: customDB}
}

func (repo *masterConfigRepository) Create(ctx *context.Context, config *entities.MasterConfig) *errs.XError {
	res := repo.WithDB(ctx).Create(config)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save master config", res.Error)
	}
	return nil
}

func (repo *masterConfigRepository) Update(ctx *context.Context, config *entities.MasterConfig) *errs.XError {
	return repo.GormDAL.Update(ctx, *config)
}

func (repo *masterConfigRepository) Get(ctx *context.Context, id uint) (*entities.MasterConfig, *errs.XError) {
	config := entities.MasterConfig{}
	res := repo.WithDB(ctx).Model(config).
		Scopes(scopes.WithAuditInfo()).
		Scopes(scopes.Channel()).
		Find(&config, id)

	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find master config", res.Error)
	}
	return &config, nil
}

// GetByNameAndType retrieves a master config by name and type
func (repo *masterConfigRepository) GetValue(ctx *context.Context, keyType string, name string) (*entities.MasterConfig, *errs.XError) {

	config := entities.MasterConfig{}
	res := repo.WithDB(ctx).
		Scopes(scopes.Channel()).
		Where("type = ? and name = ?", keyType, name).
		First(&config)

	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find master config", res.Error)
	}
	return &config, nil
}

func (repo *masterConfigRepository) LoadAll(ctx *context.Context) ([]entities.MasterConfig, *errs.XError) {

	var configs []entities.MasterConfig
	res := repo.WithDB(ctx).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Find(&configs)

	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find master configs", res.Error)
	}
	return configs, nil
}

func (repo *masterConfigRepository) GetForBrowse(ctx *context.Context, search string) ([]entities.MasterConfig, *errs.XError) {
	var configs []entities.MasterConfig
	res := repo.WithDB(ctx).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.ILike(search, "name", "type")).
		Scopes(db.Paginate(ctx)).
		Find(&configs)

	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find master configs", res.Error)
	}
	return configs, nil
}
