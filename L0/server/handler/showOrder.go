package handler

import (
	"net/http"

	"github.com/Draskown/WBL0/model"
	"github.com/gin-gonic/gin"
)

// Route that must be handled after accessing `/:id`
//
// Parses id from the route (i.e. localhost:8080/123,
// 123 would be the id)
func (h *Handler) showOrder(c *gin.Context) {
	id := c.Param("id")

	result, err := h.service.DBConv.ShowOrder(id)
	if err != nil {
		throwError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK,
		"index.html",
		gin.H{
			"id":     result.OrderId,
			"orders": []model.Order{result},
		},
	)
}
