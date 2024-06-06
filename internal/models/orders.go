package models

import "time"

type Order struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	UserID    string    `json:"user_id"`
	PizzaID   string    `json:"menu_item_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

type Filter struct {
	UserID  string `json:"user_id,omitempty"`
	PizzaID string `json:"menu_item_id,omitempty"`
}
