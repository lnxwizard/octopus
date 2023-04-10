package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

func main() {
	encrypt()
}

func encrypt() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <input_file> <output_file> <password> \n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	password := os.Args[3]

	// Open Input File
	input, err := os.Open(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer input.Close()

	// Create Output File
	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer output.Close()

	// Generate 32-Byte Key From Password
	key := make([]byte, 32)
	for i := 0; i < len(password); i++ {
		key[i%32] ^= password[i]
	}

	// Generate Random Initialization
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Write Initialization Vector to Output File
	_, err = output.Write(iv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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

	// Encrypt input file and write to output file
	encryptedWriter := &cipher.StreamWriter{S: stream, W: outputWriter}
	_, err = io.Copy(encryptedWriter, inputReader)
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

	// Encryption Complate Message
	fmt.Printf("Encryption complete: %s -> %s\n", inputFile, outputFile)
}
