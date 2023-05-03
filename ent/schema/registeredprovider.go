package schema

import (
	"net/url"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// RegisteredProvider holds the schema definition for the RegisteredProvider entity.
type RegisteredProvider struct {
	ent.Schema
}

// Fields of the RegisteredProvider.
func (RegisteredProvider) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.JSON("url", &url.URL{}),
		field.String("participant_name").NotEmpty(),
		field.String("contract_id").MaxLen(128).MinLen(128),
		field.Time("created_at").Immutable().Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}

}

// Edges of the RegisteredProvider.
func (RegisteredProvider) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("contract", RegisteredContract.Type).Field("contract_id").Unique().Required(),
	}
}
