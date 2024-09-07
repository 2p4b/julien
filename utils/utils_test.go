package utils

import "testing"

func TestArrayIncludesInt(t *testing.T) {
	type args struct {
		values []int
		target int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Int target in values",
			args: args{
				values: []int{1, 2, 3},
				target: 3,
			},
			want: true,
		},
		{
			name: "Int target not in values",
			args: args{
				values: []int{1, 2, 3},
				target: 4,
			},
			want: false,
		},
		{
			name: "Empty int values",
			args: args{
				values: []int{},
				target: 2,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayIncludes(tt.args.values, tt.args.target); got != tt.want {
				t.Errorf("ArrayIncludes[Int]() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayIncludesString(t *testing.T) {
	type args struct {
		values []string
		target string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "String target in values",
			args: args{
				values: []string{"apple", "banana", "cherry"},
				target: "banana",
			},
			want: true,
		},
		{
			name: "String target not in values",
			args: args{
				values: []string{"apple", "banana", "cherry"},
				target: "grape",
			},
			want: false,
		},
		{
			name: "Empty string values",
			args: args{
				values: []string{},
				target: "grape",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayIncludes(tt.args.values, tt.args.target); got != tt.want {
				t.Errorf("ArrayIncludes[string]() = %v, want %v", got, tt.want)
			}
		})
	}
}
