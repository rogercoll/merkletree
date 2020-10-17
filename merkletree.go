package merkletree

import (
	"encoding/hex"
	"hash"
	"io/ioutil"
	"strconv"
)

type node struct {
	//name is just for the exercise nomenclature
	name    string
	parent  *node
	left    *node
	right   *node
	leaf    bool
	content []byte
	i       int
	j       int
}

type merkletree struct {
	root     *node
	leafs    *[]node
	middle   *[]node
	hashAlgo func() hash.Hash
}

func readData(files []string) (*[][]byte, error) {
	contents := make([][]byte, len(files))
	for i, file := range files {
		dat, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		contents[i] = dat
	}
	return &contents, nil
}

func computeMiddleNodes(nodes *[]node, m *merkletree, level int) (*node, error) {
	newlevel := make([]node, len(*nodes)/2)
	level++
	for iter := 0; iter <= len(*nodes)/2; iter += 2 {
		hash := (*m).hashAlgo()
		left := iter
		right := iter + 1
		newhash := append((*nodes)[left].content, (*nodes)[right].content...)
		hash.Write(newhash)
		newNode := node{
			content: hash.Sum(nil),
			leaf:    false,
			j:       iter / 2,
			i:       level,
		}
		newlevel[iter/2] = newNode
		if len(*nodes) == 2 {
			return &newNode, nil
		}
		*m.middle = append(*m.middle, newNode)
	}
	return computeMiddleNodes(&newlevel, m, level)
}

func computeMerkleTree(data *[][]byte, hashAlgorithm func() hash.Hash) (*merkletree, error) {
	m := merkletree{hashAlgo: hashAlgorithm}
	leafs := make([]node, len(*data))
	iter := 0
	for _, fdata := range *data {
		hash := m.hashAlgo()
		hash.Write(fdata)
		leafs[iter] = node{
			content: hash.Sum(nil),
			leaf:    true,
			i:       0,
			j:       iter,
		}
		iter++
	}
	if len(leafs)%2 == 1 {
		leafs = append(leafs, leafs[len(leafs)-1])
	}
	m.leafs = &leafs
	m.middle = &[]node{}
	var err error
	m.root, err = computeMiddleNodes(&leafs, &m, 0)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (m *merkletree) GetRoot() (string, int) {
	return hex.EncodeToString(m.root.content[:]), m.root.i
}

func (m *merkletree) GetPrivateTree(hashName string) string {
	output := "MerkleTree:" + hashName + ":" + strconv.Itoa(len(*m.leafs)) + ":" + strconv.Itoa(m.root.i) + ":" + hex.EncodeToString(m.root.content[:]) + "\n"
	for _, leaf := range *m.leafs {
		output += strconv.Itoa(leaf.i) + ":" + strconv.Itoa(leaf.j) + ":" + hex.EncodeToString(leaf.content[:]) + "\n"
	}
	for _, middleNode := range *m.middle {
		output += strconv.Itoa(middleNode.i) + ":" + strconv.Itoa(middleNode.j) + ":" + hex.EncodeToString(middleNode.content[:]) + "\n"
	}
	return output[:len(output)-1]
}

func Build(files []string, hashAlgorithm func() hash.Hash) (*merkletree, error) {
	content, err := readData(files)
	if err != nil {
		return nil, err
	}
	m, err := computeMerkleTree(content, hashAlgorithm)

	return m, err
}
