package models

type Build struct {
	ID          int    `json:"id"`
	ComponentID int    `json:"component_id"`
	Name        string `json:"name"`
}
