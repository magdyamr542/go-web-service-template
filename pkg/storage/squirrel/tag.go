package squirrel

type tagTableIdentifiers struct {
	tableName string
	id        string
	createdAt string
	updatedAt string
	name      string
}

var (
	tagT = tagTableIdentifiers{
		tableName: "app.tag",
		id:        "id",
		createdAt: "created_at",
		updatedAt: "updated_at",
		name:      "name",
	}
)
