package main

import (
	"github.com/ProtonMail/ui"
)

func subSlice(slice interface{}, index int) interface{} {
	if grids, ok := slice.([]*ui.Grid); ok {
		slice = append(grids[:index], grids[index+1:]...)
	} else if entries, ok := slice.([]*SQLEntry); ok {
		var temp []*SQLEntry
		temp = append(temp, entries[:index]...)
		temp = append(temp, entries[index+1:]...)
		slice = temp
	}
	return slice
}
