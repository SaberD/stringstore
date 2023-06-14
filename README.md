![Getting Started](logo.png)

# String store

Store string arrays persistently on disk with minimal CPU ram and disk usage.

The optimization methods used in this module restricts usage to only LIFO (last in, first out) operation, using only append and truncate operations on the storage file.

## Why not just use a map or a slice

In most cases you should. In-memory value stores are fast and efficient. The only issue is that they are lost when the program stops or the hardware power cycles. Writing to disk avoid data loss in the event of software or hardware resets by persisting the data to disk, and is intended to be used sparingly for critical data/messages that cannot be lost.

## Why not just read from file and store the array manually

Reading and writing to file directly wastes CPU, ram and disk resources.

1. First you have to read the whole file from disk and load it into memory
2. Then you have to add/remove the string from the loaded array
3. Then clear the old contents of the file
4. And finally rewrite the whole string array down to disk again

## Why not use SQLite

SQLite is a whole SQL database and can be used for a lot more than this lightweight module. But if all you need is a persistent LIFO string array then this module is more efficient and easy to understand and use. This module only has 3 functions: New(), Add() and Pop(). This module in combination with a simple JSON file database like [golang-scribble](https://github.com/nanobox-io/golang-scribble) is often good enough for most projects, and maintains the human readability and ease of data recovery. I have had cases where in an event of power loss, the files where corrupted, and then I could just open the files and fix the corrupted JSON files and string store lines manually.

## Features of this module:

- **Optimized String Retrieval:** By reading strings from the end of the file we eliminate the need to load and traverse through the entire file. A 1 byte buffer is all that is used while looking for the beginning of the last line. Each byte is also saved in a text array and then the resulting string is returned. The lines start position is then used to truncate the bytes of the last string from the file, removing it without having to read or open the file again.
- **Optimized String Storage:** By using append to add a string to the end of the file, we eliminate the need to load and rewrite the entire file for every insertion.
- **Human Readable:** Maintains human readability of data by storing it as strings in a text file rather than encoded bytes. This makes the stored data easier to reason about and repair, if necessary.

## Won't constantly reading and writing to the same area in the same file damage the disk drive?

SSD drives will automatically handle this by monitoring the reads and writes, and moving the file to another location on the disk in the background when needed.

The magnetic disks of HDD drives do not get damaged by reading/writing to the same area, because the mechanical parts of the HDD, like the motor or the read/write head, will wear out long before the magnetic disks. We are however using those mechanical parts more compared to an in-memory approach, but this is kept to a minimum using this modules optimizations.

## Getting Started

To include this module in your Go project, use the following command:

```bash
go get github.com/saberd/stringstore
```

Import it into your Go code using:

```go
import "github.com/saberd/stringstore"
```

## Usage

Here's an example of how to use this module:

```go
package main

import (
	"github.com/saberd/stringstore"
)

func main() {
	// Create a new string store instance
	store := stringstore.New("/path/to/your/file")

  // Store a string
	err := store.Add("Store this message for later")
	if err != nil {
		println(err)
		return
	}

	// Retrieve the last string from the file
	line, err := store.Pop()
	if err != nil {
		println(err)
		return
	}

	println("Retrieved value: ", line)
}
```

## Testing

This module includes a set of unit tests. To run these tests, navigate to the project's root directory and execute:

```bash
go test
```

## Contribution

Contributions are always welcome. Please fork the repository and create a pull request with your changes.

## License

This module is licensed under the MIT License. See the LICENSE file for more details.

## Support

For issues, feature requests, or general inquiries, please open a GitHub issue.