package main

import "log"

func variableMissing(variable string) {
	if variable == "" {
		log.Fatal("Missing env variable", variable)
	}
}
