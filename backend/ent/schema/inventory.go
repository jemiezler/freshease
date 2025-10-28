package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Inventory struct{ ent.Schema }

func (Inventory) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable(),
		field.Int("quantity").Positive(),
		field.Int("restock_amount").Positive(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Inventory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("product", Product.Type).Ref("inventory"),
		edge.From("vendor", Vendor.Type).Ref("inventory"),
	}
}
