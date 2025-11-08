package schema

import (
	"github.com/google/uuid"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Bundle_item struct{ ent.Schema }

func (Bundle_item) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Int("qty").Default(1),
	}
}

func (Bundle_item) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("bundle", Bundle.Type).Ref("items").Unique().Required(),
		edge.From("product", Product.Type).Ref("bundle_items").Unique().Required(),
	}
}

