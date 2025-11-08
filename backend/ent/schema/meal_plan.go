package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Meal_plan struct{ ent.Schema }

func (Meal_plan) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Time("week_start"),
		field.String("goal").Nillable().Optional(),
	}
}

func (Meal_plan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("meal_plans").Unique().Required(),
		edge.To("items", Meal_plan_item.Type),
	}
}
