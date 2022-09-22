package app

import (
	"reflect"
	"testing"

	"github.com/doniyorbek7376/random_string_generator/models"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []models.Token
		wantErr bool
	}{
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:  "only constant values",
			input: "abcd",
			want: []models.Token{
				{Value: "a"},
				{Value: "b"},
				{Value: "c"},
				{Value: "d"},
			},
		},
		{
			name:  "asteriks, dot",
			input: "a*cd.",
			want: []models.Token{
				{Value: "a"},
				{Value: "*", Type: models.Asteriks},
				{Value: "c"},
				{Value: "d"},
				{Value: ".", Type: models.Dot},
			},
		},
		{
			name:  "question mark",
			input: "A&?c",
			want: []models.Token{
				{Value: "A"},
				{Value: "&"},
				{Value: "?", Type: models.QuestionMark},
				{Value: "c"},
			},
		},
		{
			name:  "plus",
			input: "10+1",
			want: []models.Token{
				{Value: "1"},
				{Value: "0"},
				{Value: "+", Type: models.Plus},
				{Value: "1"},
			},
		},
		{
			name:  "alternating branch",
			input: "ab|cd",
			want: []models.Token{
				{Value: "a"},
				{Value: "b"},
				{Value: "|", Type: models.AlternatingBranch},
				{Value: "c"},
				{Value: "d"},
			},
		},
		{
			name:  "min max block on group",
			input: "(ab){0,1}",
			want: []models.Token{
				{Value: "(", Type: models.GroupOpener},
				{Value: "a"},
				{Value: "b"},
				{Value: ")", Type: models.GroupCloser},
				{Value: "{", Type: models.CounterOpener},
				{Value: "0"},
				{Value: ",", Type: models.Comma},
				{Value: "1"},
				{Value: "}", Type: models.CounterCloser},
			},
		},
		{
			name:  "backslashed asteriks",
			input: "ab\\*",
			want: []models.Token{
				{Value: "a"},
				{Value: "b"},
				{Value: "\\", Type: models.BackSlash},
				{Value: "*"},
			},
		},
		{
			name:  "class of chars",
			input: "[^a-z*+]",
			want: []models.Token{
				{Value: "[", Type: models.ClassOpener},
				{Value: "^", Type: models.ClassNegater},
				{Value: "a"},
				{Value: "-", Type: models.ClassRange},
				{Value: "z"},
				{Value: "*"},
				{Value: "+"},
				{Value: "]", Type: models.ClassCloser},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tk := NewTokenizer()
			got, err := tk.Tokenize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("tokenizer.Tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenizer.Tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
