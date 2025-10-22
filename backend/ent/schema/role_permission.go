package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

// Role_Permission holds the schema definition for the Role_Permission entity.
type Role_Permission struct {
	ent.Schema
}

// Fields of the Role_Permission.
func (Role_Permission) Fields() []ent.Field {
	return []ent.Field{}
}

// Edges of the Role_Permission.
func (Role_Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", Role.Type).Ref("role_permission"),
		edge.From("permission", Permission.Type).Ref("role_permission"),
		edge.From("granted_by", User.Type).Ref("role_permission"),
	}
}
