package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kovalenkoyo81/weather/internal/models/"
)

func (r *Rest) createUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := r.service.CreateNewUser(c, user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *Rest) userExists(c *gin.Context) {
	name := c.Param("name")
	exists, err := r.service.UserExists(c, name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": exists})
}
