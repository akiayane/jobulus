package main

import (
	data "jobulus/internal/data"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (app *application) Healthcheck(c *gin.Context) {

	health := &data.Healthcheck{
		Status:      "available",
		Environment: app.config.Env,
		Version:     version,
	}

	c.JSON(http.StatusOK, gin.H{"payload": health})

}

func (app *application) RegisterOffer(c *gin.Context) {

	var input data.Offer

	if err := c.BindJSON(&input); err != nil {
		app.serverErrorResponse(err, c)
		return
	}

	input.CreatedTime = time.Now()

	app.offers <- input

	c.JSON(http.StatusOK, gin.H{"payload": "sent successfully"})

}
