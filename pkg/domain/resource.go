package domain

type ResourceLevel string

const (
	ADVANCED     ResourceLevel = "ADVANCED"
	BEGINNER     ResourceLevel = "BEGINNER"
	INTERMEDIATE ResourceLevel = "INTERMEDIATE"
)

type ResourceType string

const (
	ARTICLE ResourceType = "ARTICLE"
	VIDEO   ResourceType = "VIDEO"
)

type Resource struct {
	DefaultFields
	Description string        `json:"description" db:"description"`
	Reference   string        `json:"reference" db:"reference"`
	Level       ResourceLevel `json:"level" db:"level"`
	Type        ResourceType  `json:"type" db:"type"`
	Tags        []string      `json:"tags" db:"tags"`
}
