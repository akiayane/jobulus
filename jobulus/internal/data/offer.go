package data

import (
	"fmt"
	"time"
)

type Offer struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Salary      string    `json:"salary"`
	Contacts    string    `json:"contacts"`
	CreatedTime time.Time `json:"createdTime"`
}

func (o *Offer) Print() string {
	return fmt.Sprintf("Title: %s \nDescription: %s \nSalary: %s \nContacts: %s \nCreated Time: %s \n", o.Title, o.Description, o.Salary, o.Contacts, o.CreatedTime)
}
