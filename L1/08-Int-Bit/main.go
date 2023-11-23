package main

import "fmt"

func main() {

	// Variables to store input
	var (
		in       int64
		shiftPos uint
		shiftVal uint
	)

	// Get the initial number
	fmt.Println("Enter initial int64 number:")
	if _, err := fmt.Scanln(&in); err != nil {
		fmt.Println(err.Error())
		return
	}

	// Get the position in which the shift must happen
	fmt.Println("Enter shift position (0-63):")
	if _, err := fmt.Scanln(&shiftPos); err != nil {
		fmt.Println(err.Error())
		return
	}
	if shiftPos > 63 || shiftPos < 0 {
		fmt.Println("Wrong position value")
		return
	}

	// Get the value to which a bit needs to be onverted
	fmt.Println("Enter shift value (0|1):")
	if _, err := fmt.Scanln(&shiftVal); err != nil {
		fmt.Println(err.Error())
		return
	}
	if shiftVal > 1 || shiftVal < 0 {
		fmt.Println("Wrong shift value")
		return
	}

	fmt.Printf("Initial value: %b\n", in)

	// Create a shift for changing the bit
	shift := int64(1) << shiftPos

	// Apply the reversed shift to the initial number
	// By and operation which neutralises bit value one
	//
	// Example:
	// 00101011 &
	// 11101111 =
	// 00100011
	out := in & ^shift
	// Set the user value into the neutralised bit's place
	// (if user set 0 - nothing changes)
	//
	// Example:
	// 00100011 |
	// 00010000 =
	// 00110011
	out |= int64(shiftVal) << shiftPos
	fmt.Printf("%d with %d at position %d looks like this: \n%b (%d)\n", in, shiftVal, shiftPos, out, out)
}
