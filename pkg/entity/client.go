package entity

import (
	"strconv"
	"strings"
	"time"

	"neoway/pkg/lib/cnpj"
	"neoway/pkg/lib/cpf"
	"neoway/pkg/lib/regex"
)

// Client represents a client entity.
type Client struct {
	ID                int        `db:"id"`
	CPF               string     `db:"cpf"`
	Private           bool       `db:"private"`
	Incomplete        bool       `db:"incomplete"`
	LastOrderDate     *time.Time `db:"last_order_date"`
	AverageTicket     *float64   `db:"average_ticket"`
	LastOrderTicket   *float64   `db:"last_order_ticket"`
	MostFrequentStore *string    `db:"most_frequent_store"`
	LastOrderStore    *string    `db:"last_order_store"`

	ValidCPF               bool `db:"valid_cpf"`
	ValidMostFrequentStore bool `db:"valid_most_frequent_store"`
	ValidLastOrderStore    bool `db:"valid_last_order_store"`
}

// NewClient creates a new client entity.
func NewClient(data []string) (*Client, error) {
	var lastOrderDate *time.Time
	if date := data[3]; date != "NULL" {
		t, err := time.Parse(time.DateOnly, date)
		if err != nil {
			return nil, err
		}

		lastOrderDate = &t
	}

	var averageTicket *float64
	if ticket := data[4]; ticket != "NULL" {
		ticket = strings.Replace(ticket, ",", ".", 1)
		floatTicket, err := strconv.ParseFloat(ticket, 64)
		if err != nil {
			return nil, err
		}

		averageTicket = &floatTicket
	}

	var lastOrderTicket *float64
	if ticket := data[5]; ticket != "NULL" {
		ticket = strings.Replace(ticket, ",", ".", 1)
		floatTicket, err := strconv.ParseFloat(ticket, 64)
		if err != nil {
			return nil, err
		}

		lastOrderTicket = &floatTicket
	}

	var mostFrequentStore *string
	if store := data[6]; store != "NULL" {
		store = regex.OnlyDigits.ReplaceAllString(store, "")
		mostFrequentStore = &store
	}

	var lastOrderStore *string
	if store := data[7]; store != "NULL" {
		store = regex.OnlyDigits.ReplaceAllString(store, "")
		lastOrderStore = &store
	}

	return &Client{
		CPF:               regex.OnlyDigits.ReplaceAllString(data[0], ""),
		Private:           data[1] == "1",
		Incomplete:        data[2] == "1",
		LastOrderDate:     lastOrderDate,
		AverageTicket:     averageTicket,
		LastOrderTicket:   lastOrderTicket,
		MostFrequentStore: mostFrequentStore,
		LastOrderStore:    lastOrderStore,
	}, nil
}

// ValidateDocuments validates the client documents and stores the result in the entity.
func (c *Client) ValidateDocuments() {
	c.ValidCPF = cpf.Validate(c.CPF)
	c.ValidMostFrequentStore = c.MostFrequentStore != nil && cnpj.Validate(*c.MostFrequentStore)
	c.ValidLastOrderStore = c.LastOrderStore != nil && cnpj.Validate(*c.LastOrderStore)
}
