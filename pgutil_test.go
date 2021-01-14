package pgutil

import (
	"errors"
	"fmt"
	"testing"
)

func add(a, b int) int {
	return a + b
}

func TestParseToken(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name string
		args args
		want []string
		err  error
	}{
		{
			name: "normal",
			args: args{sql: `
CREATE TABLE public.orders (
	id integer NOT NULL,
	donor_id integer,
	ordered_at timestamp with time zone DEFAULT now()
)
			`},
			want: []string{
				"create", "table", "public.orders", "(",
				"id", "integer", "not", "null", ",",
				"donor_id", "integer", ",",
				"ordered_at", "timestamp", "with", "time", "zone",
				"default", "now", "(", ")", ")",
			},
			err: nil,
		},
		{
			name: "normal",
			args: args{sql: `
CREATE TABLE public.b (
	id integer NOT NULL
)
			`},
			want: []string{
				"create", "table", "public.b", "(",
				"id", "integer", "not", "null", ")",
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.args.sql)
			if err != tt.err {
				t.Errorf("Tokenize err = %v, want %v", err, tt.err)
			}
			if err1 := stringSliceEql(got, tt.want); err1 != nil {
				t.Errorf("Tokenize got = %v, want = %v err1=%v\n", got, tt.want, err1)
			}
		})
	}
}

func stringSliceEql(a, b []string) error {
	if len(a) != len(b) {
		return errors.New("slice length not equal")
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return fmt.Errorf("slice index(%d) not equal(%s != %s)", i, a[i], b[i:])
		}
	}
	return nil
}
