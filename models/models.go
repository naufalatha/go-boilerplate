/*
Models package contains all data model used in handler.
Most of the data model is json mapped and used for returning response data.
*/

package models

import (
	"fmt"
)

// Environment
const (
	ENV_LOCAL       = "LOCAL"
	ENV_DEVELOPMENT = "DEVELOPMENT"
	ENV_STAGING     = "STAGING"
	ENV_PRODUCTION  = "PRODUCTION"
)

// Status Code Mapping
const (
	CODE_WRONG_ARGS        = "WRONG_ARGS"
	CODE_NOT_FOUND         = "NOT_FOUND"
	CODE_INTERNAL_ERROR    = "INTERNAL_ERROR"
	CODE_UNAUTHORIZED      = "UNAUTHORIZED"
	CODE_FORBIDDEN         = "FORBIDDEN"
	CODE_INVALID_MIME_TYPE = "INVALID_MIME_TYPE"
)

var (
	ErrNoRows = fmt.Errorf("qrm: no rows in result set")
)

const DEFAULT_TIME_FORMAT = "2006-01-02 15:04:05"

type PaginationFilter struct {
	Page  int64
	Limit int64
}

type Pagination struct {
	Total int64
	Data  interface{}
}

type Total struct {
	Total int64 `alias:"total"`
}

type Response struct {
	Fields     map[string]string `json:"fields,omitempty"`
	Status     string            `json:"status,omitempty"`
	Success    bool              `json:"success,omitempty"`
	Message    string            `json:"message,omitempty"`
	StatusCode int               `json:"status_code,omitempty"`
	Data       interface{}       `json:"data,omitempty"`
}

func (filter *PaginationFilter) ParsePagination(datas interface{}, total int64) Pagination {
	if filter.Limit == 0 {
		filter.Limit = 10
	}

	return Pagination{
		Total: total,
		Data:  datas,
	}
}

type LogProvider struct {
	Url         string      `json:"url,omitempty"`
	Headers     interface{} `json:"headers,omitempty"`
	RequestBody interface{} `json:"request_body,omitempty"`
	// RequestBody *string     `json:"request_body,omitempty"`
	Response   interface{} `json:"response,omitempty"`
	StatusCode int         `json:"status_code,omitempty"`
}
