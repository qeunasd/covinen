package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var maxPageSize = 100

var listFilters = map[string]allowedFilter{
	"dmin":   {column: "tgl_dibuat", operator: "gte"},
	"dmax":   {column: "tgl_dibuat", operator: "lte"},
	"status": {column: "status", operator: "in"},
	"jr":     {column: "jumlah_ruangan", operator: "eq"},
}

type TableConfig struct {
	QueryCols   []string
	SortCols    []string
	DefaultSort string
}

type allowedFilter struct {
	column, operator string
}

type PaginationResult struct {
	Data      any
	TotalData int64
	TotalPage int
	Page      int
	PerPage   int
}

type filter struct {
	field    string
	operator string
	value    any
}

type PaginationParams struct {
	Page      int
	PerPage   int
	SortBy    string
	SortDir   string
	Query     string
	Filters   []filter
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

func PaginationFromRequest(r *http.Request) PaginationParams {
	q := r.URL.Query()

	page, err := strconv.Atoi(q.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(q.Get("perpage"))
	if err != nil || perPage < 1 || perPage > maxPageSize {
		perPage = 10
	}

	sortDir := q.Get("ord")
	if sortDir != "desc" && sortDir != "asc" {
		sortDir = "desc"
	}

	var filters []filter
	for key, val := range q {
		if v, ok := listFilters[key]; ok {
			filters = append(filters, filter{
				field: v.column, operator: v.operator, value: extractValue(v.operator, val)},
			)
		}
	}

	return PaginationParams{
		Page:    page,
		PerPage: perPage,
		SortBy:  q.Get("sb"),
		SortDir: sortDir,
		Query:   q.Get("q"),
		Filters: filters,
	}
}

func FiltersToMap(filters []filter) map[string]any {
	m := make(map[string]any)
	for _, f := range filters {
		m[f.field] = f.value
	}
	return m
}

func extractValue(op string, vals []string) any {
	if len(vals) == 0 {
		return nil
	}

	if op == "in" || op == "nin" {
		result := make([]any, len(vals))
		for i, v := range vals {
			result[i] = v
		}
		return result
	}

	return vals[0]
}

func BuildWhereClauses(params PaginationParams) (whereClause string, arguments []any) {
	var conditions []string
	argIndex := 1

	if len(params.Filters) > 0 {
		for _, f := range params.Filters {
			switch f.operator {
			case "in", "nin":
				op := "IN"
				if f.operator == "nin" {
					op = "NIN"
				}
				values, ok := f.value.([]any)
				if !ok || len(values) == 0 {
					continue
				}
				placeholders := []string{}
				for _, value := range values {
					arguments = append(arguments, value)
					placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
					argIndex++
				}
				conditions = append(conditions,
					fmt.Sprintf("%s %s (%s)", f.field, op, strings.Join(placeholders, ", ")))
			default:
				arguments = append(arguments, f.value)
				conditions = append(conditions,
					fmt.Sprintf("%s %s $%d", f.field, mapOperator(f.operator), argIndex))
				argIndex++
			}
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

	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	return
}

func BuildSortClause(params PaginationParams) string {
	if params.SortBy == "" {
		return ""
	}
	return fmt.Sprintf(" ORDER BY %s %s", params.SortBy, strings.ToUpper(params.SortDir))
}

func Contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
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
