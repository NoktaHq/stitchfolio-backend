package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/errs"
)

type FileStoreRepository interface {
	Create(*context.Context, *entities.FileStoreMetadata) (uint, *errs.XError)
	Update(*context.Context, *entities.FileStoreMetadata) *errs.XError
	Get(*context.Context, uint) (*entities.FileStoreMetadata, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
	GetByKey(*context.Context, string) (*entities.FileStoreMetadata, *errs.XError)
}

type fileStoreRepository struct {
	GormDAL
}

func ProvideFileStoreRepository(dal GormDAL) FileStoreRepository {
	return &fileStoreRepository{GormDAL: dal}
}

func (ur *fileStoreRepository) Create(ctx *context.Context, fileStoreMetadata *entities.FileStoreMetadata) (uint, *errs.XError) {
	res := ur.WithDB(ctx).Create(&fileStoreMetadata)
	if res.Error != nil {
		return 0, errs.NewXError(errs.DATABASE, "Unable to save file store metadata", res.Error)
	}
	return fileStoreMetadata.ID, nil
}

func (ur *fileStoreRepository) Update(ctx *context.Context, fileStoreMetadata *entities.FileStoreMetadata) *errs.XError {
	return ur.GormDAL.Update(ctx, *fileStoreMetadata)
}

func (repo *fileStoreRepository) Get(ctx *context.Context, id uint) (*entities.FileStoreMetadata, *errs.XError) {
	fileStoreMetadata := entities.FileStoreMetadata{}
	res := repo.WithDB(ctx).Scopes(scopes.Channel()).
		Find(&fileStoreMetadata, id)

	if res.Error != nil || res.RowsAffected == 0 {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find file store metadata", res.Error)
	}
	return &fileStoreMetadata, nil
}

func (repo *fileStoreRepository) GetByKey(ctx *context.Context, fileKey string) (*entities.FileStoreMetadata, *errs.XError) {
	fileStoreMetadata := entities.FileStoreMetadata{}

	res := repo.WithDB(ctx).
		Scopes(scopes.Channel()).
		Scopes(scopes.IsActive()).
		Where("file_key = ?", fileKey).
		Find(&fileStoreMetadata)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find file store metadata", res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, errs.NewXError(errs.NOT_EXIST, "Unable to find file store metadata", res.Error)
	}
	return &fileStoreMetadata, nil
}

func (repo *fileStoreRepository) Delete(ctx *context.Context, id uint) *errs.XError {
	fileStoreMetadata := &entities.FileStoreMetadata{Model: &entities.Model{ID: id, IsActive: false}}
	return repo.GormDAL.Delete(ctx, fileStoreMetadata)
}
