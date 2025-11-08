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
		field.String("name"),
		field.String("sku").Unique(),
		field.String("description").Nillable().Optional(),
		field.Float("price"),
		field.String("unit_label"),
		field.Bool("is_active").Default(true),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Product) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("vendor", Vendor.Type).Ref("products").Unique(),
		edge.To("product_categories", Product_category.Type),
		edge.To("inventories", Inventory.Type),
		edge.To("cart_items", Cart_item.Type),
		edge.To("order_items", Order_item.Type),
		edge.To("bundle_items", Bundle_item.Type),
		edge.To("recipe_items", Recipe_item.Type),
		edge.To("reviews", Review.Type),
	}
}
