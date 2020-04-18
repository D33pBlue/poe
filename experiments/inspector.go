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
  fmt.Println("End")
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
