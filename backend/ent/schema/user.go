package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name"),
		field.String("email").Unique(),
		field.String("password").Nillable().Optional().Sensitive(),
		field.String("phone").Unique().Nillable().Optional(),
		field.String("address").Nillable().Optional(),
		field.String("bio").Nillable().Optional(),
		field.String("avatar").Nillable().Optional(),
		field.String("cover").Nillable().Optional(),
		field.Time("date_of_birth").Nillable().Optional(),
		field.String("sex").Nillable().Optional(),
		field.String("goal").Nillable().Optional(),
		field.Float("height_cm").Nillable().Optional(),
		field.Float("weight_kg").Nillable().Optional(),
		field.String("status").Nillable().Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Nillable().Optional(),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").Unique(),
		index.Fields("phone").Unique(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", Role.Type).Ref("users").Unique(),
		edge.To("addresses", Address.Type),
		edge.To("carts", Cart.Type),
		edge.To("orders", Order.Type),
		edge.To("notifications", Notification.Type),
		edge.To("reviews", Review.Type),
		edge.To("meal_plans", Meal_plan.Type),
		edge.To("identities", Identity.Type),
	}
}
