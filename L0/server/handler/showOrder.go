package handler

import (
	"net/http"

	"github.com/Draskown/WBL0/model"
	"github.com/gin-gonic/gin"
)

// Route that must be handled after accessing `/:id`
func (h *Handler) showOrder(c *gin.Context) {
	var input model.Order

	if err := c.BindJSON(&input); err != nil {
		throwError(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.DBConv.ShowOrder(input)
	if err != nil {
		throwError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
