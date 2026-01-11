package gin_comp

import (
	"net/http"

	"github.com/dukk308/golang-clean-arc/internal/common"
	"github.com/dukk308/golang-clean-arc/pkgs/ddd"
	"github.com/gin-gonic/gin"
)

func ResponseError(c *gin.Context, err error) {
	if ddd.IsDomainError(err) {
		c.JSON(http.StatusBadRequest, ddd.ToDomainError(err))
	} else {
		c.JSON(http.StatusInternalServerError, ddd.ToDomainError(err))
	}

}

func ResponseSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, common.NewResponseSuccess(data))
}

func ResponseSuccessCreated(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, common.NewResponseSuccess(data))
}
