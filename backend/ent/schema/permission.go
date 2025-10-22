package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Permission struct{ ent.Schema }

func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.String("name").Unique(),
		field.String("description"),
	}
}

func Index() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}

func (Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("role_permission", Role_Permission.Type),
	}
}
