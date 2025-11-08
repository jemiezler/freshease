package schema

import (
	"github.com/google/uuid"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Category struct{ ent.Schema }

func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name").Unique(),
		field.String("slug").Unique(),
	}
}

func (Category) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
		index.Fields("slug").Unique(),
	}
}

func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("product_categories", Product_category.Type),
	}
}

