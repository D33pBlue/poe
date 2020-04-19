/**
 * @Author: d33pblue
 * @Date:   2020-Apr-19
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-19
 * @Copyright: 2020
 */

package ga

import(
  "fmt"
  "plugin"
  "os/exec"
  "go/ast"
  "go/token"
  "go/parser"
  "errors"
)

// Problem defines the interface the user has to
// implement to declare his problem.
// Initialize method is called at the beginning, and then
// New has to return an instance of the DNA interface.
type Problem interface {
	Initialize(path string)
  New()DNA
}

// Checks if the package is "main" in order
// to compile the source as plugin.
func checkPackage(f *ast.File)error{
  if f.Name.Name!="main"{
    return errors.New("The defined package is not main")
  }
  return nil
}

// Checks the user does not import forbidden modules.
func checkImport(decl *ast.ImportSpec)error{
  var whitelist = []string{"math/rand","io/ioutil",
    "encoding/json","github.com/D33pBlue/poe/op",
    "github.com/D33pBlue/poe/ga","fmt"}
  var name string = decl.Path.Value
  if len(name)>0 && name[0]=='"' {
    name = name[1:]
  }
  if len(name)>0 && name[len(name)-1]=='"' {
    name = name[:len(name)-1]
  }
  var found bool = false
  for i:=0; !found && i<len(whitelist); i++{
    if name==whitelist[i] {
      found = true
    }
  }
  if !found{
    return errors.New("Imported "+name+", which is not in white list.")
  }
  return nil
}

// Checks the functions defined by the user.
func checkFunc(decl *ast.FuncDecl)error{
  // TODO: check for conditions at least
  return nil
}

// Checks the source defined by the user.
func checkDeclarations(f *ast.File)error{
  for _,decl := range f.Decls{
    switch decl.(type) {
    case *ast.GenDecl:
      switch fmt.Sprint(decl.(*ast.GenDecl).Tok) {
      case "import":
        for _,im := range decl.(*ast.GenDecl).Specs{
          err := checkImport(im.(*ast.ImportSpec))
          if err!=nil{return err}
        }
      // case "type": fmt.Println("TYPE")
      // case "var": fmt.Println("VV")
      }
    case *ast.FuncDecl:
      err := checkFunc(decl.(*ast.FuncDecl))
      if err!=nil{return err}
    }
  }
  return nil
}

// Parse and inspect the source defined by the user.
func inspect(dir,name string)error{
  fset := token.NewFileSet()
  f, err := parser.ParseFile(fset,dir+name+".go",nil, parser.AllErrors)
	if err != nil {
    return err
  }
  // ast.Print(fset,f)
  err = checkPackage(f)
  if err!=nil {return err}
  err = checkDeclarations(f)
  if err!=nil {return err}
  return nil
}

// Compiles a plugin, after checking it.
func compilePlugin(dir,name string)error{
  err := inspect(dir,name)
  if err!=nil{
    return err
  }
  baseName := dir+name
  cmd := exec.Command("go","build","-buildmode=plugin",
    "-o",baseName+".so",baseName+".go")
  cmdstr := fmt.Sprint("go build -buildmode=plugin -o "+baseName+".so "+baseName+".go")
  // fmt.Println("Compiling plugin..")
	err = cmd.Run()
  if err!=nil{
    return errors.New("Error in plugin compilation:\nrun the following command for more details:\n\t"+cmdstr)
  }
  return nil
}

// Returns a DNA after checking and compiling
// a user-defined plugin.
func LoadGA(plugDir,plugName,path2data string)(DNA,error){
  err := compilePlugin(plugDir,plugName)
  if err != nil { return nil,err}
  plug, err := plugin.Open(plugDir+plugName+".so")
  if err != nil {return nil,err}
  definition, err := plug.Lookup("Definition")
  if err != nil {return nil,err}
  var problem Problem
  problem, ok := definition.(Problem)
  if !ok {
    err = errors.New("The module does not implement Problem interface")
    return nil,err
  }
  problem.Initialize(path2data)
  problDna := problem.New()
  if problDna==nil{
    err = errors.New("The function New defined in the module does not return a valid DNA")
    return nil,err
  }
  return problDna,nil
}
