package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Identity struct{ ent.Schema }

func (Identity) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable(),
		field.UUID("user_id", uuid.UUID{}),
		field.String("provider").NotEmpty(), // "google" | "line" | "facebook" | "apple" | "github" | "oidc"
		field.String("subject").NotEmpty(),  // provider’s user id (sub)
		field.String("email").Optional(),
		field.String("name").Optional(),
		field.String("avatar").Optional(),

		// token storage (encrypt at rest if you keep these!)
		field.String("access_token").Optional().Sensitive(),
		field.String("refresh_token").Optional().Sensitive(),
		field.Time("expires_at").Optional(),

		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Identity) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "provider", "subject").Unique(),
	}
}

func (Identity) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("identities").Unique().Required().Field("user_id"),
	}
}
