package slice

import (
	"reflect"
	"testing"
)

func TestCompareStringSlice(t *testing.T) {
	type test []struct {
		s1   []string
		s2   []string
		want bool
	}
	tests := test{
		{s1: []string{"a", "b", "c"}, s2: []string{"a", "b", "c"}, want: true},
		{s1: []string{}, s2: []string{}, want: true},
		{s1: []string{"abc", "b", "a"}, s2: []string{"a", "b", "abc"}, want: true},
		{s1: []string{"a", "b"}, s2: []string{"a", "b", "c"}, want: false},
		{s1: []string{"a", "b", "c"}, s2: []string{"a", "b", "d"}, want: false},
	}
	for _, tc := range tests {
		got := CompareStringSlice(tc.s1, tc.s2)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("compare error, s1: %#v, s2: %#v", tc.s1, tc.s2)
		}
	}
}

func TestCompareIntSlice(t *testing.T) {
	type test []struct {
		i1   []int
		i2   []int
		want bool
	}
	tests := test{
		{i1: []int{1, 2, 3}, i2: []int{1, 2, 3}, want: true},
		{i1: []int{}, i2: []int{}, want: true},
		{i1: []int{123, 2, 1}, i2: []int{1, 2, 123}, want: true},
		{i1: []int{1, 2}, i2: []int{1, 2, 3}, want: false},
		{i1: []int{1, 2, 3}, i2: []int{1, 2, 4}, want: false},
	}
	for _, tc := range tests {
		got := CompareIntSlice(tc.i1, tc.i2)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("compare error, i1: %#v, s2: %#v", tc.i1, tc.i2)
		}
	}
}
