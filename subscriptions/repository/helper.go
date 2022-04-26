package repository

func AppendUniqueToSlice[T comparable](slc *[]T, v T) *[]T {
	for _, item := range *slc {
		if item == v {
			return slc
		}
	}
	newSlice := append(*slc, v)
	return &newSlice
}

func RemoveFromSlice[T comparable](slc *[]T, v T) *[]T {
	var index int
	for i, item := range *slc {
		if item == v {
			index = i
			break
		}
	}
	newSlice := append((*slc)[:index], (*slc)[index+1:]...)
	return &newSlice
}
