package set

import (
	"testing"

	"github.com/data0007/golib/slice"
)

func TestStringSet(t *testing.T) {
	type test struct {
		got              []string
		addResult        []string
		removeResult     []string
		hasResult        bool
		countResult      int
		unionResult      []string
		minusResult      []string
		intersectResult  []string
		complementResult []string
		clearResult      []string
		emptyResult      bool
	}
	test1 := test{
		got:              []string{"a", "b", "c"},
		addResult:        []string{"a", "b", "c", "d"},
		removeResult:     []string{"a", "b", "c"},
		hasResult:        true,
		countResult:      3,
		unionResult:      []string{"a", "b", "c", "d"},
		minusResult:      []string{"c"},
		intersectResult:  []string{"a", "b"},
		complementResult: []string{"d"},
		clearResult:      []string{},
		emptyResult:      true,
	}
	test2 := test{
		got:              []string{"a", "b", "c"},
		addResult:        []string{"a", "b", "c", "d", "e"},
		removeResult:     []string{"a", "b", "c"},
		hasResult:        false,
		countResult:      3,
		unionResult:      []string{"a", "b", "c", "d", "e"},
		minusResult:      []string{"a"},
		intersectResult:  []string{"a", "b", "c"},
		complementResult: []string{},
		clearResult:      []string{},
		emptyResult:      true,
	}
	// test1
	got1 := NewStringSet(test1.got...)
	got1.add("d")
	if !slice.CompareStringSlice(test1.addResult, got1.list()) {
		t.Errorf("string set add errorf: want %v , got %v", test1.addResult, got1.list())
	}
	got1.remove("d")
	if !slice.CompareStringSlice(test1.removeResult, got1.list()) {
		t.Errorf("string set remove errorf: want %v , got %v", test1.removeResult, got1.list())
	}
	if test1.hasResult != got1.has("c") {
		t.Errorf("string set has errorf: want %v , got %v", test1.hasResult, got1.has("c"))
	}
	if test1.countResult != got1.count() {
		t.Errorf("string set count errorf: want %v , got %v", test1.countResult, got1.count())
	}
	result := got1.union(NewStringSet("d")).list()
	if !slice.CompareStringSlice(test1.unionResult, result) {
		t.Errorf("string set union errorf: want %v , got %v", test1.unionResult, result)
	}
	result = got1.minus(NewStringSet("a", "b")).list()
	if !slice.CompareStringSlice(test1.minusResult, result) {
		t.Errorf("string set minus errorf: want %v , got %v", test1.minusResult, result)
	}
	result = got1.intersect(NewStringSet("a", "b", "d")).list()
	if !slice.CompareStringSlice(test1.intersectResult, result) {
		t.Errorf("string set intersect errorf: want %v , got %v", test1.intersectResult, result)
	}
	result = got1.complement(NewStringSet("c", "d")).list()
	if !slice.CompareStringSlice(test1.complementResult, result) {
		t.Errorf("string set complement errorf: want %v , got %v", test1.complementResult, result)
	}
	got1.clear()
	if !slice.CompareStringSlice(test1.clearResult, got1.list()) {
		t.Errorf("string set intersect errorf: want %v , got %v", test1.clearResult, got1.list())
	}
	if test1.emptyResult != got1.empty() {
		t.Errorf("string set complement errorf: want %v , got %v", test1.emptyResult, got1.empty())
	}
	// test2
	got2 := NewStringSet(test2.got...)
	got2.add("d", "e")
	if !slice.CompareStringSlice(test2.addResult, got2.list()) {
		t.Errorf("string set add errorf: want %v , got %v", test2.addResult, got1.list())
	}
	got2.remove("d", "e")
	if !slice.CompareStringSlice(test2.removeResult, got2.list()) {
		t.Errorf("string set remove errorf: want %v , got %v", test2.removeResult, got1.list())
	}
	if test2.hasResult != got2.has("d") {
		t.Errorf("string set has errorf: want %v , got %v", test2.hasResult, got2.list())
	}
	if test2.countResult != got2.count() {
		t.Errorf("string set count errorf: want %v , got %v", test2.countResult, got1.count())
	}
	result = got2.union(NewStringSet("d", "e")).list()
	if !slice.CompareStringSlice(test2.unionResult, result) {
		t.Errorf("string set union errorf: want %v , got %v", test2.unionResult, result)
	}
	result = got2.minus(NewStringSet("b", "c")).list()
	if !slice.CompareStringSlice(test2.minusResult, result) {
		t.Errorf("string set minus errorf: want %v , got %v", test2.minusResult, result)
	}
	result = got2.intersect(got2).list()
	if !slice.CompareStringSlice(test2.intersectResult, result) {
		t.Errorf("string set intersect errorf: want %v , got %v", test2.intersectResult, result)
	}
	result = got2.complement(NewStringSet("a", "b")).list()
	if !slice.CompareStringSlice(test2.complementResult, result) {
		t.Errorf("string set complement errorf: want %v , got %v", test2.complementResult, result)
	}
	got2.clear()
	if !slice.CompareStringSlice(test2.clearResult, got2.list()) {
		t.Errorf("string set clear errorf: want %v , got %v", test2.clearResult, got2.list())
	}
	if test2.emptyResult != got2.empty() {
		t.Errorf("string set empty errorf: want %v , got %v", test2.emptyResult, got2.empty())
	}
}

