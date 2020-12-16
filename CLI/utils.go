package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Reads the key file and returns public key
// Returns as byte slice
func getPubKey(filename string) ([]byte, error) {
	pubFile, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to read file. Err: %s", err.Error())
	}
	key, err := ioutil.ReadAll(pubFile)
	if err != nil {
		return nil, fmt.Errorf("error reading file Err: %s", err.Error())
	}
	return key, nil
}
