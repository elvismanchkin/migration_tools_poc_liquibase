package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TemplateCategory holds the schema definition for the TemplateCategory entity.
type TemplateCategory struct {
	ent.Schema
}

// Fields of the TemplateCategory.
func (TemplateCategory) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			StructTag(`json:"id,omitempty"`),
		field.String("name").
			Unique().
			NotEmpty().
			MaxLen(100),
		field.Text("description").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the TemplateCategory.
func (TemplateCategory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("templates", Template.Type),
	}
}
