package swort

import (
	"cmp"
	"slices"
)

type Slice[IT any, ST cmp.Ordered] struct {
	sortValues            []ST
	sortValue2Item        map[ST][]IT
	items                 []IT
	originalSortValue2Idx map[ST]int
	ascSortedItems        []IT
	descSortValue2Idx     map[ST]int
	descSortedItems       []IT
	ascSortValue2Idx      map[ST]int
	sortValueExtractor    func(item IT) ST
}

func MakeSlice[IT any, ST cmp.Ordered](items []IT, sortValueExtractor func(item IT) ST) *Slice[IT, ST] {
	return &Slice[IT, ST]{
		items:              items,
		sortValueExtractor: sortValueExtractor,
	}
}

func (s *Slice[IT, ST]) Len() int {
	return len(s.items)
}

func (s *Slice[IT, ST]) SearchFromOriginal(item IT) (int, bool) {
	s.initialize()

	idx, ok := s.originalSortValue2Idx[s.sortValueExtractor(item)]
	if !ok {
		idx = s.Len()
	}
	return idx, ok

}

func (s *Slice[IT, ST]) SearchFromSortedByAsc(item IT) (int, bool) {
	if s.ascSortValue2Idx == nil {
		s.sortByAsc()
	}

	idx, ok := s.ascSortValue2Idx[s.sortValueExtractor(item)]
	if !ok {
		idx = s.Len()
	}
	return idx, ok
}

func (s *Slice[IT, ST]) SearchFromSortedByDesc(item IT) (int, bool) {
	if s.descSortValue2Idx == nil {
		s.sortByDesc()
	}

	idx, ok := s.descSortValue2Idx[s.sortValueExtractor(item)]
	if !ok {
		idx = s.Len()
	}
	return idx, ok
}

func (s *Slice[IT, ST]) SortByAsc() []IT {
	return s.sortByAsc()
}

func (s *Slice[IT, ST]) SortByDesc() []IT {
	return s.sortByDesc()
}

func (s *Slice[IT, ST]) sortByAsc() []IT {
	if ascSorted := s.ascSortedItems; ascSorted != nil {
		return ascSorted
	}

	s.initialize()

	slices.Sort(s.sortValues)

	s.ascSortValue2Idx = make(map[ST]int)
	s.ascSortedItems = make([]IT, 0, len(s.items))
	idx := 0
	for _, sortValue := range s.sortValues {
		for _, item := range s.sortValue2Item[sortValue] {
			s.ascSortedItems = append(s.ascSortedItems, item)

			if _, ok := s.ascSortValue2Idx[sortValue]; !ok {
				s.ascSortValue2Idx[sortValue] = idx
			}
			idx++
		}
	}

	return s.ascSortedItems
}

func (s *Slice[IT, ST]) sortByDesc() []IT {
	if descSorted := s.descSortedItems; descSorted != nil {
		return descSorted
	}

	ascSorted := s.sortByAsc()
	l := len(ascSorted)

	s.descSortValue2Idx = make(map[ST]int)
	s.descSortedItems = make([]IT, l)
	idx := 0
	for i := l - 1; i >= 0; i-- {
		s.descSortedItems[idx] = ascSorted[i]

		sortValue := s.sortValueExtractor(ascSorted[i])
		if _, ok := s.descSortValue2Idx[sortValue]; !ok {
			s.descSortValue2Idx[sortValue] = idx
		}

		idx++
	}

	return s.descSortedItems
}

func (s *Slice[IT, ST]) initialize() {
	if s.originalSortValue2Idx != nil {
		return
	}

	sortValues := make([]ST, 0, len(s.items))
	sortValue2Item := map[ST][]IT{}
	originalSortValue2Idx := map[ST]int{}

	for i, item := range s.items {
		sortVal := s.sortValueExtractor(item)

		gotItems, ok := sortValue2Item[sortVal]
		if !ok {
			sortValues = append(sortValues, sortVal)
			sortValue2Item[sortVal] = []IT{item}
			originalSortValue2Idx[sortVal] = i
			continue
		}

		sortValue2Item[sortVal] = append(gotItems, item)
	}

	s.sortValues = sortValues
	s.sortValue2Item = sortValue2Item
	s.originalSortValue2Idx = originalSortValue2Idx
}
