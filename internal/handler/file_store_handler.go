package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/service"
	"github.com/imkarthi24/sf-backend/internal/utils"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/response"
	"github.com/loop-kar/pixie/util"
)

var _ = (*responseModel.TempFileUpload)(nil) // for swagger

type FileStoreHandler struct {
	fileStoreSvc service.FileStoreService
	resp         response.Response
	dataResp     response.DataResponse
}

func ProvideFileStoreHandler(svc service.FileStoreService) *FileStoreHandler {
	return &FileStoreHandler{fileStoreSvc: svc}
}

// UploadTemp uploads a file to temporary storage and returns metadata for later commit.
//
//	@Summary		Upload file to temp storage
//	@Description	Uploads a file to temporary storage. Returns id and tempFileKey to use when committing the file to an entity (e.g. in tempFileRefs).
//	@Tags			FileStore
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file						true	"File to upload"
//	@Success		201		{object}	responseModel.DataResponse	"Returns { data: { id, tempFileKey, fileName } }"
//	@Failure		400		{object}	responseModel.DataResponse
//	@Failure		500		{object}	responseModel.DataResponse
//	@Router			/file-store/temp [post]
//	@Security		BearerAuth
func (h *FileStoreHandler) UploadTemp(ctx *gin.Context) {
	appCtx := util.CopyContextFromGin(ctx)
	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, "Error parsing multipart form", err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&appCtx, ctx, http.StatusBadRequest)
		return
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, "Missing or invalid file", err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&appCtx, ctx, http.StatusBadRequest)
		return
	}

	file, err := utils.ExtractFile(fileHeader)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, "Error extracting file", err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&appCtx, ctx, http.StatusBadRequest)
		return
	}

	result, xerr := h.fileStoreSvc.UploadTemp(&appCtx, *file)
	if xerr != nil {
		code := http.StatusInternalServerError
		if xerr.Type == errs.VALIDATION || xerr.Type == errs.INVALID_REQUEST {
			code = http.StatusBadRequest
		}
		h.resp.DefaultFailureResponse(xerr).FormatAndSend(&appCtx, ctx, code)
		return
	}

	h.dataResp.DefaultSuccessResponse(result).FormatAndSend(&appCtx, ctx, http.StatusCreated)
}

// UploadTempBulk uploads multiple files to temporary storage and returns metadata for each.
//
//	@Summary		Bulk upload files to temp storage
//	@Description	Uploads multiple files to temporary storage. Response array order matches request order (result[i] = i-th file). Each item includes fileName for matching when names are unique.
//	@Tags			FileStore
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			files	formData	file						true	"Files to upload (multiple with same field name 'files')"
//	@Success		201		{object}	responseModel.DataResponse	"Returns { data: [ { id, tempFileKey, fileName }, ... ] }"
//	@Failure		400		{object}	responseModel.DataResponse
//	@Failure		500		{object}	responseModel.DataResponse
//	@Router			/file-store/temp/bulk [post]
//	@Security		BearerAuth
func (h *FileStoreHandler) UploadTempBulk(ctx *gin.Context) {
	appCtx := util.CopyContextFromGin(ctx)
	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, "Error parsing multipart form", err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&appCtx, ctx, http.StatusBadRequest)
		return
	}

	fileHeaders := ctx.Request.MultipartForm.File["files"]
	if len(fileHeaders) == 0 {
		x := errs.NewXError(errs.INVALID_REQUEST, "No files provided; use form field 'files' with one or more files", nil)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&appCtx, ctx, http.StatusBadRequest)
		return
	}

	results := make([]responseModel.TempFileUpload, 0, len(fileHeaders))
	for _, fileHeader := range fileHeaders {
		file, err := utils.ExtractFile(fileHeader)
		if err != nil {
			x := errs.NewXError(errs.INVALID_REQUEST, "Error extracting file: "+fileHeader.Filename, err)
			h.resp.DefaultFailureResponse(x).FormatAndSend(&appCtx, ctx, http.StatusBadRequest)
			return
		}
		result, xerr := h.fileStoreSvc.UploadTemp(&appCtx, *file)
		if xerr != nil {
			code := http.StatusInternalServerError
			if xerr.Type == errs.VALIDATION || xerr.Type == errs.INVALID_REQUEST {
				code = http.StatusBadRequest
			}
			h.resp.DefaultFailureResponse(xerr).FormatAndSend(&appCtx, ctx, code)
			return
		}
		results = append(results, *result)
	}

	h.dataResp.DefaultSuccessResponse(results).FormatAndSend(&appCtx, ctx, http.StatusCreated)
}
