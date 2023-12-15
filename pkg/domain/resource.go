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
	Description string        `json:"description"`
	Reference   string        `json:"reference"`
	Level       ResourceLevel `json:"level"`
	Type        ResourceType  `json:"type"`
	Tags        []Tag         `json:"tags"`
}
