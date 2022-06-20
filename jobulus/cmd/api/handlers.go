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

	go func() {
		app.offers <- input
	}()

	_, err := app.models.Offer.Insert(input.Title, input.Description, input.Salary, input.Contacts, input.Schedule, input.EmploymentType)
	if err != nil {
		app.serverErrorResponse(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"payload": "sent successfully"})

}

func (app *application) GetAllOffers(c *gin.Context) {

	offers, err := app.models.Offer.GetAll()
	if err != nil {
		app.serverErrorResponse(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"payload": offers})

}
