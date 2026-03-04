package service

import (
	"context"
	"fmt"
	"time"

	"github.com/imkarthi24/sf-backend/internal/config"
	"github.com/imkarthi24/sf-backend/internal/mapper"
	"github.com/imkarthi24/sf-backend/internal/model/models"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/storage"
)

var defaultContentType string = "application/octet-stream"

type FileStoreService interface {
	SaveFileStoreMetadata(*context.Context, requestModel.FileStoreMetadata) (uint, *errs.XError)
	UpdateFileStoreMetadata(*context.Context, requestModel.FileStoreMetadata, uint) *errs.XError
	GetFileStoreMetadata(*context.Context, uint) (*responseModel.FileStoreMetadata, *errs.XError)
	DeleteFileStoreMetadata(*context.Context, uint) *errs.XError

	GetFileStoreMetadataByKey(ctx *context.Context, entityName string, entityId uint, kind string) (*responseModel.FileStoreMetadata, *errs.XError)
	GetFileKey(ctx *context.Context, entityName string, entityId uint, fileType string) string
	Upload(ctx *context.Context, file models.FileUpload) *errs.XError
	GetFileMetadataIfExists(ctx *context.Context, entityType string, id uint, kind string) (bool, *responseModel.FileStoreMetadata, *errs.XError)
}

type fileStoreService struct {
	fileStoreRepo repository.FileStoreRepository
	mapper        mapper.Mapper
	config        config.AppConfig
	respMapper    mapper.ResponseMapper
	s3Storage     storage.CloudStorageProvider
}

func ProvideFileStoreService(repo repository.FileStoreRepository, s3Storage storage.CloudStorageProvider, mapper mapper.Mapper, config config.AppConfig, respMapper mapper.ResponseMapper) FileStoreService {
	return fileStoreService{
		fileStoreRepo: repo,
		mapper:        mapper,
		config:        config,
		respMapper:    respMapper,
		s3Storage:     s3Storage,
	}
}

func (svc fileStoreService) SaveFileStoreMetadata(ctx *context.Context, fileStoreMetadata requestModel.FileStoreMetadata) (uint, *errs.XError) {

	dbFileStoreMetadata, err := svc.mapper.FileStoreMetadata(fileStoreMetadata)
	if err != nil {
		return 0, errs.NewXError(errs.INVALID_REQUEST, "Unable to save File Store Metadata", err)
	}
	return svc.fileStoreRepo.Create(ctx, dbFileStoreMetadata)

}

func (svc fileStoreService) UpdateFileStoreMetadata(ctx *context.Context, fileStoreMetadata requestModel.FileStoreMetadata, id uint) *errs.XError {

	dbFileStoreMetadata, err := svc.mapper.FileStoreMetadata(fileStoreMetadata)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to map File Store Metadata", err)
	}
	dbFileStoreMetadata.ID = id
	return svc.fileStoreRepo.Update(ctx, dbFileStoreMetadata)
}

func (svc fileStoreService) GetFileStoreMetadata(ctx *context.Context, id uint) (*responseModel.FileStoreMetadata, *errs.XError) {

	fileStoreMetadata, err := svc.fileStoreRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	respFileStoreMetadata, errr := svc.respMapper.FileStoreMetadata(fileStoreMetadata)
	if errr != nil {
		return nil, errs.Wrap(errr, "Unable to map File Store Metadata")
	}
	return respFileStoreMetadata, nil

}

func (svc fileStoreService) DeleteFileStoreMetadata(ctx *context.Context, id uint) *errs.XError {
	return svc.fileStoreRepo.Delete(ctx, id)
}

func (svc fileStoreService) GetFileStoreMetadataByKey(ctx *context.Context, entityType string, id uint, kind string) (*responseModel.FileStoreMetadata, *errs.XError) {

	key := svc.GetFileKey(ctx, entityType, id, kind)
	fileStoreMetadata, err := svc.fileStoreRepo.GetByKey(ctx, key)
	if err != nil {
		if err.Type == errs.NOT_EXIST {
			return nil, nil
		} else {
			return nil, err
		}
	}

	respFileStoreMetadata, errr := svc.respMapper.FileStoreMetadata(fileStoreMetadata)
	if errr != nil {
		return nil, errs.Wrap(errr, "Unable to map FileStore Metadata")
	}
	return respFileStoreMetadata, nil
}

func (svc fileStoreService) Upload(ctx *context.Context, file models.FileUpload) *errs.XError {

	fileKey := svc.GetFileKey(ctx, file.EntityType, file.EntityId, file.Kind)

	var contentType = defaultContentType

	// Get content type, default is application/octet-stream if not provided
	if file.Metadata.Header != nil {
		if contentTypes, ok := file.Metadata.Header["Content-Type"]; ok && len(contentTypes) > 0 {
			contentType = contentTypes[0]
		}
	}

	// Upload file to S3
	if err := svc.s3Storage.Upload(*ctx, fileKey, file.Content, contentType); err != nil {
		return errs.NewXError(errs.STORAGE, fmt.Sprintf("Failed to upload file  to S3: %v", err), err)
	}

	// Get the S3 URL
	offerUrl, err := svc.s3Storage.GetURL(*ctx, fileKey, 24*3*time.Hour)
	if err != nil {
		return errs.NewXError(errs.STORAGE, "Failed to get offer file presigned URL", err)
	}

	fileStoreMetadata := requestModel.FileStoreMetadata{
		IsActive:   true,
		FileName:   file.Metadata.Filename,
		FileSize:   file.Metadata.Size,
		FileType:   contentType,
		EntityId:   file.EntityId,
		EntityType: file.EntityType,
		FileUrl:    offerUrl,
		FileKey:    fileKey,
		FileBucket: svc.s3Storage.GetBucket(),
	}

	ok, metadata, xerr := svc.GetFileMetadataIfExists(ctx, file.EntityType, file.EntityId, file.Kind)
	if xerr != nil {
		return xerr
	}

	// if metadata has already existed, then update if not fresh save
	if !ok {
		_, xerr = svc.SaveFileStoreMetadata(ctx, fileStoreMetadata)
		if xerr != nil {
			return xerr
		}
	} else {
		if xerr := svc.UpdateFileStoreMetadata(ctx, fileStoreMetadata, metadata.Id); xerr != nil {
			return xerr
		}
	}

	return nil

}

func (svc fileStoreService) GetFileKey(ctx *context.Context, entityName string, entityId uint, fileType string) string {
	return fmt.Sprintf("%s/%d/%s", entityName, entityId, fileType)
}

func (svc fileStoreService) GetFileMetadataIfExists(ctx *context.Context, entityType string, id uint, kind string) (bool, *responseModel.FileStoreMetadata, *errs.XError) {
	fileKey := svc.GetFileKey(ctx, entityType, id, kind)
	fileStoreMetadata, err := svc.fileStoreRepo.GetByKey(ctx, fileKey)
	if err != nil {
		if err.Type == errs.NOT_EXIST {
			return false, nil, nil
		} else {
			return false, nil, err
		}
	}

	respFileStoreMetadata, errr := svc.respMapper.FileStoreMetadata(fileStoreMetadata)
	if errr != nil {
		return false, nil, errs.NewXError(errs.INTERNAL, "Unable to map FileStoreMetadata", errr)
	}
	return true, respFileStoreMetadata, nil
}
