package handler

// DEBUG

import (
	"net/http"
	"strconv"

	"github.com/Draskown/WBL0/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) testGetDB(c *gin.Context) {
	type id struct {
		Id int `json:"id"`
	}

	var input id

	resultId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		if err = c.BindJSON(&input); err != nil {
			throwError(c, http.StatusBadRequest, "No id provided")
			return
		}
		resultId = input.Id
	}

	test, err := h.service.DBConv.TestGetDB(resultId)
	if err != nil {
		throwError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, test)
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

func (h *Handler) showTestDB(c *gin.Context) {
	result, err := h.service.DBConv.ShowTestDB()
	if err != nil {
		throwError(c, http.StatusNoContent, err.Error())
		return
	}

	c.HTML(http.StatusOK,
		"index.html",
		gin.H{
			"result": result,
		},
	)
}

func (h *Handler) showTestDBbyId(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	result, err := h.service.DBConv.ShowTestDBbyId(id)
	if err != nil {
		throwError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK,
		"get_id.html",
		gin.H{
			"result": result,
		},
	)
}
