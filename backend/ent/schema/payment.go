package schema

import (
	"github.com/google/uuid"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Payment struct{ ent.Schema }

func (Payment) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable(),
		field.String("provider"),
		field.String("provider_ref").Nillable().Optional(),
		field.String("status"),
		field.Float("amount").Default(0.0),
		field.Time("paid_at").Nillable().Optional(),
	}
}

func (Payment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).Ref("payments").Required(),
	}
}
