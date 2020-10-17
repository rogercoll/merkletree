package main

import (
	"crypto/sha256"
	"github.com/rogercoll/merkletree"
	"log"
)

func main() {
	var sha256func = sha256.New
	files := []string{
		"./data/docs/doc0.dat",
		"./data/docs/doc1.dat",
		"./data/docs/doc2.dat",
		"./data/docs/doc3.dat",
	}
	m, err := merkletree.Build(files, sha256func)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(m.GetRoot())
}
