package ga

import(
  "fmt"
  "plugin"
  "os/exec"
)

type Problem interface {
	Initialize(path string)
  New()DNA
}

func compilePlugin(dir,name string)error{
  baseName := dir+name
  cmd := exec.Command("go","build","-buildmode=plugin",
    "-o",baseName+".so",baseName+".go")
  fmt.Println("go","build","-buildmode=plugin",
    "-o",baseName+".so",baseName+".go")
  fmt.Println("Compiling plugin..")
	return cmd.Run()
}

func LoadGA(plugDir,plugName,path2data string)DNA{
  err := compilePlugin(plugDir,plugName)
  if err != nil {panic(err)}
  plug, err := plugin.Open(plugDir+plugName+".so")
  if err != nil {panic(err)}
  definition, err := plug.Lookup("Definition")
  if err != nil {panic(err)}
  var problem Problem
  problem, ok := definition.(Problem)
  if !ok {
    panic("The module does not implement Problem interface")
  }
  problem.Initialize(path2data)
  problDna := problem.New()
  return problDna
}
