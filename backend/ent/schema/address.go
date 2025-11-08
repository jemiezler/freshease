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
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("line1"),
		field.String("line2").Nillable().Optional(),
		field.String("city"),
		field.String("province"),
		field.String("postal_code"),
		field.String("country"),
		field.Float("lat").Nillable().Optional(),
		field.Float("lng").Nillable().Optional(),
		field.Bool("is_default").Default(false),
	}
}

// Edges of the Address.
func (Address) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("addresses").Unique().Required(),
		edge.To("shipping_orders", Order.Type),
		edge.To("billing_orders", Order.Type),
	}
}
