package main

var rule Config

type Config struct {
	Version     string
	Junk        []string
	Rules       []Rule
	Destination string
}

type Rule struct {
	Match  []string
	Rename string
	Dir    string
}
