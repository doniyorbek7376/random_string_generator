package app

import (
	"math/rand"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name    string
		regex   string
		wantErr bool
	}{
		{
			name:  "simple text",
			regex: "hello world",
		},
		{
			name:  "numbers",
			regex: "[-+]?[0-9]{1,16}[.][0-9]{1,6}",
		},
		{
			name:  "uuids",
			regex: "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{8}",
		},
		{
			name:  "min max block with dot",
			regex: ".{8,12}",
		},
		{
			name:  "negated class",
			regex: "[^aeiouAEIOU0-9]{5}",
		},
		{
			name:  "class with special characters",
			regex: "[a-f&-]{5}",
		},
		{
			name:  "bonus: alternating branch, sub patterns",
			regex: "(1[0-2]|0[1-9])(:[0-5][0-9]){2} (A|P)M",
		},
		{
			name:  "back reference",
			regex: "(.{5})\\1\\1",
		},
	}
	g := NewGenerator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, err := g.Generate(tt.regex, rand.Intn(10)+1)
			if (err != nil) != tt.wantErr {
				t.Errorf("generator.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				for _, value := range values {
					ok, _ := regexp.MatchString(tt.regex, value)
					assert.True(t, ok)
				}
			}
		})
	}
}
