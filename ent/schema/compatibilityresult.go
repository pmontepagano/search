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

// Fields of the CompatibilityResult.
func (CompatibilityResult) Fields() []ent.Field {
	return []ent.Field{
		field.String("req_contract_id").MaxLen(128).MinLen(128).NotEmpty(),
		field.String("prov_contract_id").MaxLen(128).MinLen(128).NotEmpty(),
		field.String("participant_name_req").NotEmpty(),
		field.String("participant_name_prov").NotEmpty(),
		field.Bool("result"),

		field.Time("created_at").Immutable().Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the CompatibilityResult.
func (CompatibilityResult) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("requirement_contract", RegisteredContract.Type).Required().Field("req_contract_id").Unique(),
		edge.To("provider_contract", RegisteredContract.Type).Required().Field("prov_contract_id").Unique(),
	}
}

func (CompatibilityResult) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("participant_name_req", "participant_name_prov", "req_contract_id", "prov_contract_id").Unique(),
	}
}
