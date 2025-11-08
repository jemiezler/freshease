package schema

import (
	"github.com/google/uuid"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Meal_plan_item struct{ ent.Schema }

func (Meal_plan_item) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Time("day"),
		field.String("slot"),
	}
}

func (Meal_plan_item) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("meal_plan", Meal_plan.Type).Ref("items").Unique().Required(),
		edge.From("recipe", Recipe.Type).Ref("meal_plan_items").Unique().Required(),
	}
}

