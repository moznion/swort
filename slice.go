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
	sortValueExtractor    func(item IT) ST
	sortedByAsc           *Sorted[IT, ST]
	sortedByDesc          *Sorted[IT, ST]
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
	s.initializeIfNeeded()

	idx, ok := s.originalSortValue2Idx[s.sortValueExtractor(item)]
	if !ok {
		idx = s.Len()
	}
	return idx, ok

}

func (s *Slice[IT, ST]) SortByAsc() *Sorted[IT, ST] {
	return s.sortByAsc()
}

func (s *Slice[IT, ST]) SortByDesc() *Sorted[IT, ST] {
	if sortedByDesc := s.sortedByDesc; sortedByDesc != nil {
		return sortedByDesc
	}

	sortedByAsc := s.sortByAsc()
	l := len(sortedByAsc.Items)

	descSortValue2Idx := make(map[ST]int)
	descSortedItems := make([]IT, l)
	idx := 0
	for i := l - 1; i >= 0; i-- {
		item := sortedByAsc.Items[i]

		descSortedItems[idx] = item
		sortValue := s.sortValueExtractor(item)
		if _, ok := descSortValue2Idx[sortValue]; !ok {
			descSortValue2Idx[sortValue] = idx
		}

		idx++
	}

	sortedByDesc := &Sorted[IT, ST]{
		Items:              descSortedItems,
		sortValue2Idx:      descSortValue2Idx,
		sortValueExtractor: s.sortValueExtractor,
	}
	s.sortedByDesc = sortedByDesc
	return sortedByDesc
}

func (s *Slice[IT, ST]) sortByAsc() *Sorted[IT, ST] {
	if sortedByAsc := s.sortedByAsc; sortedByAsc != nil {
		return sortedByAsc
	}

	s.initializeIfNeeded()

	slices.Sort(s.sortValues)

	ascSortValue2Idx := make(map[ST]int)
	ascSortedItems := make([]IT, 0, len(s.items))
	idx := 0
	for _, sortValue := range s.sortValues {
		for _, item := range s.sortValue2Item[sortValue] {
			ascSortedItems = append(ascSortedItems, item)

			if _, ok := ascSortValue2Idx[sortValue]; !ok {
				ascSortValue2Idx[sortValue] = idx
			}
			idx++
		}
	}

	sortedByAsc := &Sorted[IT, ST]{
		Items:              ascSortedItems,
		sortValue2Idx:      ascSortValue2Idx,
		sortValueExtractor: s.sortValueExtractor,
	}
	s.sortedByAsc = sortedByAsc
	return sortedByAsc
}

func (s *Slice[IT, ST]) initializeIfNeeded() {
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

type Sorted[IT any, ST cmp.Ordered] struct {
	Items              []IT
	sortValue2Idx      map[ST]int
	sortValueExtractor func(item IT) ST
}

func (sd *Sorted[IT, ST]) Search(item IT) (int, bool) {
	idx, ok := sd.sortValue2Idx[sd.sortValueExtractor(item)]
	if !ok {
		idx = len(sd.Items)
	}
	return idx, ok
}
