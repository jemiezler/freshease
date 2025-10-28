package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Cart_item struct{ ent.Schema }

func (Cart_item) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable(),
		field.String("name").NotEmpty(),
		field.String("description").NotEmpty(),
	}
}

func (Cart_item) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("cart", Cart.Type).Unique().Required(),
	}
}
