package schema

import (
	"github.com/google/uuid"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Recipe struct{ ent.Schema }

func (Recipe) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name"),
		field.String("instructions").Nillable().Optional(),
		field.Int("kcal").Default(0),
	}
}

func (Recipe) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("items", Recipe_item.Type),
		edge.To("meal_plan_items", Meal_plan_item.Type),
	}
}

