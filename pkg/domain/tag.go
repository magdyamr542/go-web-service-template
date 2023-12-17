package domain

type Tag struct {
	DefaultFields
	Name string `json:"name" db:"name"`
}
