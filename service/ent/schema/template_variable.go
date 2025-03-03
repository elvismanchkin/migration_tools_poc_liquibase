package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TemplateVariable holds the schema definition for the TemplateVariable entity.
type TemplateVariable struct {
	ent.Schema
}

// Fields of the TemplateVariable.
func (TemplateVariable) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			StructTag(`json:"id,omitempty"`),
		field.String("variable_name").
			NotEmpty().
			MaxLen(100),
		field.Text("description").
			Optional().
			Nillable(),
		field.Text("default_value").
			Optional().
			Nillable(),
		field.Bool("is_required").
			Default(false),
		field.String("variable_type").
			Default("string").
			MaxLen(50),
	}
}

// Edges of the TemplateVariable.
func (TemplateVariable) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("template", Template.Type).
			Ref("variables").
			Unique().
			Required(),
	}
}
