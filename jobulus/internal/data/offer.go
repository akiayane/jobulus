package data

import (
	"database/sql"
	"fmt"
	"time"
)

type Offer struct {
	Id             int64     `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Salary         string    `json:"salary"`
	Contacts       string    `json:"contacts"`
	Schedule       string    `json:"schedule"`
	EmploymentType string    `json:"employmentType"`
	CreatedTime    time.Time `json:"createdTime"`
}

func (o *Offer) Print() string {
	return fmt.Sprintf("Title: %s \nDescription: %s \nSchedule: %s \nEmployment Type: %s \nSalary: %s \nContacts: %s \nCreated Time: %s \n",
		o.Title, o.Description, o.Schedule, o.EmploymentType, o.Salary, o.Contacts, o.CreatedTime)
}

type OfferModel struct {
	DB *sql.DB
}

func (m *OfferModel) GetAll() ([]*Offer, error) {
	stmt := `SELECT id, title, description, salary, contacts, schedule, employmentType, createdTime FROM offers;`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []*Offer
	for rows.Next() {
		u := &Offer{}
		err = rows.Scan(&u.Id, &u.Title, &u.Description, &u.Salary, &u.Contacts, &u.Schedule, &u.EmploymentType, &u.CreatedTime)
		if err != nil {
			return nil, err
		}
		offers = append(offers, u)
	}

	if len(offers) == 0 {
		return nil, sql.ErrNoRows
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return offers, nil
}

func (m *OfferModel) Insert(title, desciption, salary, contacts, schedule, employmentType string) (int, error) {
	stmt := `INSERT INTO offers (title, description, salary, contacts, schedule, employmentType) VALUES (?,?,?,?,?,?);`

	result, err := m.DB.Exec(stmt, title, desciption, salary, contacts, schedule, employmentType)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
