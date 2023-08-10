package main

import "log"

func invariant(variable string) {
	if variable == "" {
		log.Fatal("Missing env variable", variable)
	}
}
