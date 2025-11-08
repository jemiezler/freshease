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
		field.Int("quantity").Default(0),
		field.Int("reorder_level").Default(0),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Inventory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("product", Product.Type).Ref("inventories").Unique().Required(),
		edge.From("vendor", Vendor.Type).Ref("inventories").Unique().Required(),
	}
}
