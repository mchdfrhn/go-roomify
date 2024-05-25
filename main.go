package main

import (
	"go-roomify/delivery"
)

func main() {
	delivery.NewServer().Run()
}