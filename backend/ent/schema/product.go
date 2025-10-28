package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Product struct {
	ent.Schema
}

func (Product) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name").NotEmpty(),
		field.Float("price").Positive(),
		field.String("description").NotEmpty(),
		field.String("image_url").NotEmpty(),
		field.String("unit_label").NotEmpty(),
		field.String("is_active").NotEmpty(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Nillable().Optional(),
	}
}

func (Product) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("catagory", Product_category.Type).Ref("product"),
		edge.From("vendor", Vendor.Type).Ref("product"),
		edge.To("inventory", Inventory.Type).Unique().Required(),
	}
}
