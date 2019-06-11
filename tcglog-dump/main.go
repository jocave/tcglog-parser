package main

import (
	"flag"
	"fmt"
	"github.com/chrisccoulson/tcglog-parser"
	"io"
	"os"
)

var (
	alg string
)

func init() {
	flag.StringVar(&alg, "alg", "sha1", "Name of the hash algorithm to display")
}

func main() {
	flag.Parse()

	var algorithmId tcglog.AlgorithmId
	switch alg {
	case "sha1":
		algorithmId = tcglog.AlgorithmSha1
	case "sha256":
		algorithmId = tcglog.AlgorithmSha256
	case "sha384":
		algorithmId = tcglog.AlgorithmSha384
	case "sha512":
		algorithmId = tcglog.AlgorithmSha512
	default:
		fmt.Fprintf(os.Stderr, "Unrecognized algorithm\n")
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "Too many arguments\n")
		os.Exit(1)
	}

	var path string
	if len(args) == 1 {
		path = args[0]
	} else {
		path = "/sys/kernel/security/tpm0/binary_bios_measurements"
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
		os.Exit(1)
	}

	log, err := tcglog.NewLogFromFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse log file: %v\n", err)
		os.Exit(1)
	}

	if !log.HasAlgorithm(algorithmId) {
		fmt.Fprintf(os.Stderr, "The log doesn't contain entries for the specified digest algorithm\n")
		os.Exit(1)
	}

	for {
		event, err := log.NextEvent()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%2d %s %s\n", event.PCRIndex, event.Digests[algorithmId], event.EventType)
	}
}
