package placeholdedvalues_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mondegor/go-storage/mrpostgres/stream/placeholdedvalues"
)

func TestSQL_CountLineArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		countArgs int
		want      int
	}{
		{
			name:      "test1",
			countArgs: 0,
			want:      1,
		},
		{
			name:      "test2",
			countArgs: 1,
			want:      1,
		},
		{
			name:      "test3",
			countArgs: 5,
			want:      5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sqlHelper := placeholdedvalues.New(
				placeholdedvalues.WithCountLineArgs(tt.countArgs),
			)

			assert.Equal(t, tt.want, sqlHelper.CountLineArgs())
		})
	}
}

func TestSQL_WriteFirstLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		countArgs int
		spans     []string
		want      string
	}{
		{
			name:      "test1",
			countArgs: 1,
			spans:     nil,
			want:      "($1)",
		},
		{
			name:      "test2",
			countArgs: 2,
			spans:     nil,
			want:      "($1, $2)",
		},
		{
			name:      "test3",
			countArgs: 2,
			spans:     []string{"left_arg1::"},
			want:      "(left_arg1::$1, $2)",
		},
		{
			name:      "test4",
			countArgs: 2,
			spans:     []string{"", "", "::right_arg2"},
			want:      "($1, $2::right_arg2)",
		},
		{
			name:      "test5",
			countArgs: 2,
			spans:     []string{"left_arg1::", "::right_arg1, left_arg2::", "::right_arg2, NOW()"},
			want:      "(left_arg1::$1::right_arg1, left_arg2::$2::right_arg2, NOW())",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := strings.Builder{}

			sqlHelper := placeholdedvalues.New(
				placeholdedvalues.WithCountLineArgs(tt.countArgs),
				placeholdedvalues.WithLine(tt.spans...),
			)

			sqlHelper.WriteFirstLine(&buf)

			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestSQL_WriteNextLine(t *testing.T) {
	t.Parallel()

	buf := strings.Builder{}

	sqlHelper := placeholdedvalues.New(
		placeholdedvalues.WithCountLineArgs(3),
		placeholdedvalues.WithLine("", "::x, y::", "::y, z::"),
	)

	number := sqlHelper.WriteFirstLine(&buf)
	number = sqlHelper.WriteNextLine(&buf, number)
	sqlHelper.WriteNextLine(&buf, number)

	assert.Equal(t, "($1::x, y::$2::y, z::$3), ($2::x, y::$3::y, z::$4), ($3::x, y::$4::y, z::$5)", buf.String())
}
