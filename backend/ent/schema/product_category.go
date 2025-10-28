package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Product_category struct{ ent.Schema }

func (Product_category) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable(),
		field.String("name").NotEmpty(),
		field.String("description").NotEmpty(),
		field.String("slug").NotEmpty(),
	}
}

func (Product_category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("product", Product.Type),
	}
}
