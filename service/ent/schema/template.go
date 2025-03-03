package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Template holds the schema definition for the Template entity.
type Template struct {
	ent.Schema
}

// Fields of the Template.
func (Template) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("name").
			NotEmpty().
			MaxLen(200),
		field.Text("content").
			NotEmpty(),
		field.String("format").
			Default("html").
			MaxLen(50),
		field.Int("version").
			Default(1),
		field.Bool("is_active").
			Default(true),
		field.String("created_by").
			NotEmpty().
			MaxLen(100),
		field.Time("created_at").
			Default(time.Now),
		field.String("updated_by").
			Optional().
			Nillable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Template.
func (Template) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("variables", TemplateVariable.Type),
		edge.From("category", TemplateCategory.Type).
			Ref("templates").
			Unique().
			Required(),
	}
}
