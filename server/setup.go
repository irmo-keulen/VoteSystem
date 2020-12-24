package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Handles setting up the environment
// Does the following things:
// Checks if private key file exists
// Creates a new key pair if key doesn't exists
// Writes the keyBytes to key files
func setup() error {
	setupVote()
	fpPub, _ := filepath.Abs(filenamePub)
	fpPriv, _ := filepath.Abs(filenamePriv)
	fmt.Printf("Generating keys: \n- Privpath: \t%s\n- Pubpath:\t%s\n", fpPriv, fpPub)

	if _, err := os.Stat(fpPriv); !os.IsNotExist(err) { // Checks if private key already exists
		fmt.Printf("%v PrivKey Path Exists.\nskipping generating keys.\n", ck)
		return nil
	}
	privKey, pubKey := GenerateKeyPair(4096)
	privFile, err := os.OpenFile(filenamePriv, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println("Error generating privKey file")
		log.Fatal(err)
	}
	_, err = privFile.Write(PrivateKeyToBytes(privKey))
	if err != nil {
		fmt.Println("Error writing to privFile")
		log.Fatal(err)
	}
	defer privFile.Close()

	pubFile, err := os.OpenFile(filenamePub, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println("Error generating pubKey file")
		log.Fatal(err)
	}
	_, err = pubFile.Write(PublicKeyToBytes(pubKey))
	if err != nil {
		fmt.Println("Error writing to pubFile")
		log.Fatal(err)
	}
	return testDBConnect()
}

// Prompts user for the vote subject
// Right now there is only a binary vote system
// TODO multi answer voting <- Not sure if this is necessary for the assignment, might implement later
func setupVote() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What is the subject of the vote?")
	voteSubject, _ = reader.ReadString('\n')
	fmt.Printf("The vote will be about %s\n", voteSubject)
	return
}

// A VERY basic test (writes to database and checks error, and retreives the data)
// TODO Randomize data written
func testDBConnect() error {
	fmt.Println("Testing Connection")

	err := rdb.Set(ctx, "Hello", "World", 0).Err()
	if err != nil {
		return fmt.Errorf("redis error during write: %s", err.Error())
	}

	val, err := rdb.Get(ctx, "Hello").Result()
	if err != nil {
		return fmt.Errorf("redis error durign GET: %s", err.Error())
	}

	if val != "World" {
		return fmt.Errorf("redis error, result not equal")
	}
	fmt.Printf("%s Connection Succes\n", ck)
	return nil
}
