package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// CompatibilityResult holds the schema definition for the CompatibilityResult entity.
type CompatibilityResult struct {
	ent.Schema
}

type ParticipantNameMapping map[string]string

// Fields of the CompatibilityResult.
func (CompatibilityResult) Fields() []ent.Field {
	return []ent.Field{
		field.String("requirement_contract_id").MaxLen(128).MinLen(128),
		field.String("provider_contract_id").MaxLen(128).MinLen(128),
		field.Bool("result"),
		field.JSON("mapping", ParticipantNameMapping{}).Optional(),

		field.Time("created_at").Immutable().Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the CompatibilityResult.
func (CompatibilityResult) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("requirement_contract", RegisteredContract.Type).Required().Unique().Field("requirement_contract_id"),
		edge.To("provider_contract", RegisteredContract.Type).Required().Unique().Field("provider_contract_id"),
	}
}

func (CompatibilityResult) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("requirement_contract_id", "provider_contract_id").Unique(),
	}
}
