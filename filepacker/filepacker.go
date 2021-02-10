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
	data, err := readData(infile)
	if err != nil {
		return err
	}

	size := fmt.Sprintf("%d", len(data))

	chksum, err := checksum(data)
	if err != nil {
		return err
	}

	compressedData, err := compress(data)
	if err != nil {
		return err
	}
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

	var buf bytes.Buffer
	if err := t.Execute(&buf, content); err != nil {
		return err
	}

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

// this is probably feasible in a more efficient way, but hey, it's something
func bytesToHexString(data []byte) string {
	str := ""
	var value []byte = make([]byte, 1, 1)
	for _, b := range data {
		value[0] = b
		str += fmt.Sprintf("%s%v", prefix, hex.EncodeToString(value))
	}
	return str
}

func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
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
