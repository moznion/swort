package swort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Struct struct {
	intValue int64
}

func TestSlice_Len(t *testing.T) {
	given := []Struct{{intValue: 5}, {intValue: 2}, {intValue: 6}, {intValue: 3}, {intValue: 1}, {intValue: 5} /* <= duplicated */, {intValue: 4}}
	slice := MakeSlice(given, func(s Struct) int64 {
		return s.intValue
	})

	assert.Equal(t, 7, slice.Len())
}

func TestSlice_SortByAsc(t *testing.T) {
	given := []Struct{{intValue: 5}, {intValue: 2}, {intValue: 6}, {intValue: 3}, {intValue: 1}, {intValue: 5} /* <= duplicated */, {intValue: 4}}
	slice := MakeSlice(given, func(s Struct) int64 {
		return s.intValue
	})

	assert.EqualValues(t, []Struct{{intValue: 1}, {intValue: 2}, {intValue: 3}, {intValue: 4}, {intValue: 5}, {intValue: 5}, {intValue: 6}}, slice.SortByAsc().Items, "should be sorted by ASC")
	assert.EqualValues(t, slice.SortByAsc(), slice.SortByAsc(), "should be sorted by ASC; check for retrieval from cache")
	assert.EqualValues(t, []Struct{{intValue: 5}, {intValue: 2}, {intValue: 6}, {intValue: 3}, {intValue: 1}, {intValue: 5}, {intValue: 4}}, given, "should given be immutable")
}

func TestSlice_SortByDesc(t *testing.T) {
	given := []Struct{{intValue: 5}, {intValue: 2}, {intValue: 6}, {intValue: 3}, {intValue: 1}, {intValue: 5} /* <= duplicated */, {intValue: 4}}
	slice := MakeSlice(given, func(s Struct) int64 {
		return s.intValue
	})

	assert.EqualValues(t, []Struct{{intValue: 6}, {intValue: 5}, {intValue: 5}, {intValue: 4}, {intValue: 3}, {intValue: 2}, {intValue: 1}}, slice.SortByDesc().Items, "should be sorted by Desc")
	assert.EqualValues(t, slice.SortByDesc(), slice.SortByDesc(), "should be sorted by Desc, check for retrieval from cache")
	assert.EqualValues(t, []Struct{{intValue: 5}, {intValue: 2}, {intValue: 6}, {intValue: 3}, {intValue: 1}, {intValue: 5}, {intValue: 4}}, given, "should given be immutable")
}

func TestSlice_SearchFromOriginal(t *testing.T) {
	given := []Struct{{intValue: 5}, {intValue: 2}, {intValue: 6}, {intValue: 3}, {intValue: 1}, {intValue: 5} /* <= duplicated */, {intValue: 4}}
	slice := MakeSlice(given, func(s Struct) int64 {
		return s.intValue
	})

	type testCase struct {
		intValue    int64
		expectedIdx int
		expectedOK  bool
	}
	testCases := []testCase{
		{
			intValue:    5,
			expectedIdx: 0,
			expectedOK:  true,
		},
		{
			intValue:    2,
			expectedIdx: 1,
			expectedOK:  true,
		},
		{
			intValue:    6,
			expectedIdx: 2,
			expectedOK:  true,
		},
		{
			intValue:    3,
			expectedIdx: 3,
			expectedOK:  true,
		},
		{
			intValue:    1,
			expectedIdx: 4,
			expectedOK:  true,
		},
		{
			intValue:    4,
			expectedIdx: 6,
			expectedOK:  true,
		},
		{
			intValue:    7,
			expectedIdx: 7,
			expectedOK:  false,
		},
	}

	for _, tc := range testCases {
		got, ok := slice.SearchFromOriginal(Struct{intValue: tc.intValue})
		assert.Equal(t, got, tc.expectedIdx)
		assert.Equal(t, ok, tc.expectedOK)
	}
}

func TestSlice_SearchFromSortedByAsc(t *testing.T) {
	given := []Struct{{intValue: 5}, {intValue: 2}, {intValue: 6}, {intValue: 3}, {intValue: 1}, {intValue: 5} /* <= duplicated */, {intValue: 4}}
	slice := MakeSlice(given, func(s Struct) int64 {
		return s.intValue
	})

	type testCase struct {
		intValue    int64
		expectedIdx int
		expectedOK  bool
	}
	testCases := []testCase{
		{
			intValue:    1,
			expectedIdx: 0,
			expectedOK:  true,
		},
		{
			intValue:    2,
			expectedIdx: 1,
			expectedOK:  true,
		},
		{
			intValue:    3,
			expectedIdx: 2,
			expectedOK:  true,
		},
		{
			intValue:    4,
			expectedIdx: 3,
			expectedOK:  true,
		},
		{
			intValue:    5,
			expectedIdx: 4,
			expectedOK:  true,
		},
		{
			intValue:    6,
			expectedIdx: 6,
			expectedOK:  true,
		},
		{
			intValue:    7,
			expectedIdx: 7,
			expectedOK:  false,
		},
	}

	sortedByAsc := slice.SortByAsc()
	for _, tc := range testCases {
		got, ok := sortedByAsc.Search(Struct{intValue: tc.intValue})
		assert.Equal(t, got, tc.expectedIdx)
		assert.Equal(t, ok, tc.expectedOK)
	}
}

func TestSlice_SearchFromSortedByDesc(t *testing.T) {
	given := []Struct{{intValue: 5}, {intValue: 2}, {intValue: 6}, {intValue: 3}, {intValue: 1}, {intValue: 5} /* <= duplicated */, {intValue: 4}}
	slice := MakeSlice(given, func(s Struct) int64 {
		return s.intValue
	})

	type testCase struct {
		intValue    int64
		expectedIdx int
		expectedOK  bool
	}
	testCases := []testCase{
		{
			intValue:    6,
			expectedIdx: 0,
			expectedOK:  true,
		},
		{
			intValue:    5,
			expectedIdx: 1,
			expectedOK:  true,
		},
		{
			intValue:    4,
			expectedIdx: 3,
			expectedOK:  true,
		},
		{
			intValue:    3,
			expectedIdx: 4,
			expectedOK:  true,
		},
		{
			intValue:    2,
			expectedIdx: 5,
			expectedOK:  true,
		},
		{
			intValue:    1,
			expectedIdx: 6,
			expectedOK:  true,
		},
		{
			intValue:    7,
			expectedIdx: 7,
			expectedOK:  false,
		},
	}

	sortedByDesc := slice.SortByDesc()
	for _, tc := range testCases {
		got, ok := sortedByDesc.Search(Struct{intValue: tc.intValue})
		assert.Equal(t, got, tc.expectedIdx)
		assert.Equal(t, ok, tc.expectedOK)
	}
}
