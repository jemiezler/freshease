package schema

import (
	"github.com/google/uuid"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Recipe_item struct{ ent.Schema }

func (Recipe_item) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Float("amount").Default(0.0),
		field.String("unit"),
	}
}

func (Recipe_item) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("recipe", Recipe.Type).Ref("items").Unique().Required(),
		edge.From("product", Product.Type).Ref("recipe_items").Unique().Required(),
	}
}

