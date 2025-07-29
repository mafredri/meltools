package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func padKey(key []byte) []byte {
	// Pad the key to 16 bytes with zeroes if it's shorter.
	p := make([]byte, 16)
	copy(p, key)
	return p
}

func padPKCS7(data []byte, blockSize int) []byte {
	n := blockSize - len(data)%blockSize
	return append(data, bytes.Repeat([]byte{byte(n)}, n)...)
}

func encryptAESCBC(plain, key []byte) ([]byte, error) {
	key = padKey(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	padded := padPKCS7(plain, aes.BlockSize)
	enc := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(enc, padded)
	return append(iv, enc...), nil
}

func decryptAESCBC(enc, key []byte) ([]byte, error) {
	if len(enc) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := enc[:aes.BlockSize]
	enc = enc[aes.BlockSize:]
	key = padKey(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(enc)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("enc is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	dec := make([]byte, len(enc))
	mode.CryptBlocks(dec, enc)
	return unpadPKCS7(dec), nil
}

func unpadPKCS7(data []byte) []byte {
	// ISO/IEC 7816-4: last 0x80, then zero or more 0x00.
	if len(data) == 0 {
		return data
	}
	end := len(data)
	for end > 0 && data[end-1] == 0x00 {
		end--
	}
	if end > 0 && data[end-1] == 0x80 {
		end--
	}
	return data[:end]
}

func decodeESV(esv string, key []byte) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(esv)
	if err != nil {
		return nil, err
	}
	plain, err := decryptAESCBC(b, key)
	if err != nil {
		return nil, err
	}
	return plain, nil
}

type CSV struct {
	XMLName xml.Name `xml:"CSV"`
	Connect string   `xml:"CONNECT,omitempty"`
	Reset   string   `xml:"RESET,omitempty"`
	Code    []string `xml:"CODE>VALUE,omitempty"`
	ECHONET string   `xml:"ECHONET,omitempty"`
}

type ESV struct {
	XMLName xml.Name `xml:"ESV"`
	Data    string   `xml:",chardata"`
}

func main() {
	key := flag.String("key", "unregistered", "AES key")
	reset := flag.Bool("reset", false, "Reset the device (default false)")
	echonetOn := flag.Bool("enable-echonet", false, "Enable ECHONET (default false)")
	echonetOff := flag.Bool("disable-echonet", false, "Disable ECHONET (default false)")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <host>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	host := flag.Arg(0)

	csv := CSV{Connect: "ON"}
	if *reset {
		csv.Reset = "ON"
	}
	if *echonetOn {
		csv.ECHONET = "ON"
	}
	if *echonetOff {
		csv.ECHONET = "OFF"
	}

	// Marshal CSV request and encrypt it into ESV local command.
	csvData, err := xml.MarshalIndent(csv, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "CSV marshal error: %v\n", err)
		os.Exit(1)
	}
	// The XML header is not needed for the CSV data
	// but if you want to include it, uncomment the next line:
	// csvData = append([]byte(xml.Header), csvData...)
	fmt.Printf("CSV data to encrypt:\n%s\n", string(csvData))

	enc, err := encryptAESCBC(csvData, []byte(*key))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Encryption error: %v\n", err)
		os.Exit(1)
	}
	b64 := base64.StdEncoding.EncodeToString(enc)
	esv := ESV{Data: b64}
	esvData, err := xml.MarshalIndent(esv, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "XML marshal error: %v\n", err)
		os.Exit(1)
	}
	esvData = append([]byte(xml.Header), esvData...)
	fmt.Printf("Sending encrypted data:\n%s\n", string(esvData))

	url := fmt.Sprintf("http://%s/smart", host)
	resp, err := http.Post(url, "text/xml", bytes.NewReader(esvData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP POST error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "Error: received non-OK HTTP status")
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Println(string(respBody))
		os.Exit(1)
	}

	// Decrypt ESV response.
	dec := xml.NewDecoder(resp.Body)
	var esvResp ESV
	err = dec.Decode(&esvResp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ESV decode error: %v\n", err)
		os.Exit(1)
	}

	if esvResp.Data != "" {
		plain, err := decodeESV(esvResp.Data, []byte(*key))
		if err != nil {
			fmt.Fprintf(os.Stderr, "ESV decode error: %v\n", err)
		} else {
			fmt.Printf("\nDecrypted ESV response:\n%s\n", plain)
		}
	}
}
