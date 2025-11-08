package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Cart struct{ ent.Schema }

func (Cart) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable(),
		field.String("status"),
		field.Float("subtotal").Default(0.0),
		field.Float("discount").Default(0.0),
		field.Float("total").Default(0.0),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Cart) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("carts").Required(),
		edge.To("items", Cart_item.Type),
	}
}
