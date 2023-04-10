package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"os"
)

func main() {
	decrypt()
}

func decrypt() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <input file> <output file> <password>\n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	password := os.Args[3]

	// Open input file
	input, err := os.Open(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer input.Close()

	// Create output file
	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer output.Close()

	// Read initialization vector from input file
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(input, iv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Generate 32-byte key from password
	key := make([]byte, 32)
	for i := 0; i < len(password); i++ {
		key[i%32] ^= password[i]
	}

	// Create AES cipher with key and initialization vector
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	stream := cipher.NewCTR(block, iv)

	// Create buffered reader and writer for input and output files
	inputReader := bufio.NewReader(input)
	outputWriter := bufio.NewWriter(output)

	// Decrypt input file and write to output file
	decryptedReader := &cipher.StreamReader{S: stream, R: inputReader}
	_, err = io.Copy(outputWriter, decryptedReader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Flush output buffer to ensure all data is written to file
	err = outputWriter.Flush()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Decryption Complate Message
	fmt.Printf("Decryption complete: %s -> %s\n", inputFile, outputFile)
}
