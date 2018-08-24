package library

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"

	"github.com/decosblockchain/audittrail-client/library/fastsha256"
)

type SortedHashSlice []*big.Int

func (a SortedHashSlice) Len() int           { return len(a) }
func (a SortedHashSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortedHashSlice) Less(i, j int) bool { return a[i].Cmp(a[j]) == -1 }

type MerkleNode struct {
	leaf      bool
	leftVal   *big.Int
	rightVal  *big.Int
	root      *big.Int
	leftNode  *MerkleNode
	rightNode *MerkleNode
}

func NewMerkleLeaf(left, right *big.Int) *MerkleNode {
	node := new(MerkleNode)
	node.leaf = true
	node.leftVal = left
	node.rightVal = right
	node.root = node.calcRoot()
	return node
}

func NewMerkleNode(left, right *MerkleNode) *MerkleNode {
	node := new(MerkleNode)
	node.leaf = false
	node.leftNode = left
	node.rightNode = right
	node.root = node.calcRoot()
	return node
}

func (n *MerkleNode) getLeftVal() *big.Int {
	if n.leaf {
		return n.leftVal
	} else {
		return n.leftNode.MerkleRoot()
	}
}

func (n *MerkleNode) getRightVal() *big.Int {
	if n.leaf {
		return n.rightVal
	} else {
		return n.rightNode.MerkleRoot()
	}
}

func (n *MerkleNode) MerkleRoot() *big.Int {
	return n.root
}

func (n *MerkleNode) calcRoot() *big.Int {
	left := n.getLeftVal()
	right := n.getRightVal()

	hash := fastsha256.Sum256([]byte(fmt.Sprintf("%x%x", left, right)))

	root, _ := new(big.Int).SetString(hex.EncodeToString(hash[:]), 16)
	return root
}

func MakeMerkleTree(hashes []*big.Int) *MerkleNode {
	sort.Sort(SortedHashSlice(hashes))
	leafQueue := make([]*big.Int, 0)
	nodes := make([]*MerkleNode, 0)
	for _, hash := range hashes {
		if len(leafQueue) < 2 {
			leafQueue = append(leafQueue, hash)
		} else {
			nodes = append(nodes, NewMerkleLeaf(leafQueue[0], leafQueue[1]))
			leafQueue = []*big.Int{hash}
		}
	}

	if len(leafQueue) > 0 {
		nodes = append(nodes, NewMerkleLeaf(leafQueue[0], leafQueue[len(leafQueue)-1]))
	}

	nodeQueue := make([]*MerkleNode, 0)
	newNodes := make([]*MerkleNode, 0)
	for len(nodes) > 1 {
		for _, node := range nodes {
			if len(nodeQueue) < 2 {
				nodeQueue = append(nodeQueue, node)
			} else {
				newNodes = append(newNodes, NewMerkleNode(nodeQueue[0], nodeQueue[1]))
				nodeQueue = []*MerkleNode{node}
			}
		}

		if len(nodeQueue) > 0 {
			newNodes = append(newNodes, NewMerkleNode(nodeQueue[0], nodeQueue[len(nodeQueue)-1]))
		}

		nodes = newNodes
	}

	return nodes[0]
}
