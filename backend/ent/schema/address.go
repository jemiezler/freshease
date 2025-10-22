package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Address holds the schema definition for the Address entity.
type Address struct {
	ent.Schema
}

// Fields of the Address.
func (Address) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.String("line1").NotEmpty(),
		field.String("line2").Optional(),
		field.String("city").NotEmpty(),
		field.String("province").NotEmpty(),
		field.String("country").NotEmpty(),
		field.String("zip").NotEmpty(),
		field.Bool("is_default").Default(false),
	}
}

// Edges of the Address.
func (Address) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("address").Required(),
	}
}
