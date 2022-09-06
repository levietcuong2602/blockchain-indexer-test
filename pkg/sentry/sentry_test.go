package sentry

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConditionAnd(t *testing.T) {
	tests := []struct {
		name       string
		conditions []Condition
		expected   bool
	}{
		{
			name: "all conditions satisfied",
			conditions: []Condition{
				func(res *http.Response, url string) bool {
					return true
				},
				func(res *http.Response, url string) bool {
					return true
				},
			},
			expected: true,
		},
		{
			name: "all conditions unsatisfied",
			conditions: []Condition{
				func(res *http.Response, url string) bool {
					return false
				},
				func(res *http.Response, url string) bool {
					return false
				},
			},
			expected: false,
		},
		{
			name: "first of two conditions is satisfied",
			conditions: []Condition{
				func(res *http.Response, url string) bool {
					return true
				},
				func(res *http.Response, url string) bool {
					return false
				},
			},
			expected: false,
		},
		{
			name: "second of two conditions is satisfied",
			conditions: []Condition{
				func(res *http.Response, url string) bool {
					return false
				},
				func(res *http.Response, url string) bool {
					return true
				},
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			condition := ConditionAnd(tc.conditions...)
			actual := condition(nil, "")

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestConditionOr(t *testing.T) {
	tests := []struct {
		name       string
		conditions []Condition
		expected   bool
	}{
		{
			name: "two out of two are satisfied",
			conditions: []Condition{
				func(res *http.Response, url string) bool {
					return true
				},
				func(res *http.Response, url string) bool {
					return true
				},
			},
			expected: true,
		},
		{
			name: "first out of two is satisfied",
			conditions: []Condition{
				func(res *http.Response, url string) bool {
					return true
				},
				func(res *http.Response, url string) bool {
					return false
				},
			},
			expected: true,
		},
		{
			name: "second out of two is satisfied",
			conditions: []Condition{
				func(res *http.Response, url string) bool {
					return false
				},
				func(res *http.Response, url string) bool {
					return true
				},
			},
			expected: true,
		},
		{
			name: "none of conditions is satisfied",
			conditions: []Condition{
				func(res *http.Response, url string) bool {
					return false
				},
				func(res *http.Response, url string) bool {
					return false
				},
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			condition := ConditionOr(tc.conditions...)
			actual := condition(nil, "")

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestConditionNotStatusOk(t *testing.T) {
	tests := []struct {
		name     string
		resp     *http.Response
		expected bool
	}{
		{
			name:     "response code is below 200",
			resp:     &http.Response{StatusCode: 100},
			expected: true,
		},
		{
			name:     "response code is 200",
			resp:     &http.Response{StatusCode: 200},
			expected: false,
		},
		{
			name:     "response code is between 200 and 300",
			resp:     &http.Response{StatusCode: 201},
			expected: false,
		},
		{
			name:     "response code is between 300",
			resp:     &http.Response{StatusCode: 300},
			expected: true,
		},
		{
			name:     "response code is above 300",
			resp:     &http.Response{StatusCode: 303},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := ConditionNotStatusOk(tc.resp, "")
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestConditionNotStatusBadRequest(t *testing.T) {
	tests := []struct {
		name     string
		resp     *http.Response
		expected bool
	}{
		{
			name:     "response code is bad request",
			resp:     &http.Response{StatusCode: http.StatusBadRequest},
			expected: false,
		},
		{
			name:     "response code is not bad request",
			resp:     &http.Response{StatusCode: http.StatusOK},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := ConditionNotStatusBadRequest(tc.resp, "")
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestConditionNotStatusNotFound(t *testing.T) {
	tests := []struct {
		name     string
		resp     *http.Response
		expected bool
	}{
		{
			name:     "response code is not found",
			resp:     &http.Response{StatusCode: http.StatusNotFound},
			expected: false,
		},
		{
			name:     "response code is not \"not found\"",
			resp:     &http.Response{StatusCode: http.StatusBadRequest},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := ConditionNotStatusNotFound(tc.resp, "")
			assert.Equal(t, tc.expected, actual)
		})
	}
}
