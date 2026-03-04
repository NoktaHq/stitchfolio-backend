package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	"github.com/imkarthi24/sf-backend/internal/service"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/response"
	"github.com/loop-kar/pixie/util"
)

type ProductHandler struct {
	productSvc service.ProductService
	resp       response.Response
	dataResp   response.DataResponse
}

func ProvideProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{productSvc: svc}
}

//	@Summary		Save Product
//	@Description	Saves an instance of Product
//	@Tags			Product
//	@Accept			json
//	@Success		201		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		500		{object}	responseModel.Response
//	@Param			product	body		requestModel.Product	true	"product"
//	@Router			/product [post]
func (h ProductHandler) SaveProduct(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var product requestModel.Product
	err := ctx.Bind(&product)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	errr := h.productSvc.SaveProduct(&context, product)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Save success").FormatAndSend(&context, ctx, http.StatusCreated)
}

//	@Summary		Update Product
//	@Description	Updates an instance of Product
//	@Tags			Product
//	@Accept			json
//	@Success		202		{object}	responseModel.Response
//	@Failure		400		{object}	responseModel.Response
//	@Failure		500		{object}	responseModel.Response
//	@Param			product	body		requestModel.Product	true	"product"
//	@Param			id		path		int						true	"Product id"
//	@Router			/product/{id} [put]
func (h ProductHandler) UpdateProduct(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)
	var product requestModel.Product
	err := ctx.Bind(&product)
	if err != nil {
		x := errs.NewXError(errs.INVALID_REQUEST, errs.MALFORMED_REQUEST, err)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	errr := h.productSvc.UpdateProduct(&context, product, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusInternalServerError)
		return
	}

	h.resp.SuccessResponse("Update success").FormatAndSend(&context, ctx, http.StatusAccepted)
}

//	@Summary		Get a specific Product
//	@Description	Get an instance of Product with inventory
//	@Tags			Product
//	@Accept			json
//	@Success		200	{object}	responseModel.Product
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			id	path		int	true	"Product id"
//	@Router			/product/{id} [get]
func (h ProductHandler) Get(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))

	product, errr := h.productSvc.Get(&context, uint(id))
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(product).FormatAndSend(&context, ctx, http.StatusOK)
}

//	@Summary		Get all active products
//	@Description	Get all active products with current stock
//	@Tags			Product
//	@Accept			json
//	@Success		200		{object}	responseModel.Product
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/product [get]
func (h ProductHandler) GetAllProducts(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	products, errr := h.productSvc.GetAll(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(products).FormatAndSend(&context, ctx, http.StatusOK)
}

//	@Summary		Delete Product
//	@Description	Deletes an instance of Product
//	@Tags			Product
//	@Accept			json
//	@Success		200	{object}	responseModel.Response
//	@Failure		400	{object}	responseModel.Response
//	@Param			id	path		int	true	"product id"
//	@Router			/product/{id} [delete]
func (h ProductHandler) Delete(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	id, _ := strconv.Atoi(ctx.Param("id"))
	err := h.productSvc.Delete(&context, uint(id))
	if err != nil {
		h.resp.DefaultFailureResponse(err).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.resp.SuccessResponse("Delete Success").FormatAndSend(&context, ctx, http.StatusOK)
}

//	@Summary		Autocomplete for products
//	@Description	Autocomplete for products with stock info
//	@Tags			Product
//	@Accept			json
//	@Success		200		{object}	responseModel.ProductAutoComplete
//	@Failure		400		{object}	responseModel.DataResponse
//	@Param			search	query		string	false	"search"
//	@Router			/product/autocomplete [get]
func (h ProductHandler) AutocompleteProduct(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	search := ctx.Query("search")
	search = util.EncloseWithSingleQuote(search)

	products, errr := h.productSvc.AutocompleteProduct(&context, search)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(products).FormatAndSend(&context, ctx, http.StatusOK)
}

//	@Summary		Get product by SKU
//	@Description	Get product details by SKU
//	@Tags			Product
//	@Accept			json
//	@Success		200	{object}	responseModel.Product
//	@Failure		400	{object}	responseModel.DataResponse
//	@Param			sku	query		string	true	"Product SKU"
//	@Router			/product/sku [get]
func (h ProductHandler) GetBySKU(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	sku := ctx.Query("sku")
	if sku == "" {
		x := errs.NewXError(errs.INVALID_REQUEST, "SKU is required", nil)
		h.resp.DefaultFailureResponse(x).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	product, errr := h.productSvc.GetBySKU(&context, sku)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(product).FormatAndSend(&context, ctx, http.StatusOK)
}

//	@Summary		Get low stock products
//	@Description	Get all products with stock below threshold
//	@Tags			Product
//	@Accept			json
//	@Success		200	{object}	responseModel.Product
//	@Failure		400	{object}	responseModel.DataResponse
//	@Router			/product/low-stock [get]
func (h ProductHandler) GetLowStockProducts(ctx *gin.Context) {
	context := util.CopyContextFromGin(ctx)

	products, errr := h.productSvc.GetLowStockProducts(&context)
	if errr != nil {
		h.resp.DefaultFailureResponse(errr).FormatAndSend(&context, ctx, http.StatusBadRequest)
		return
	}

	h.dataResp.DefaultSuccessResponse(products).FormatAndSend(&context, ctx, http.StatusOK)
}
