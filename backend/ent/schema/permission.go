package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Permission struct{ ent.Schema }

func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Annotations(),
		field.String("name").Unique(),
		field.String("description"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Default(time.Now).UpdateDefault(time.Now),
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
