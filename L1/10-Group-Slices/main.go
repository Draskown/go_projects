package main

import (
	"fmt"
	"sort"
	"strconv"
)

// Struct to hold the groups of inital values
type group struct {
	// Name (-20, -10, 10, 20 etc)
	groupName string
	// Values inside of the group
	groupValues []float32
}

func main() {
	// Aray of initial values
	// Extended
	initArr := []float32{-25.4, -27.0, 13.0, 19.0,
		15.5, 24.5, -21.0, 32.5,
		0.0, 2.5, 90.1, 9.5,
		-12.6, -18.2, -16.9, -10.0,
		-9.5, -3.6, 0.0, -2.2,
	}

	// Slice for all of the groups
	groups := []group{}

	// Loop through every element
	for _, val := range initArr {
		// Calculate to which group does
		// a value belong
		groupInt := int(val/10) * 10

		// Convert group int to string
		// for -0
		var groupStr string
		if groupInt == 0 && val < 0 {
			// Negative values go under
			// -0 group
			groupStr = "-0"
		} else {
			groupStr = strconv.Itoa(groupInt)
		}

		// Get index of the existing group
		ix, ok := exists(groupStr, groups)
		// If a group with this name does not exist
		if !ok {
			// Append it to the slice
			groups = append(groups, group{
				groupName: groupStr,
			})
			// Get the index of the last added group
			ix = len(groups) - 1
		}
		// Append slice value to a group
		groups[ix].groupValues = append(groups[ix].groupValues, val)
	}

	// Sort ascending
	sortGroups(groups)

	fmt.Println(groups)
}

// Checks if the group with the given name alread exists
//
// Returns an index of the group and bool indicating that
// a group exists or not
func exists(name string, groups []group) (int, bool) {
	// Iterate through every group
	for i, g := range groups {
		// If name exists
		if g.groupName == name {
			return i, true
		}
	}

	// If not
	return 0, false
}

// Sorts groups by ascending order
func sortGroups(groups []group) {
	// Use in-built sort package
	// with a custom condition
	sort.Slice(groups, func(i, j int) bool {
		a, _ := strconv.Atoi(groups[i].groupName)
		b, _ := strconv.Atoi(groups[j].groupName)
		return a < b
	})

	// Swap -0 and 0 group if they are in the wrong order
	zeroIx, _ := exists("0", groups)
	negZeroIx, _ := exists("-0", groups)
	if zeroIx < negZeroIx {
		groups[zeroIx], groups[negZeroIx] = groups[negZeroIx], groups[zeroIx]
	}
}
