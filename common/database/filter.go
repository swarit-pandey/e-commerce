package database

// QueryFilter provides a clean model for clients to interact with
// database and abstracts away the DB specific concerns.
type QueryFilter struct {
	// WhereMap is used for equality conditions on fields.
	WhereMap map[string]any

	// WhereQuery is a custom WHERE clause with placeholders.
	WhereQuery string

	// WhereArgs are the arguments for placeholders in WhereQuery.
	WhereArgs []any

	// HavingMap is used for HAVING conditions on fields, typically with GROUP BY.
	HavingMap map[string]any

	// HavingQuery is a custom HAVING clause with placeholders, used with GROUP BY.
	HavingQuery string

	// HavingArgs are the arguments for placeholders in HavingQuery.
	HavingArgs []any

	// UpdateMap is used for specifying fields to update in an UPDATE operation.
	UpdateMap map[string]any

	// OrderBy specifies the ordering of the result set.
	OrderBy string

	// GroupBy specifies the grouping of the result set.
	GroupBy string

	// Limit specifies the maximum number of records to return.
	Limit int

	// Offset specifies the number of records to skip before starting to return records.
	Offset int

	// Joins specifies the tables and conditions for JOIN operations.
	Joins []Join

	// RegexFilter specifies regex conditions for fields.
	RegexFilter map[string]string

	// Distinct specifies if only distinct records should be returned.
	Distinct bool

	// Pluck allows selecting a single column from the results.
	Pluck string

	// Or specifies OR conditions, allowing for complex query building.
	Or []QueryFilter
}

// Join represents a JOIN operation in SQL, specifying the table, condition, and arguments.
type Join struct {
	Table     string // The table to join.
	Condition string // The JOIN condition.
	Args      []any  // Arguments for placeholders in Condition.
	Type      string // Type of join: INNER, LEFT, RIGHT, OUTER, etc.
}
