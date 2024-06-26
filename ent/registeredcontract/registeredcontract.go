// Code generated by ent, DO NOT EDIT.

package registeredcontract

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the registeredcontract type in the database.
	Label = "registered_contract"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldFormat holds the string denoting the format field in the database.
	FieldFormat = "format"
	// FieldContract holds the string denoting the contract field in the database.
	FieldContract = "contract"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// EdgeProviders holds the string denoting the providers edge name in mutations.
	EdgeProviders = "providers"
	// EdgeCompatibilityResultsAsRequirement holds the string denoting the compatibility_results_as_requirement edge name in mutations.
	EdgeCompatibilityResultsAsRequirement = "compatibility_results_as_requirement"
	// EdgeCompatibilityResultsAsProvider holds the string denoting the compatibility_results_as_provider edge name in mutations.
	EdgeCompatibilityResultsAsProvider = "compatibility_results_as_provider"
	// Table holds the table name of the registeredcontract in the database.
	Table = "registered_contracts"
	// ProvidersTable is the table that holds the providers relation/edge.
	ProvidersTable = "registered_providers"
	// ProvidersInverseTable is the table name for the RegisteredProvider entity.
	// It exists in this package in order to avoid circular dependency with the "registeredprovider" package.
	ProvidersInverseTable = "registered_providers"
	// ProvidersColumn is the table column denoting the providers relation/edge.
	ProvidersColumn = "contract_id"
	// CompatibilityResultsAsRequirementTable is the table that holds the compatibility_results_as_requirement relation/edge.
	CompatibilityResultsAsRequirementTable = "compatibility_results"
	// CompatibilityResultsAsRequirementInverseTable is the table name for the CompatibilityResult entity.
	// It exists in this package in order to avoid circular dependency with the "compatibilityresult" package.
	CompatibilityResultsAsRequirementInverseTable = "compatibility_results"
	// CompatibilityResultsAsRequirementColumn is the table column denoting the compatibility_results_as_requirement relation/edge.
	CompatibilityResultsAsRequirementColumn = "requirement_contract_id"
	// CompatibilityResultsAsProviderTable is the table that holds the compatibility_results_as_provider relation/edge.
	CompatibilityResultsAsProviderTable = "compatibility_results"
	// CompatibilityResultsAsProviderInverseTable is the table name for the CompatibilityResult entity.
	// It exists in this package in order to avoid circular dependency with the "compatibilityresult" package.
	CompatibilityResultsAsProviderInverseTable = "compatibility_results"
	// CompatibilityResultsAsProviderColumn is the table column denoting the compatibility_results_as_provider relation/edge.
	CompatibilityResultsAsProviderColumn = "provider_contract_id"
)

// Columns holds all SQL columns for registeredcontract fields.
var Columns = []string{
	FieldID,
	FieldFormat,
	FieldContract,
	FieldCreatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// ContractValidator is a validator for the "contract" field. It is called by the builders before save.
	ContractValidator func([]byte) error
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// IDValidator is a validator for the "id" field. It is called by the builders before save.
	IDValidator func(string) error
)

// OrderOption defines the ordering options for the RegisteredContract queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByFormat orders the results by the format field.
func ByFormat(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFormat, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByProvidersCount orders the results by providers count.
func ByProvidersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newProvidersStep(), opts...)
	}
}

// ByProviders orders the results by providers terms.
func ByProviders(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newProvidersStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByCompatibilityResultsAsRequirementCount orders the results by compatibility_results_as_requirement count.
func ByCompatibilityResultsAsRequirementCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newCompatibilityResultsAsRequirementStep(), opts...)
	}
}

// ByCompatibilityResultsAsRequirement orders the results by compatibility_results_as_requirement terms.
func ByCompatibilityResultsAsRequirement(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCompatibilityResultsAsRequirementStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByCompatibilityResultsAsProviderCount orders the results by compatibility_results_as_provider count.
func ByCompatibilityResultsAsProviderCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newCompatibilityResultsAsProviderStep(), opts...)
	}
}

// ByCompatibilityResultsAsProvider orders the results by compatibility_results_as_provider terms.
func ByCompatibilityResultsAsProvider(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCompatibilityResultsAsProviderStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newProvidersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProvidersInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, ProvidersTable, ProvidersColumn),
	)
}
func newCompatibilityResultsAsRequirementStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CompatibilityResultsAsRequirementInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, CompatibilityResultsAsRequirementTable, CompatibilityResultsAsRequirementColumn),
	)
}
func newCompatibilityResultsAsProviderStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CompatibilityResultsAsProviderInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, CompatibilityResultsAsProviderTable, CompatibilityResultsAsProviderColumn),
	)
}
