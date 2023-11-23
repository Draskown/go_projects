package main

import "fmt"

// Interface to hold actions that
// Human struct can perform
type HumanAction interface {
	Birthday()
}

// Parent struct Human that implements
// HumanAction interface
type Human struct {
	Name    string
	FName   string
	Age     int
	working bool
	HumanAction
}

// Interface implementation
func (h *Human) Birthday() {
	h.Age++
}

// Simple composition 'inheritance'
//
// Child struct ActionFields inherit parent's
// fileds and methods.
//
// So interface HumanAction implementation can
// be called directly from the ActionFields variable
type ActionFields struct {
	*Human
}

// Additional method to alter the parent's struct field
//
// Which demonstrates how composition can be difficult
// to maintain due to the passing of fields along with methods
// of the parental struct
func (a *ActionFields) ChangeName(name string) {
	a.Human.Name = name
}

// Child struct that inherits only parent's methods
// by composing HumanAction interface.
//
// Because this struct does not implement HumanAction,
// methods are taken from the passed parental struct that does
type ActionMethods struct {
	HumanAction
	Age int
}

// The approach of inheriting only methods
// allows us to overwrite them for new logic
// by still having access to other parent's struct methods
func (a *ActionMethods) Birthday() {
	if a.Age > 0 {
		a.Age++
	} else {
		// Call parental method if Age wasn't set
		a.HumanAction.Birthday()
	}
}

func main() {
	// Create parental struct
	human := Human{
		Name:  "Andrei",
		FName: "Kabalin",
		Age:   22,
	}
	fmt.Printf("%+v\n", human)

	// Create child struct with only methods
	actionMethods := ActionMethods{HumanAction: &human}

	// Call HumanAction implementation
	// which is actually performed from
	// the parental struct
	actionMethods.Birthday()
	fmt.Printf("%+v\n", human)

	// Create child struct with methods AND fields
	actionFields := ActionFields{&human}

	// Call methods that change the fields
	// of the parental struct
	actionFields.Birthday()
	actionFields.ChangeName("Vladimir")
	fmt.Printf("%+v\n", human)
}
