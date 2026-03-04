package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/config"
	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/storage"
)

type EntityDocumentService interface {
	SaveEntityDocument(*context.Context, requestModel.EntityDocuments) *errs.XError
	UpdateEntityDocument(*context.Context, requestModel.EntityDocuments, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.EntityDocument, *errs.XError)
	Delete(*context.Context, uint) *errs.XError

	GetEntityDocumentsByEntity(*context.Context, uint, entities.EntityName, entities.EntityDocumentsType) ([]responseModel.EntityDocument, *errs.XError)
}

type entityDocumentService struct {
	entityDocumentRepo repository.EntityDocumentRepository
	fileStoreSvc       FileStoreService
	mapper             mapper.Mapper
	respMapper         mapper.ResponseMapper
	config             config.AppConfig
	s3Storage          storage.CloudStorageProvider
}

func ProvideEntityDocumentService(entityDocumentRepo repository.EntityDocumentRepository, fileStoreSvc FileStoreService, mapper mapper.Mapper, respMapper mapper.ResponseMapper, config config.AppConfig, s3Storage storage.CloudStorageProvider) EntityDocumentService {
	return entityDocumentService{
		entityDocumentRepo: entityDocumentRepo,
		fileStoreSvc:       fileStoreSvc,
		mapper:             mapper,
		respMapper:         respMapper,
		config:             config,
		s3Storage:          s3Storage,
	}
}

func (svc entityDocumentService) SaveEntityDocument(ctx *context.Context, entityDocument requestModel.EntityDocuments) *errs.XError {
	dbEntityDocument, err := svc.mapper.EntityDocument(entityDocument)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save Entity Documents", err)
	}

	if xerr := svc.entityDocumentRepo.Create(ctx, dbEntityDocument); xerr != nil {
		return xerr
	}

	if entityDocument.Document != nil && entityDocument.Document.HasContent() {
		ok := entityDocument.Document.AddEntityInfo(dbEntityDocument.ID, string(entities.Entity_EntityDocuments), entityDocument.Type)
		if ok {
			xerr := svc.fileStoreSvc.Upload(ctx, *entityDocument.Document)
			if xerr != nil {
				return xerr
			}
		}
	}

	return nil
}

func (svc entityDocumentService) UpdateEntityDocument(ctx *context.Context, entityDocument requestModel.EntityDocuments, id uint) *errs.XError {
	dbEntityDocument, err := svc.mapper.EntityDocument(entityDocument)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update Entity Documents", err)
	}

	dbEntityDocument.ID = id
	if xerr := svc.entityDocumentRepo.Update(ctx, dbEntityDocument); xerr != nil {
		return xerr
	}

	if entityDocument.Document != nil && entityDocument.Document.HasContent() {
		ok := entityDocument.Document.AddEntityInfo(dbEntityDocument.ID, string(entities.Entity_EntityDocuments), entityDocument.Type)
		if ok {
			xerr := svc.fileStoreSvc.Upload(ctx, *entityDocument.Document)
			if xerr != nil {
				return xerr
			}
		}
	}

	return nil
}

// NOT USED
func (svc entityDocumentService) Get(ctx *context.Context, id uint) (*responseModel.EntityDocument, *errs.XError) {
	entityDocument, err := svc.entityDocumentRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	respEntityDocument, errr := svc.respMapper.EntityDocument(entityDocument)
	if errr != nil {
		return nil, err
	}

	// Get document metadata and URL
	documentMetadata, xerr := svc.fileStoreSvc.GetFileStoreMetadataByKey(ctx, string(entityDocument.EntityName), entityDocument.EntityId, string(entityDocument.Type))
	if xerr != nil {
		return nil, xerr
	}

	if documentMetadata != nil {
		url, err := svc.s3Storage.GetCurrentOrRenewedURL(ctx, documentMetadata.FileUrl, documentMetadata.FileKey)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		respEntityDocument.Document = &responseModel.FileResponse{
			FileUrl:  url,
			FileName: documentMetadata.FileName,
		}
	}

	return respEntityDocument, nil
}

func (svc entityDocumentService) Delete(ctx *context.Context, id uint) *errs.XError {
	return svc.entityDocumentRepo.Delete(ctx, id)
}

func (svc entityDocumentService) GetEntityDocumentsByEntity(ctx *context.Context, entityId uint, entityName entities.EntityName, typ entities.EntityDocumentsType) ([]responseModel.EntityDocument, *errs.XError) {

	entityDocuments, err := svc.entityDocumentRepo.GetEntityDocumentsByEntity(ctx, entityId, string(entityName), string(typ))
	if err != nil {
		return nil, err
	}

	respEntityDocuments, errr := svc.respMapper.EntityDocuments(entityDocuments)
	if errr != nil {
		return nil, err
	}

	// Fetch file metadata for each entity document
	for i := range respEntityDocuments {
		ok, documentMetadata, xerr := svc.fileStoreSvc.GetFileMetadataIfExists(ctx, string(entities.Entity_EntityDocuments), entityDocuments[i].ID, string(entityDocuments[i].Type))
		if xerr != nil {
			return nil, xerr
		}
		if ok {
			url, err := svc.s3Storage.GetCurrentOrRenewedURL(ctx, documentMetadata.FileUrl, documentMetadata.FileKey)
			if err != nil {
				return nil, errs.Wrap(err)
			}
			respEntityDocuments[i].Document = &responseModel.FileResponse{
				FileUrl:  url,
				FileName: documentMetadata.FileName,
			}
		}
	}

	return respEntityDocuments, nil
}
