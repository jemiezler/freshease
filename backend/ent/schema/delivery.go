package schema

import (
	"github.com/google/uuid"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Delivery struct{ ent.Schema }

func (Delivery) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("provider"),
		field.String("tracking_no").Nillable().Optional(),
		field.String("status"),
		field.Time("eta").Nillable().Optional(),
		field.Time("delivered_at").Nillable().Optional(),
	}
}

func (Delivery) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).Ref("deliveries").Required(),
	}
}

