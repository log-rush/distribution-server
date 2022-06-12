package repository

import "github.com/google/uuid"

func GenerateID() string {
	return uuid.NewString()
}

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
	index := -1
	for i, item := range *slc {
		if item == v {
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
