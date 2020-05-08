/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-22
 * @Project: Proof of Evolution
 * @Filename: merkle.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-08
 * @Copyright: 2020
 */

// Package merkle defines the Merkle tree that
// is used to store the transactions of the transactchain
package blockchain

import(
	"fmt"
	"encoding/json"
	"github.com/D33pBlue/poe/utils"
)

type Node struct{
	parent,L,R *Node
	Type string
	Transaction Transaction
	Hash string
	Children int
}

type Tree struct{
	Root *Node
	Nleaves int
	transactions []Transaction
}

func BuildMerkleTree()*Tree{
	m := new(Tree)
	m.Nleaves = 0
	return m
}

func (self *Tree)GetHash()string{
	if self.Root==nil{ return "" }
	return self.Root.Hash
}

// Checks the hashes of the Merkle tree.
func (self *Tree)Check()bool{
	return checkSubTree(self.Root)
}

// Add a Transaction to the merkle tree.
// The transaction is not checked here: it is assumed
// to be valid.
func (self *Tree)Add(trans Transaction){
	self.transactions = append(self.transactions,trans)
	var n *Node = new(Node)
	n.Transaction = trans
	n.Type = trans.GetType()
	n.Hash = trans.GetHashCached()
	n.Children = 0
	self.insertNode(n)
}

// Returns an array with all the Transactions inside the tree.
func (self *Tree)GetTransactionArray()[]Transaction{
	return self.transactions
}

func (self *Tree)PruneSpentTransactions()  {
	// TODO: implement later
	// remove also from self.transactions
}

func marshalTransaction(data []byte,tp string)Transaction{
	if len(data)<=0{ return nil }
	var transact Transaction = nil
	switch tp {
	case TrCoin:
		transact = MarshalCoinTransaction(data)
	case TrStd:
		transact = MarshalStdTransaction(data)
	case TrJob:
		transact = MarshalJobTransaction(data)
	}
	return transact
}

func marshalMerkleNode(data []byte,parent *Node)(node *Node,transactions []Transaction){
	if len(data)<=0{ return }
	node = new(Node)
	var objmap map[string]json.RawMessage
  json.Unmarshal(data, &objmap)
	node.parent = parent
	json.Unmarshal(objmap["Type"],&node.Type)
	json.Unmarshal(objmap["Children"],&node.Children)
	json.Unmarshal(objmap["Hash"],&node.Hash)
	node.Transaction = marshalTransaction(objmap["Transaction"],node.Type)
	if node.Transaction!=nil{
		transactions = append(transactions,node.Transaction)
	}
	var trs []Transaction
	node.L,trs = marshalMerkleNode(objmap["L"],node)
	for i:=0;i<len(trs);i++{
		if trs[i]!=nil{
			transactions = append(transactions,trs[i])
		}
	}
	node.R,trs = marshalMerkleNode(objmap["R"],node)
	for i:=0;i<len(trs);i++{
		if trs[i]!=nil{
			transactions = append(transactions,trs[i])
		}
	}
	return
}

func MarshalMerkleTree(data []byte)*Tree {
	tree := new(Tree)
	var objmap map[string]json.RawMessage
  json.Unmarshal(data, &objmap)
  json.Unmarshal(objmap["Nleaves"],&tree.Nleaves)
	tree.Root,tree.transactions = marshalMerkleNode(objmap["Root"],nil)
	return tree
}

func checkSubTree(n *Node)bool{
	if n==nil {return true}
	if n.L==nil && n.R==nil {return true}
	hashBuilder := new(utils.HashBuilder)
	hashBuilder.Add(n.L.Hash)
	hashBuilder.Add(n.R.Hash)
	result := fmt.Sprintf("%x",hashBuilder.GetHash())
	if result!=n.Hash{return false}
	return checkSubTree(n.L) && checkSubTree(n.R)
}

func (self *Node)isFull()bool{
	if self.Children==0{
		return true
	}
	return self.L.Children==self.R.Children
}

func (self *Tree)insertNode(n *Node){
	if self.Root==nil{
		self.Root = n
		self.Nleaves += 1
		return
	}
	p := self.Root
	if p.L==nil{
		self.Root = new(Node)
		self.Root.L = p
		p.parent = self.Root
		self.Root.R = n
		n.parent = self.Root
		self.Root.Children = 2
		self.Nleaves += 1
		updateHashes(self.Root)
		return
	}
	if p.isFull(){
		self.Root = new(Node)
		self.Root.L = p
		p.parent = self.Root
		self.Root.R = n
		n.parent = self.Root
		self.Root.Children = self.Root.L.Children+2
		self.Nleaves += 1
		updateHashes(self.Root)
		return
	}
	for ;!p.isFull();{
		p = p.R
	}
	x1 := p
	x2 := n
	p = new(Node)
	p.parent = x1.parent
	p.parent.L = p
	p.L = x1
	x1.parent = p
	p.R = x2
	x2.parent = p
	p.Children = x1.Children+x2.Children+2
	self.Nleaves += 1
	updateHashes(p)
}

func updateHashes(n *Node){
	for ;n!=nil;{
		hashBuilder := new(utils.HashBuilder)
		hashBuilder.Add(n.L.Hash)
		hashBuilder.Add(n.R.Hash)
		n.Hash = fmt.Sprintf("%x",hashBuilder.GetHash())
		n.Children = n.L.Children+n.R.Children+2
		n = n.parent
	}
}
