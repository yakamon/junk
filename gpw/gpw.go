package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
)

const (
	// MinASCII is min int to convert to character in ASCII
	MinASCII = 33
	// MaxASCII is max int to convert to character in ASCII
	MaxASCII = 126
)

// ParseArgs parses command line arguments
func ParseArgs() (length int, useSpecials bool, err error) {
	length, useSpecials = 32, true
	if len(os.Args) < 3 {
		fmt.Printf("Missing arguments, use default values: length = %v, useSpecials = %v\n", length, useSpecials)
		return
	}

	length, err = strconv.Atoi(os.Args[1])
	if err != nil {
		return
	}

	useSpecialsInt, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return
	}
	if useSpecialsInt > 0 {
		useSpecials = true
	} else {
		useSpecials = false
	}

	return
}

// GeneratePassword generates random password
func GeneratePassword(length int, useSpecials bool) (string, error) {
	b := make([]byte, length)
	for i := range b {
		// Generate random ASCII int code
		r, err := rand.Int(rand.Reader, big.NewInt(MaxASCII-MinASCII))
		if err != nil {
			return "", err
		}
		n, _ := strconv.Atoi(r.String())
		n += MinASCII

		// Set to byte slice
		b[i] = byte(n)
	}

	return string(b), nil
}

func main() {
	length, useSpecials, err := ParseArgs()
	if err != nil {
		panic(err)
	}

	pw, err := GeneratePassword(length, useSpecials)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Generated Password: %v\n", pw)
}
