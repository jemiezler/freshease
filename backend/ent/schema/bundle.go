package schema

import (
	"github.com/google/uuid"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Bundle struct{ ent.Schema }

func (Bundle) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name"),
		field.String("description").Nillable().Optional(),
		field.Float("price"),
		field.Bool("is_active").Default(true),
	}
}

func (Bundle) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("items", Bundle_item.Type),
	}
}

