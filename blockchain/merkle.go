/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-22
 * @Project: Proof of Evolution
 * @Filename: merkle.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-24
 * @Copyright: 2020
 */

// Package merkle defines the Merkle tree that
// is used to store the transactions of the transactchain
package blockchain

import(
	"github.com/D33pBlue/poe/utils"
)

type Node struct{
	Parent,L,R *Node
	Transaction Transaction
	Hash []byte
	Children int
}

type Tree struct{
	Root *Node
	Nleaves int
}

func New()*Tree{
	m := new(Tree)
	m.Nleaves = 0
	return m
}

func (self *Tree)Check()bool{
	return checkSubTree(self.Root)
}

func checkSubTree(n *Node)bool{
	if n==nil {return true}
	if n.L==nil && n.R==nil {return true}
	hashBuilder := new(utils.HashBuilder)
	hashBuilder.Add(n.L.Hash)
	hashBuilder.Add(n.R.Hash)
	var result []byte = hashBuilder.GetHash()
	for i:=0;i<len(n.Hash);i++{
		if n.Hash[i]!=result[i]{
			return false
		}
	}
	hashBuilder = nil
	return checkSubTree(n.L) && checkSubTree(n.R)
}

func (self *Tree)Add(trans Transaction){
	var n *Node = new(Node)
	n.Transaction = trans
	n.Hash = trans.GetHash()
	n.Children = 0
	self.insertNode(n)
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
		p.Parent = self.Root
		self.Root.R = n
		n.Parent = self.Root
		self.Root.Children = 2
		self.Nleaves += 1
		updateHashes(self.Root)
		return
	}
	if p.isFull(){
		self.Root = new(Node)
		self.Root.L = p
		p.Parent = self.Root
		self.Root.R = n
		n.Parent = self.Root
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
	p.Parent = x1.Parent
	p.Parent.L = p
	p.L = x1
	x1.Parent = p
	p.R = x2
	x2.Parent = p
	p.Children = x1.Children+x2.Children+2
	self.Nleaves += 1
	updateHashes(p)
}

func updateHashes(n *Node){
	for ;n!=nil;{
		hashBuilder := new(utils.HashBuilder)
		hashBuilder.Add(n.L.Hash)
		hashBuilder.Add(n.R.Hash)
		n.Hash = hashBuilder.GetHash()
		n.Children = n.L.Children+n.R.Children+2
		n = n.Parent
	}
}
