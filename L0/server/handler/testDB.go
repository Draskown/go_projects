package handler

// DEBUG

import (
	"net/http"

	"github.com/Draskown/WBL0/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) testGetDB(c *gin.Context) {
	type id struct {
		Id int `json:"id"`
	}

	var input id

	if err := c.BindJSON(&input); err != nil {
		throwError(c, http.StatusBadRequest, err.Error())
		return
	}

	test, err := h.service.DBConv.TestGetDB(input.Id)
	if err != nil {
		throwError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"value": test.Value,
		"text": test.Text,
		"arr_one": test.Arr_One,
	})
}

func (h *Handler) testPostDB(c *gin.Context) {
	var input model.Test

	if err := c.BindJSON(&input); err != nil {
		throwError(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.DBConv.TestPostDB(input)
	if err != nil {
		throwError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}