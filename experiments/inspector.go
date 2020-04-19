/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-19
 * @Project: Proof of Evolution
 * @Filename: inspector.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-19
 * @Copyright: 2020
 */



// To test source inspection
package main

import (
	"fmt"
	"go/parser"
	"go/token"
  "go/ast"
)

func main() {
	fset := token.NewFileSet() // positions are relative to fset

	src := `package foo

import (
	"fmt"
	"time"
)

func bar() {
  for i:=0;i<5;i++{
    fmt.Println(time.Now())
  }
	var f int = 5
	switch f{
	case 1: fmt.Println("1")
	case 2: fmt.Println("2")
	default: fmt.Println("0")
	}
  fmt.Println("End")
	if f==5{
		f = 0
	}else{
		f = 9
	}
}`

	// Parse src but stop after processing the imports.
	f, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		fmt.Println("MY error:",err)
		return
	}

	// Print the imports from the file's AST.
	for _, s := range f.Decls {
		fmt.Printf("%t\n\n",s)
	}

  ast.Inspect(f, func(n ast.Node) bool {
		var s string
		switch x := n.(type) {
		case *ast.BasicLit:
			s = "LIT_"+x.Value
		case *ast.Ident:
			s = "ID_"+x.Name
		}
		if s != "" {
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
		}
		return true
	})
  // Print the AST.
  ast.Print(fset, f.Decls[1].(*ast.FuncDecl).Body.List)

  fmt.Println(f.Decls[1].(*ast.FuncDecl).Name.Name)

}
