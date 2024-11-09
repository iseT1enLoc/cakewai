package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/component/response"
	"cakewai/cakewai.com/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductTypeHandler struct {
	Product_type_usecase domain.ProductTypeUsecase
	Env                  *appconfig.Env
}

func (p *ProductTypeHandler) CreateProductType() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var product_type domain.ProductType
		if err := ctx.ShouldBindBodyWithJSON(&product_type); err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Can not parsing the file",
				Error:   err.Error(),
			})
			return
		}
		err := p.Product_type_usecase.CreateProductType(ctx, product_type)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while querying database",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully adding productype",
			},
			Data: nil,
		})

	}
}
