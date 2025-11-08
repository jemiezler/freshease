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
		field.Int("qty").Default(1),
		field.Float("unit_price").Default(0.0),
		field.Float("line_total").Default(0.0),
	}
}

func (Cart_item) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("cart", Cart.Type).Ref("items").Unique().Required(),
		edge.From("product", Product.Type).Ref("cart_items").Unique().Required(),
	}
}
