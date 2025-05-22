package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var MaxPageSize = 100

type allowedFilter struct {
	column   string
	operator string
}

var filterConfigs = map[string]allowedFilter{
	"dmin":   {column: "tgl_dibuat", operator: "gte"},
	"dmax":   {column: "tgl_dibuat", operator: "lte"},
	"status": {column: "status", operator: "in"},
	"jr":     {column: "jumlah_ruangan", operator: "eq"},
}

type TableConfig struct {
	QueryCols   []string
	SortCols    []AllowedSort
	DefaultSort string
}

type AllowedSort struct {
	Name   string
	Column string
}

func (tc TableConfig) GetSortColumn(sortBy string) (string, bool) {
	for _, sc := range tc.SortCols {
		if sc.Name == sortBy {
			return sc.Column, true
		}
	}
	return "", false
}

type PaginationResult struct {
	Data      interface{}
	TotalData int64
	TotalPage int
	Page      int
	PerPage   int
}

type Filter struct {
	Field    string
	Operator string
	Value    interface{}
}

type PaginationParams struct {
	Page      int
	PerPage   int
	SortBy    string
	SortDir   string
	Query     string
	Filters   []Filter
	QueryCols []string
}

func (p PaginationParams) getOffset() int {
	return (p.Page - 1) * p.PerPage
}

func (p PaginationParams) getLimit() int {
	return p.PerPage
}

func (p *PaginationParams) SetColumnSearch(cols ...string) {
	p.QueryCols = cols
}

func PaginationFromRequest(r *http.Request) (PaginationParams, error) {
	q := r.URL.Query()

	page, err := parsePage(q.Get("page"))
	if err != nil {
		return PaginationParams{}, fmt.Errorf("invalid page parameter: %w", err)
	}

	perPage, err := parsePerPage(q.Get("perpage"))
	if err != nil {
		return PaginationParams{}, fmt.Errorf("invalid per-page parameter: %w", err)
	}

	sortDir, err := validateSortDirection(q.Get("ord"))
	if err != nil {
		return PaginationParams{}, fmt.Errorf("invalid sort direction: %w", err)
	}

	filters, err := buildFilters(q)
	if err != nil {
		return PaginationParams{}, fmt.Errorf("invalid filters: %w", err)
	}

	return PaginationParams{
		Page:    page,
		PerPage: perPage,
		SortBy:  q.Get("sb"),
		SortDir: sortDir,
		Query:   q.Get("q"),
		Filters: filters,
	}, nil
}

func parsePage(pageStr string) (int, error) {
	if pageStr == "" {
		return 1, nil
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 0, fmt.Errorf("page must be a positive integer")
	}
	return page, nil
}

func parsePerPage(perPageStr string) (int, error) {
	if perPageStr == "" {
		return 10, nil
	}
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 || perPage > MaxPageSize {
		return 0, fmt.Errorf("per-page must be between 1 and %d", MaxPageSize)
	}
	return perPage, nil
}

func validateSortDirection(sortDir string) (string, error) {
	if sortDir == "" {
		return "desc", nil
	}
	if sortDir != "asc" && sortDir != "desc" {
		return "", fmt.Errorf("sort direction must be 'asc' or 'desc'")
	}
	return sortDir, nil
}

func buildFilters(query url.Values) ([]Filter, error) {
	var filters []Filter
	for key, vals := range query {
		if cfg, ok := filterConfigs[key]; ok {
			if len(vals) == 0 {
				continue
			}
			filters = append(filters, Filter{
				Field:    cfg.column,
				Operator: cfg.operator,
				Value:    extractFilterValue(cfg.operator, vals),
			})
		}
	}
	return filters, nil
}

func FiltersToMap(filters []Filter) map[string]interface{} {
	m := make(map[string]interface{})
	for _, f := range filters {
		m[f.Field] = f.Value
	}
	return m
}

func extractFilterValue(operator string, vals []string) interface{} {
	if len(vals) == 0 {
		return nil
	}
	if operator == "in" || operator == "nin" {
		result := make([]interface{}, len(vals))
		for i, v := range vals {
			result[i] = v
		}
		return result
	}
	return vals[0]
}

func BuildWhereClauses(params PaginationParams) (string, []interface{}) {
	var conditions []string
	var arguments []interface{}
	argIndex := 1

	for _, f := range params.Filters {
		switch f.Operator {
		case "in", "nin":
			values, ok := f.Value.([]any)
			if !ok || len(values) == 0 {
				continue
			}

			op := "IN"
			if f.Operator == "nin" {
				op = "NIN"
			}
			placeholders := make([]string, len(values))
			for i, value := range values {
				arguments = append(arguments, value)
				placeholders[i] = fmt.Sprintf("$%d", argIndex)
				argIndex++
			}
			conditions = append(conditions, fmt.Sprintf("%s %s (%s)", f.Field, op, strings.Join(placeholders, ", ")))
		default:
			arguments = append(arguments, f.Value)
			conditions = append(conditions, fmt.Sprintf("%s %s $%d", f.Field, mapOperator(f.Operator), argIndex))
			argIndex++
		}
	}

	if params.Query != "" && len(params.QueryCols) > 0 {
		var searchConditions []string
		for _, col := range params.QueryCols {
			arguments = append(arguments, "%"+params.Query+"%")
			searchConditions = append(searchConditions, fmt.Sprintf("%s ILIKE $%d", col, argIndex))
			argIndex++
		}
		conditions = append(conditions, "("+strings.Join(searchConditions, " OR ")+")")
	}

	if len(conditions) == 0 {
		return "", arguments
	}

	return " WHERE " + strings.Join(conditions, " AND "), arguments
}

func BuildSortClause(params PaginationParams, config TableConfig) string {
	col, ok := config.GetSortColumn(params.SortBy)
	if !ok {
		col, _ = config.GetSortColumn(config.DefaultSort)
	}
	return fmt.Sprintf(" ORDER BY %s %s", col, strings.ToUpper(params.SortDir))
}

func BuildLimitClause(params PaginationParams) string {
	return fmt.Sprintf(" LIMIT %d OFFSET %d", params.getLimit(), params.getOffset())
}

func mapOperator(op string) string {
	switch op {
	case "eq":
		return "="
	case "neq":
		return "!="
	case "gte":
		return ">="
	case "lte":
		return "<="
	default:
		return "="
	}
}
