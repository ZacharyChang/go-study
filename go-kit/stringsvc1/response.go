package main

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}

type countResponse struct {
	V int `json:"v"`
}
