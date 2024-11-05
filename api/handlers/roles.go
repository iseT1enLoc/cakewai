package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/component/response"
	"cakewai/cakewai.com/domain"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	RoleUsecase domain.RoleUsecase
	Env         *appconfig.Env
}

func (rc *RoleController) GetAllRoles() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roles, err := rc.RoleUsecase.GetAllRoles(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Fail to get all roles",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully get all roles",
			},
			Data: roles,
		})
	}
}
func (rc *RoleController) GetRoleByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleID := ctx.Param("id")

		hexId, err := primitive.ObjectIDFromHex(roleID)
		role, err := rc.RoleUsecase.GetRoleByID(ctx, hexId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.BasicResponse{
				Code:    0,
				Message: "Error happened while get role by id in database",
				Error:   err.Error(),
				Data:    nil,
			})
			return
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Get role by id successfully",
			},
			Data: role,
		})
	}
}
func (rc *RoleController) CreateRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var role domain.Role
		err := ctx.ShouldBindJSON(&role)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.BasicResponse{
				Code:    0,
				Message: "Can not parse json",
				Error:   "",
				Data:    nil,
			})
			return
		}
		role.RoleID = primitive.NewObjectID()
		err = rc.RoleUsecase.CreateRole(ctx, role)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.BasicResponse{
				Code:    0,
				Message: "Error happened while inserting into database",
				Error:   "",
				Data:    nil,
			})
			return
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully create new role",
			},
			Data: nil,
		})
	}
}
func (rc *RoleController) UpdateRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var role domain.Role
		roleID := ctx.Param("id")
		hex, err := primitive.ObjectIDFromHex(roleID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "can not get the id",
				Error:   err.Error(),
			})
			return
		}
		if err = ctx.ShouldBindJSON(&role); err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    0,
				Message: "Can not parsing json file",
				Error:   err.Error(),
			})
		}
		role.RoleID = hex
		updatedrole, err := rc.RoleUsecase.UpdateRole(ctx, role)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.BasicResponse{
				Code:    0,
				Message: "Error happened while inserting into database",
				Error:   "",
				Data:    nil,
			})
			return
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    200,
				Message: "Successfully update role",
			},
			Data: updatedrole,
		})
	}
}
func (rc *RoleController) GetRoleByRoleName() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		rolename := ctx.Param("name")

		role, err := rc.RoleUsecase.GetRoleByRoleName(ctx, rolename)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error happened while query database",
				Error:   err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully get role by name",
			},
			Data: role,
		})
	}
}
