package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// RegisteredContract holds the schema definition for the RegisteredContract entity.
type RegisteredContract struct {
	ent.Schema
}

// Fields of the RegisteredContract.
func (RegisteredContract) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").MaxLen(128).MinLen(128).NotEmpty().Unique().Immutable(),
		field.Int("format"), // TODO: replace with Enum?
		field.Bytes("contract").NotEmpty(),
		field.Time("created_at").Immutable().Default(time.Now),
		// field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the RegisteredContract.
func (RegisteredContract) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("providers", RegisteredProvider.Type).Ref("contract"),
		edge.From("compatibility_results_as_requirement", CompatibilityResult.Type).Ref("requirement_contract"),
		edge.From("compatibility_results_as_provider", CompatibilityResult.Type).Ref("provider_contract"),
	}
}
