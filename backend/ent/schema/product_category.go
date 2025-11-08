package schema

import (
	"github.com/google/uuid"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Product_category struct{ ent.Schema }

func (Product_category) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
	}
}

func (Product_category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("product", Product.Type).Ref("product_categories").Unique().Required(),
		edge.From("category", Category.Type).Ref("product_categories").Unique().Required(),
	}
}