func TestIntSet(t *testing.T) {
	type test struct {
		got              []int
		addResult        []int
		removeResult     []int
		hasResult        bool
		countResult      int
		unionResult      []int
		minusResult      []int
		intersectResult  []int
		complementResult []int
		clearResult      []int
		emptyResult      bool
	}
	test1 := test{
		got:              []int{1, 2, 3},
		addResult:        []int{1, 2, 3, 4},
		removeResult:     []int{1, 2, 3},
		hasResult:        true,
		countResult:      3,
		unionResult:      []int{1, 2, 3, 4},
		minusResult:      []int{3},
		intersectResult:  []int{1, 2},
		complementResult: []int{4},
		clearResult:      []int{},
		emptyResult:      true,
	}
	test2 := test{
		got:              []int{1, 2, 3},
		addResult:        []int{1, 2, 3, 4, 5},
		removeResult:     []int{1, 2, 3},
		hasResult:        false,
		countResult:      3,
		unionResult:      []int{1, 2, 3, 4, 5},
		minusResult:      []int{1},
		intersectResult:  []int{1, 2, 3},
		complementResult: []int{},
		clearResult:      []int{},
		emptyResult:      true,
	}
	// test1
	got1 := NewIntSet(test1.got...)
	got1.add(4)
	if !slice.CompareIntSlice(test1.addResult, got1.list()) {
		t.Errorf("int set add errorf: want %v , got %v", test1.addResult, got1.list())
	}
	got1.remove(4)
	if !slice.CompareIntSlice(test1.removeResult, got1.list()) {
		t.Errorf("int set remove errorf: want %v , got %v", test1.removeResult, got1.list())
	}
	if test1.hasResult != got1.has(3) {
		t.Errorf("int set has errorf: want %v , got %v", test1.hasResult, got1.has(3))
	}
	if test1.countResult != got1.count() {
		t.Errorf("int set count errorf: want %v , got %v", test1.countResult, got1.count())
	}
	result := got1.union(NewIntSet(4)).list()
	if !slice.CompareIntSlice(test1.unionResult, result) {
		t.Errorf("int set union errorf: want %v , got %v", test1.unionResult, result)
	}
	result = got1.minus(NewIntSet(1, 2)).list()
	if !slice.CompareIntSlice(test1.minusResult, result) {
		t.Errorf("int set minus errorf: want %v , got %v", test1.minusResult, result)
	}
	result = got1.intersect(NewIntSet(1, 2, 4)).list()
	if !slice.CompareIntSlice(test1.intersectResult, result) {
		t.Errorf("int set intersect errorf: want %v , got %v", test1.intersectResult, result)
	}
	result = got1.complement(NewIntSet(3, 4)).list()
	if !slice.CompareIntSlice(test1.complementResult, result) {
		t.Errorf("int set complement errorf: want %v , got %v", test1.complementResult, result)
	}
	got1.clear()
	if !slice.CompareIntSlice(test1.clearResult, got1.list()) {
		t.Errorf("int set intersect errorf: want %v , got %v", test1.clearResult, got1.list())
	}
	if test1.emptyResult != got1.empty() {
		t.Errorf("int set complement errorf: want %v , got %v", test1.emptyResult, got1.empty())
	}
	// test2
	got2 := NewIntSet(test2.got...)
	got2.add(4, 5)
	if !slice.CompareIntSlice(test2.addResult, got2.list()) {
		t.Errorf("int set add errorf: want %v , got %v", test2.addResult, got1.list())
	}
	got2.remove(4, 5)
	if !slice.CompareIntSlice(test2.removeResult, got2.list()) {
		t.Errorf("int set remove errorf: want %v , got %v", test2.removeResult, got1.list())
	}
	if test2.hasResult != got2.has(4) {
		t.Errorf("int set has errorf: want %v , got %v", test2.hasResult, got2.list())
	}
	if test2.countResult != got2.count() {
		t.Errorf("int set count errorf: want %v , got %v", test2.countResult, got1.count())
	}
	result = got2.union(NewIntSet(4, 5)).list()
	if !slice.CompareIntSlice(test2.unionResult, result) {
		t.Errorf("int set union errorf: want %v , got %v", test2.unionResult, result)
	}
	result = got2.minus(NewIntSet(2, 3)).list()
	if !slice.CompareIntSlice(test2.minusResult, result) {
		t.Errorf("int set minus errorf: want %v , got %v", test2.minusResult, result)
	}
	result = got2.intersect(got2).list()
	if !slice.CompareIntSlice(test2.intersectResult, result) {
		t.Errorf("int set intersect errorf: want %v , got %v", test2.intersectResult, result)
	}
	result = got2.complement(NewIntSet(1, 2)).list()
	if !slice.CompareIntSlice(test2.complementResult, result) {
		t.Errorf("int set complement errorf: want %v , got %v", test2.complementResult, result)
	}
	got2.clear()
	if !slice.CompareIntSlice(test2.clearResult, got2.list()) {
		t.Errorf("int set clear errorf: want %v , got %v", test2.clearResult, got2.list())
	}
	if test2.emptyResult != got2.empty() {
		t.Errorf("int set empty errorf: want %v , got %v", test2.emptyResult, got2.empty())
	}
}
