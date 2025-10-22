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
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.String("email").Unique(),
		field.String("password").MinLen(8).MaxLen(100),
		field.String("name").MinLen(2).MaxLen(100),
		field.String("phone").MinLen(10).MaxLen(20),
		field.String("bio").MinLen(10).MaxLen(500).Nillable().Optional(),
		field.String("avatar").MinLen(10).MaxLen(200).Nillable().Optional(),
		field.String("cover").MinLen(10).MaxLen(200).Nillable().Optional(),
		field.String("status").Default("active"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Nillable().Optional(),
	}
}

func Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").Unique(),
		index.Fields("phone").Unique(),
		index.Fields("status").Unique(),
		index.Fields("role").Unique(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", Role.Type).Ref("users").Unique(),
		edge.To("address", Address.Type),
		edge.To("role_permission", Role_Permission.Type),
	}
}
