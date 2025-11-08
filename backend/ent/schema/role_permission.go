package schema

import (
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
	}
}

// Edges of the Role_Permission.
func (Role_Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", Role.Type).Ref("role_permissions").Unique().Required(),
		edge.From("permission", Permission.Type).Ref("role_permissions").Unique().Required(),
	}
}
