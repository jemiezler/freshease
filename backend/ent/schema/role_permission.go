package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Role_Permission holds the schema definition for the Role_Permission entity.
type Role_Permission struct {
	ent.Schema
}

// Fields of the Role_Permission.
func (Role_Permission) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Role_Permission.
func (Role_Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", Role.Type).Ref("role_permission"),
		edge.From("permission", Permission.Type).Ref("role_permission"),
		edge.From("granted_by", User.Type).Ref("role_permission"),
	}
}
