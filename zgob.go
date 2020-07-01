// zgob is a tiny utility library for saving and loading gobs with zlib compression
package zgob

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"
	"os"
	"reflect"
)

// RegisterTypes ...
func RegisterTypes(types ...interface{}) {
	for _, t := range types {
		gob.Register(t)
	}
}

func isPointer(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.Ptr
}

// Save encodes v and compresses it with zlib. Compressed data is written to filename.
func Save(v interface{}, filename string) error {
	fp, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer fp.Close()

	encBuffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&encBuffer)

	if err := encoder.Encode(v); err != nil {
		return err
	}

	zBuffer := bytes.Buffer{}
	zWriter := zlib.NewWriter(&zBuffer)
	zWriter.Write(encBuffer.Bytes())

	if err := zWriter.Close(); err != nil {
		return err
	}

	fp.Write(zBuffer.Bytes())
	fp.Sync()

	return nil
}

// Load takes a pointer v and attempts to decompress and decode data from filename.
func Load(v interface{}, filename string) (bytesRead int64, err error) {
	if !isPointer(v) {
		panic("v is not a pointer!")
	}

	fp, err := os.Open(filename)

	if err != nil {
		return -1, err
	}

	defer fp.Close()

	zBuffer := bytes.Buffer{}
	zReader, err := zlib.NewReader(fp)

	if err != nil {
		return -1, err
	}

	defer zReader.Close()

	bytesRead, err = zBuffer.ReadFrom(zReader)

	if err != nil {
		return -1, err
	}

	decoder := gob.NewDecoder(&zBuffer)
	err = decoder.Decode(v)

	if err != nil {
		return -1, err
	}

	return bytesRead, nil
}
