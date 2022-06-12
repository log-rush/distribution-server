package repository

func AppendUniqueToSlice[T any](slc *[]T, v T, comp func(v T) bool) *[]T {
	for _, item := range *slc {
		if comp(item) {
			return slc
		}
	}
	newSlice := append(*slc, v)
	return &newSlice
}

func RemoveFromSlice[T any](slc *[]T, comp func(v T) bool) *[]T {
	index := -1
	for i, item := range *slc {
		if comp(item) {
			index = i
			break
		}
	}
	if index == -1 {
		return slc
	}
	newSlice := append((*slc)[:index], (*slc)[index+1:]...)
	return &newSlice
}
