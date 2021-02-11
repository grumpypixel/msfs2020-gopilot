package filepacker

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Content struct {
	Tool     string
	Size     string
	Checksum string
	Package  string
	Func     string
	Data     string
}

const (
	toolName = "filepacker"
	prefix   = "\\x"
)

// https://gist.github.com/alex-ant/aeaaf497055590dacba760af24839b8d
func Pack(infile, outfile, templateFile, packageName, funcName string) error {
	// fmt.Println(" Reading...")
	data, err := readData(infile)
	if err != nil {
		return err
	}

	size := fmt.Sprintf("%d", len(data))

	// fmt.Println(" Checksumming...")
	chksum, err := checksum(data)
	if err != nil {
		return err
	}

	// fmt.Println(" Compressing...")
	compressionLevel := gzip.BestCompression
	compressedData, err := compress(data, compressionLevel)
	if err != nil {
		return err
	}

	// fmt.Println(" Hexifying...")
	hexifiedData := bytesToHexString(compressedData)

	templateBytes, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return err
	}

	t, err := template.New("gopher").Parse(string(templateBytes))
	if err != nil {
		return err
	}

	content := Content{
		Tool:     toolName,
		Size:     size,
		Checksum: chksum,
		Package:  packageName,
		Func:     funcName,
		Data:     hexifiedData,
	}

	// fmt.Println(" Templating...")
	var buf bytes.Buffer
	if err := t.Execute(&buf, content); err != nil {
		return err
	}

	// fmt.Println(" Writing...")
	if err := writeData(outfile, buf.String()); err != nil {
		return err
	}

	fmt.Println(" Input file:", infile)
	fmt.Println(" Ouput file:", outfile)
	fmt.Println(" Template:", templateFile)
	fmt.Println(" Checksum:", chksum)
	fmt.Println(" Size:", size)
	fmt.Println(" Package:", packageName)
	fmt.Println(" Getter:", funcName)

	return nil
}

func Unpack(content []byte) ([]byte, error) {
	data, err := decompress([]byte(content))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func readData(filename string) ([]byte, error) {
	input, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer input.Close()

	fileStats, err := input.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileStats.Size()
	data := make([]byte, fileSize)

	bytesRead, err := input.Read(data)
	if err != nil {
		return nil, err
	}
	if int64(bytesRead) != fileSize {
		return nil, fmt.Errorf("Bytes read to filesize mismatch")
	}
	return data, nil
}

func writeData(filename, content string) error {
	output, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer output.Close()
	output.WriteString(content)
	return nil
}

func checksum(data []byte) (string, error) {
	hasher := sha256.New()
	if _, err := hasher.Write(data); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func bytesToHexString(data []byte) string {
	var sb strings.Builder
	// count := len(data) + len(prefix)*len(data)
	// sb.Grow(count)
	// fmt.Println("Count", count)

	var value []byte = make([]byte, 1, 1)
	for _, b := range data {
		value[0] = b
		sb.WriteString(fmt.Sprintf("%s%v", prefix, hex.EncodeToString(value)))
	}
	return sb.String()
}

func compress(data []byte, compressionLevel int) ([]byte, error) {
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(&buf, compressionLevel)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}
	if err = w.Flush(); err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decompress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	var r io.Reader
	r, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	var res bytes.Buffer
	_, err = res.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}
