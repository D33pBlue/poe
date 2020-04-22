/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-22
 * @Project: Proof of Evolution
 * @Filename: merkle.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-23
 * @Copyright: 2020
 */

// Package merkle defines the Merkle tree that
// is used to store the transactions of the blockchain
package merkle

import(
	"fmt"
	"crypto/sha256"
 	"github.com/D33pBlue/poe/block"
)

type Node struct{
	Parent,L,R *Node
	Transaction block.Transact
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
	var input1 []byte = append(n.L.Hash,n.R.Hash...)
	fmt.Printf("%x\n",input1)
	encoder := sha256.New()
	encoder.Write([]byte(input1))
	var result []byte = encoder.Sum(nil)
	for i:=0;i<len(n.Hash);i++{
		if n.Hash[i]!=result[i]{
			return false
		}
	}
	return checkSubTree(n.L) && checkSubTree(n.R)
}

func (self *Tree)Add(trans block.Transact){
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
		var input1 []byte = append(n.L.Hash,n.R.Hash...)
		encoder := sha256.New()
		encoder.Write([]byte(input1))
		n.Hash = encoder.Sum(nil)
		n.Children = n.L.Children+n.R.Children+2
		n = n.Parent
	}
}
