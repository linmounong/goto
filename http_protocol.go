package main

import (
	"encoding/json"
	"io"
)

type Response struct {
	Username string  `json:"username,omitempty"`
	Ok       bool    `json:"ok,omitempty"`
	ErrMsg   string  `json:"errMsg,omitempty"`
	Link     *Link   `json:"link,omitempty"`
	Links    []*Link `json:"links,omitempty"`
}

func (r *Response) marshal(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(r)
}

type Link struct {
	Key      string `json:"key,omitempty"`
	Url      string `json:"url,omitempty"`
	Owner    string `json:"owner,omitempty"`
	UseCount int    `json:"useCount,omitempty"`
}
