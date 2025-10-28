package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Vendor struct {
	ent.Schema
}

func (Vendor) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name").Nillable().Optional(),
		field.String("email").Nillable().Optional(),
		field.String("phone").Nillable().Optional(),
		field.String("address").Nillable().Optional(),
		field.String("city").Nillable().Optional(),
		field.String("state").Nillable().Optional(),
		field.String("country").Nillable().Optional(),
		field.String("postal_code").Nillable().Optional(),
		field.String("website").Nillable().Optional(),
		field.String("logo_url").Nillable().Optional(),
		field.String("description").Nillable().Optional(),
		field.String("is_active").NotEmpty(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Nillable().Optional(),
	}
}

func (Vendor) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("product", Product.Type),
		edge.To("inventory", Inventory.Type),
	}
}
