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
	middle   map[int][]node
	hashAlgo func() hash.Hash
}

func ReadFile(file string) ([]byte, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func readData(files []string) (*[][]byte, error) {
	contents := make([][]byte, len(files))
	for i, file := range files {
		dat, err := ReadFile(file)
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
		(*m).middle[level] = append((*m).middle[level], newNode)
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
	middle := make(map[int][]node, len(leafs)/2-2)
	m.middle = middle
	var err error
	m.root, err = computeMiddleNodes(&leafs, &m, 0)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (m *merkletree) GetRoot() string {
	return hex.EncodeToString(m.root.content[:])
}

func (m *merkletree) GetPrivateTree(hashName string) string {
	output := "MerkleTree:" + hashName + ":" + strconv.Itoa(len(*m.leafs)) + ":" + strconv.Itoa(m.root.i) + ":" + hex.EncodeToString(m.root.content[:]) + "\n"
	for _, leaf := range *m.leafs {
		output += strconv.Itoa(leaf.i) + ":" + strconv.Itoa(leaf.j) + ":" + hex.EncodeToString(leaf.content[:]) + "\n"
	}
	for iter := 1; iter < m.root.i; iter++ {
		for _, node := range (*m).middle[iter] {
			output += strconv.Itoa(node.i) + ":" + strconv.Itoa(node.j) + ":" + hex.EncodeToString(node.content[:]) + "\n"
		}
	}
	return output[:len(output)-1]
}

func (m *merkletree) proof(data []byte, level, col int) []byte {
	if level == m.root.i {
		return data
	}
	if level == 0 {
		pleafs := (*m).leafs
		if col%2 == 0 {
			data = append(data, (*pleafs)[col+1].content[:]...)
		} else {
			data = append((*pleafs)[col-1].content[:], data...)
		}
	} else {
		if col%2 == 0 {
			data = append(data, (*m).middle[level][col+1].content[:]...)
		} else {
			data = append((*m).middle[level][col-1].content[:], data...)
		}
	}
	hash := m.hashAlgo()
	hash.Write(data)
	nodeHash := hash.Sum(nil)
	level++
	return m.proof(nodeHash, level, col/2)
}

func (m *merkletree) ProofMembership(data []byte, leaf int) string {
	hash := m.hashAlgo()
	hash.Write(data)
	nodeHash := hash.Sum(nil)
	tmpRoot := m.proof(nodeHash, 0, leaf)
	return hex.EncodeToString(tmpRoot[:])
}

func Build(files []string, hashAlgorithm func() hash.Hash) (*merkletree, error) {
	content, err := readData(files)
	if err != nil {
		return nil, err
	}
	m, err := computeMerkleTree(content, hashAlgorithm)

	return m, err
}
