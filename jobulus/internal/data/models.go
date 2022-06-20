package data

import "database/sql"

type Healthcheck struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type Models struct {
	Offer OfferModel
}

func NewModes(db *sql.DB) Models {
	return Models{
		OfferModel{DB: db},
	}
}
