package schema

import (
	"time"

	"github.com/google/uuid"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Order struct{ ent.Schema }

func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable(),
		field.String("order_no").Unique(),
		field.String("status"),
		field.Float("subtotal").Default(0.0),
		field.Float("shipping_fee").Default(0.0),
		field.Float("discount").Default(0.0),
		field.Float("total").Default(0.0),
		field.Time("placed_at").Nillable().Optional(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Order) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("order_no").Unique(),
	}
}

func (Order) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("orders").Required(),
		edge.From("shipping_address", Address.Type).Ref("shipping_orders"),
		edge.From("billing_address", Address.Type).Ref("billing_orders"),
		edge.To("items", Order_item.Type),
		edge.To("payments", Payment.Type),
		edge.To("deliveries", Delivery.Type),
	}
}
