package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	// Create a temporary file
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	content := []byte("video/riverStyx/1011_19912_1233.mp4\n" +
		"video/riverStyx/1011_19912_1235.mp4\n" +
		"video/riverStyx/1011_19912_1212.mp4\n")
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Call the function with the temp file
	store, err := New(filePath)
	if err != nil {
		t.Fatal(err)
	}

	// Check the result
	if length(store) != 3 {
		t.Error("Length is not 3", length(store))
	}
}

// TestPop tests the Pop method on a store until its empty
func TestPop(t *testing.T) {
	// Create a temporary file
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	content := []byte("video/riverStyx/1011_19912_1233.mp4\n" +
		"video/riverStyx/1011_19912_1235.mp4\n" +
		"video/riverStyx/1011_19912_1212.mp4\n")
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		t.Fatal(err)
	}

	store, err := New(filePath)
	if err != nil {
		t.Fatal(err)
	}

	// Test the Pop method
	for length(store) > 0 {
		// read the last line from file
		tmp, err := readFile(store.path)
		if err != nil {
			t.Fatal(err)
		}
		last := tmp[len(tmp)-1]
		// test Pop method
		line, err := store.Pop()
		if err != nil {
			t.Fatal(err)
		}
		if line != last {
			t.Errorf("pop returned wrong value wanted %s, got %s", last, line)
		}
		// read file and check that deletion worked
		tmp, err = readFile(filePath)
		if err != nil {
			t.Fatal(err)
		}
		// check slice and file match
		fileSlice, err := readFile(store.path)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Join(tmp, ",") != strings.Join(fileSlice, ",") {
			t.Error(tmp)
		}
	}
	// check that the file is empty
	tmp, err := readFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if len(tmp) != 0 {
		t.Error("wanted empty store, got", tmp)
	}
}

// TestPopEmpty tests what happens when Pop is called on an empty store
func TestPopEmpty(t *testing.T) {
	// Create a temporary file
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	content := []byte("")
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		t.Fatal(err.Error())
	}

	store, err := New(filePath)
	if err != nil {
		t.Fatal(err.Error())
	}

	// test Pop method
	value, err := store.Pop()
	if err != nil {
		t.Fatal(err.Error())
	}
	if value != "" {
		t.Error("value is not empty", value)
	}
}

// TestAdd tests the Add method when store contains values
func TestAdd(t *testing.T) {
	// Create a temporary file
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	content := []byte("video/riverStyx/1011_19912_1233.mp4\n" +
		"video/riverStyx/1011_19912_1235.mp4\n" +
		"video/riverStyx/1011_19912_1212.mp4\n")
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		t.Fatal(err.Error())
	}

	store, err := New(filePath)
	if err != nil {
		t.Fatal(err.Error())
	}
	file := "video/riverStyx/1011_19912_3235.mp4"
	err = store.Add(file)
	if err != nil {
		t.Fatal(err.Error())
	}
	want := []string{"video/riverStyx/1011_19912_1233.mp4", "video/riverStyx/1011_19912_1235.mp4", "video/riverStyx/1011_19912_1212.mp4", "video/riverStyx/1011_19912_3235.mp4"}
	got, err := readFile(filePath)
	if err != nil {
		t.Fatal(err.Error())
	}
	if strings.Join(want, ",") != strings.Join(got, ",") {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

// TestAddEmpty tests the Add method when the store is empty
func TestAddEmpty(t *testing.T) {
	// Create a temporary file
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	content := []byte("")
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		t.Fatal(err.Error())
	}

	store, err := New(filePath)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = store.Add("video/riverStyx/1011_19912_3235.mp4")
	if err != nil {
		t.Fatal(err.Error())
	}

	//remove values
	_, err = store.Pop()
	if err != nil {
		t.Fatal(err.Error())
	}
	err = store.Add("video/riverStyx/1011_19912_3235.mp4")
	if err != nil {
		t.Fatal(err.Error())
	}
	want := []string{"video/riverStyx/1011_19912_3235.mp4"}
	got, err := readFile(filePath)
	if err != nil {
		t.Fatal(err.Error())
	}
	if strings.Join(want, ",") != strings.Join(got, ",") {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestContains(t *testing.T) {
	// Create a temporary file
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	content := []byte("")
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		t.Fatal(err.Error())
	}

	store, err := New(filePath)
	if err != nil {
		t.Fatal(err.Error())
	}

	data := []string{"video/riverStyx/1011_19912_1233.mp4\n",
		"video/riverStyx/1011_19912_1235.mp4\n",
		"video/riverStyx/1011_19912_1212.mp4\n"}

	for _, file := range data {
		err = store.Add(file)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
	_, err = readFile(store.path)
	if err != nil {
		t.Error(err)
	}

	for _, line := range data {
		if !contains(store, line) {
			t.Error("failed to find", line)
		}
	}

	for line, err := store.Pop(); line != ""; line, err = store.Pop() {
		if err != nil {
			t.Fatal(err.Error())
		}
		if contains(store, line) {
			t.Error("failed to pop", line)
		}
	}
}

func readFile(path string) (out []string, err error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return out, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		out = append(out, scanner.Text())
	}
	err = scanner.Err()
	if err != nil {
		return out, err
	}
	return out, nil
}

func length(s *Store) int {
	tmp, _ := readFile(s.path)
	return len(tmp)
}

func contains(s *Store, line string) bool {
	tmp, _ := readFile(s.path)
	for _, val := range tmp {
		// Prints out unicode representation (shows hidden characters)
		// fmt.Printf("val, %+q\n", val)
		// fmt.Printf("line, %+q\n", line)
		if strings.TrimSuffix(val, "\n") == strings.TrimSuffix(line, "\n") {
			return true
		}
	}
	return false
}
