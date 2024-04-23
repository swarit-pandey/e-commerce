package database

// Projection defines how data should be projected from the database,
// including specific fields, aggregates, conditionals, and subqueries.
type Projection struct {
	Fields       []string                // Fields to include in the projection.
	Aggregates   []AggregateProjection   // Aggregate functions to apply.
	Conditionals []ConditionalProjection // Conditional field values.
	Subqueries   []SubqueryProjection    // Subqueries to include in the projection.
	Exclude      []string                // Fields to exclude.
}

// AggregateProjection specifies an aggregation operation on a field.
type AggregateProjection struct {
	Field    string // Field to aggregate.
	Function string // Aggregate function (e.g., "SUM", "COUNT", "AVG").
	Alias    string // Optional alias for the aggregated result.
}

// ConditionalProjection allows for conditional logic in field values.
type ConditionalProjection struct {
	Condition  string // SQL condition for the projection.
	TrueValue  string // Value or field when the condition is true.
	FalseValue string // Value or field when the condition is false.
	Alias      string // Alias for the conditional projection result.
}

// SubqueryProjection enables including subqueries in the projection.
type SubqueryProjection struct {
	Query string // The subquery itself.
	Args  []any  // Arguments for placeholders within the subquery.
	Alias string // Alias for the subquery result in the projection.
}
