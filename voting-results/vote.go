package main

import (
	"bytes"
	"encoding/gob"
)

type userVotes struct {
	Votes []string
}

func decode(data []byte) (*userVotes, error) {
	var uv *userVotes
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&uv)
	if err != nil {
		return nil, err
	}
	return uv, nil
}

func (uv *userVotes) encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(uv)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
