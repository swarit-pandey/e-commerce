package database

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// applyQueryFilters will apply all the filters (if given any)
func applyQueryFilters(query *gorm.DB, qf QueryFilter) *gorm.DB {
	// 1. Apply JOINs first as they expand or limit the scope of data to be queried.
	query = applyJoins(query, qf.Joins)

	// 2. Apply DISTINCT if required to ensure unique rows; beneficial post-JOINs.
	query = applyDistinct(query, qf.Distinct)

	// 3. Apply WHERE conditions to filter the dataset.
	query = applyWhereConditions(query, qf)

	// 4. Applying regex conditions for advanced pattern matching.
	query = applyRegexConditions(query, qf.RegexFilter)

	// 5. GROUP BY to aggregate results based on specified columns.
	query = applyGroupBy(query, qf.GroupBy)

	// 6. HAVING conditions to filter groups.
	query = applyHavingConditions(query, qf)

	// 7. ORDER BY to sort the results.
	query = applyOrderBy(query, qf.OrderBy)

	// 8. Limit and Offset for pagination.
	query = applyLimitOffset(query, qf.Limit, qf.Offset)

	// 9. OR conditions for alternative matching.
	query = applyOrConditions(query, qf.Or)

	// 10. Pluck is used typically at the end to select specific fields.
	query = applyPluck(query, qf.Pluck)

	return query
}

func applyJoins(query *gorm.DB, joins []Join) *gorm.DB {
	for _, join := range joins {
		joinQuery := fmt.Sprintf("JOIN %s ON %s", join.Table, join.Condition) // Default JOIN
		if join.Type != "" {
			joinQuery = fmt.Sprintf("%s JOIN %s ON %s", strings.ToUpper(join.Type), join.Table, join.Condition)
		}
		query = query.Joins(joinQuery, join.Args...)
	}
	return query
}

func applyWhereConditions(query *gorm.DB, qf QueryFilter) *gorm.DB {
	if len(qf.WhereMap) > 0 {
		for key, value := range qf.WhereMap {
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
		}

		if len(qf.UpdateMap) > 0 {
			for key, value := range qf.UpdateMap {
				query = query.Update(key, value)
			}
		}
	}

	if qf.WhereQuery != "" {
		query = query.Where(qf.WhereQuery, qf.WhereArgs...)

		if len(qf.UpdateMap) > 0 {
			for key, value := range qf.UpdateMap {
				query = query.Update(key, value)
			}
		}
	}

	return query
}

func applyHavingConditions(query *gorm.DB, qf QueryFilter) *gorm.DB {
	if len(qf.HavingMap) > 0 {
		for key, value := range qf.HavingMap {
			query = query.Having(fmt.Sprintf("%s = ?", key), value)
		}
	}

	if qf.HavingQuery != "" {
		query = query.Having(qf.HavingQuery, qf.HavingArgs...)
	}

	return query
}

func applyOrderBy(query *gorm.DB, orderBy string) *gorm.DB {
	if orderBy != "" {
		query = query.Order(orderBy)
	}
	return query
}

func applyGroupBy(query *gorm.DB, groupBy string) *gorm.DB {
	if groupBy != "" {
		query = query.Group(groupBy)
	}
	return query
}

func applyLimitOffset(query *gorm.DB, limit, offset int) *gorm.DB {
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	return query
}

func applyRegexConditions(query *gorm.DB, regexFilter map[string]string) *gorm.DB {
	for field, regex := range regexFilter {
		query = query.Where(fmt.Sprintf("%s REGEXP ?", field), regex)
	}
	return query
}

func applyDistinct(query *gorm.DB, distinct bool) *gorm.DB {
	if distinct {
		query = query.Distinct()
	}
	return query
}

func applyPluck(query *gorm.DB, pluck string) *gorm.DB {
	if pluck != "" {
		query = query.Pluck(pluck, nil)
	}
	return query
}

func applyOrConditions(query *gorm.DB, orConditions []QueryFilter) *gorm.DB {
	for _, condition := range orConditions {
		query = query.Or(applyQueryFilters(query, condition))
	}
	return query
}

// applyProjections will project the required columns as provided by clients
func applyProjections(query *gorm.DB, p *Projection) *gorm.DB {
	if p != nil {
		var projections []string

		// Apply field projections
		if len(p.Fields) > 0 {
			projections = append(projections, p.Fields...)
		}

		// Apply aggregate projections
		for _, aggregate := range p.Aggregates {
			projection := fmt.Sprintf("%s(%s) AS %s", aggregate.Function, aggregate.Field, aggregate.Alias)
			projections = append(projections, projection)
		}

		// Apply conditional projections
		for _, conditional := range p.Conditionals {
			projection := fmt.Sprintf("CASE WHEN %s THEN %s ELSE %s END AS %s", conditional.Condition, conditional.TrueValue, conditional.FalseValue, conditional.Alias)
			projections = append(projections, projection)
		}

		// Apply subquery projections
		for _, subquery := range p.Subqueries {
			projection := fmt.Sprintf("(%s) AS %s", subquery.Query, subquery.Alias)
			projections = append(projections, projection)
		}

		// Apply projections to the query
		query = query.Select(projections)

		// Apply field exclusions
		if len(p.Exclude) > 0 {
			excludeFields := make([]string, len(p.Exclude))
			for i, field := range p.Exclude {
				excludeFields[i] = "-" + field
			}
			query = query.Omit(excludeFields...)
		}
	}
	return query
}
