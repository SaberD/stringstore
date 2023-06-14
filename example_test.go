package stringstore

import "fmt"

func ExampleStore() {
	// Create a new string store instance
	store, err := New("example.store")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Store a string
	err = store.Add("Store this message for later")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Retrieve the last string from the file
	line, err := store.Pop()
	if err != nil {
		fmt.Println(err)
		return
	}

	// check returned value is not empty (store was empty)
	if line != "" {
		fmt.Println("Retrieved value:", line)
	}

	// Output: Retrieved value: Store this message for later
}
