package main

var rule Config

type Config struct {
	Version string
	Rules   []Rule
}

type Rule struct {
	Match  string
	Rename string
	Junk   []string
}
