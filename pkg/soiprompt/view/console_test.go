package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected SoiLine
	}{
		{
			name: "normal",
			s:    "15ea9676 [  2] blogs/Designing_Edge_Gateway,_Uber’s_API_Lifecycle_Management_Platform_-_Uber_Engineering_Blog [#aaa #bbb=ccc]",
			expected: SoiLine{
				Hash:     "15ea9676",
				NumViews: 2,
				Path:     "blogs/Designing_Edge_Gateway,_Uber’s_API_Lifecycle_Management_Platform_-_Uber_Engineering_Blog",
				Tags:     []string{"#aaa", "#bbb=ccc"},
			},
		},
		{
			name: "no-tags",
			s:    "15ea9676 [  2] blogs/Designing_Edge_Gateway,_Uber’s_API_Lifecycle_Management_Platform_-_Uber_Engineering_Blog []",
			expected: SoiLine{
				Hash:     "15ea9676",
				NumViews: 2,
				Path:     "blogs/Designing_Edge_Gateway,_Uber’s_API_Lifecycle_Management_Platform_-_Uber_Engineering_Blog",
				Tags:     nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l, err := ParseLine(test.s)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, *l)
		})
	}

}
