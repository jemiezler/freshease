package schema

import (
	"time"

	"github.com/google/uuid"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Review struct{ ent.Schema }

func (Review) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Int("rating"),
		field.String("comment").Nillable().Optional(),
		field.Time("created_at").Default(time.Now),
	}
}

func (Review) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("reviews").Required(),
		edge.From("product", Product.Type).Ref("reviews").Required(),
	}
}

