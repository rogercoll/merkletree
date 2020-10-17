package main

import (
	"crypto/sha256"
	"fmt"
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
	fmt.Printf("   ---MERKLE TREE FILE---   \n")
	fmt.Println(m.GetPrivateTree("sha256"))

	fmt.Printf("\n\n   ---MEMBERSHIP PROOFS---   \n")
	fmt.Printf("First test: doc3.dat tampered => ./data/docs/baddoc3.dat\n")
	bad, err := merkletree.ReadFile("./data/docs/baddoc3.dat")
	if err != nil {
		log.Fatal(err)
	}
	tmpRoot := m.ProofMembership(bad, 3)
	fmt.Printf("Returned root hash: %v\n", tmpRoot)
	fmt.Printf("Expected root hash: %v\n", m.GetRoot())

	fmt.Printf("\n\nSecond test: doc3.dat => ./data/docs/doc3.dat\n")
	doc3, err := merkletree.ReadFile("./data/docs/doc3.dat")
	if err != nil {
		log.Fatal(err)
	}
	tmpRoot2 := m.ProofMembership(doc3, 3)
	fmt.Printf("Returned root hash: %v\n", tmpRoot2)
	fmt.Printf("Expected root hash: %v\n", m.GetRoot())
}
