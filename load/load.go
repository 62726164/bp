package main

import (
	"bufio"
	"fmt"
	"github.com/willf/bloom"
	"log"
	"math"
	"os"
)

var fp = .001                                                                      // False Positive Rate (1 in 1,000)
var n = 517238891.0                                                                // Hash list cardinality
var m = math.Ceil((n * math.Log(fp)) / math.Log(1.0/math.Pow(2.0, math.Log(2.0)))) // Number of bits in the filter
var k = uint(10)                                                                   // Number of hash functions
var filter = bloom.New(uint(m), k)

// load creates a bloom filter from the partial hashes
// and saves the filter to a file. The hashes must be UPPERCASE or the checks will fail.
func main() {
	usage := "load /path/to/partial/hashes.txt /path/to/output.filter"

	if len(os.Args) != 3 {
		fmt.Println(usage)
		return
	}

	hashFile := os.Args[1]
	filterFile := os.Args[2]

	// Populate the bloom filter
	file, err := os.Open(hashFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		filter.Add(scanner.Bytes())
	}

	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	// Save the bloom filter to a file
	f, err := os.Create(filterFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	bytesWritten, err := filter.WriteTo(f)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bytes written to filter: %d\n", bytesWritten)
}
