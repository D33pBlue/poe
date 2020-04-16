package op

import(
  "fmt"
)

func (self *State)Eq(a,b any)bool{
  self.IncOperations(self.coeff["=="]+self.off["=="])
  var t string = fmt.Sprintf("%T", a)
  switch t {
    case "int": return a.(int)==b.(int)
    case "int8": return a.(int8)==b.(int8)
    case "int16": return a.(int16)==b.(int16)
    case "int32": return a.(int32)==b.(int32)
    case "int64": return a.(int64)==b.(int64)
    case "uint": return a.(uint)==b.(uint)
    case "uint8": return a.(uint8)==b.(uint8)
    case "uint16": return a.(uint16)==b.(uint16)
    case "uint32": return a.(uint32)==b.(uint32)
    case "uint64": return a.(uint64)==b.(uint64)
    case "float32": return a.(float32)==b.(float32)
    case "float64": return a.(float64)==b.(float64)
    case "complex64": return a.(complex64)==b.(complex64)
    case "complex128": return a.(complex128)==b.(complex128)
    case "string": return a.(string)==b.(string)
    default: fmt.Println("Invalid type")
    }
  return false
}

func (self *State)Le(a,b any)bool{
  self.IncOperations(self.coeff["<="]+self.off["<="])
  var t string = fmt.Sprintf("%T", a)
  switch t {
    case "int": return a.(int)<=b.(int)
    case "int8": return a.(int8)<=b.(int8)
    case "int16": return a.(int16)<=b.(int16)
    case "int32": return a.(int32)<=b.(int32)
    case "int64": return a.(int64)<=b.(int64)
    case "uint": return a.(uint)<=b.(uint)
    case "uint8": return a.(uint8)<=b.(uint8)
    case "uint16": return a.(uint16)<=b.(uint16)
    case "uint32": return a.(uint32)<=b.(uint32)
    case "uint64": return a.(uint64)<=b.(uint64)
    case "float32": return a.(float32)<=b.(float32)
    case "float64": return a.(float64)<=b.(float64)
    case "string": return a.(string)<=b.(string)
    default: fmt.Println("Invalid type")
    }
  return false
}

func (self *State)Ge(a,b any)bool{
  self.IncOperations(self.coeff[">="]+self.off[">="])
  var t string = fmt.Sprintf("%T", a)
  switch t {
    case "int": return a.(int)>=b.(int)
    case "int8": return a.(int8)>=b.(int8)
    case "int16": return a.(int16)>=b.(int16)
    case "int32": return a.(int32)>=b.(int32)
    case "int64": return a.(int64)>=b.(int64)
    case "uint": return a.(uint)>=b.(uint)
    case "uint8": return a.(uint8)>=b.(uint8)
    case "uint16": return a.(uint16)>=b.(uint16)
    case "uint32": return a.(uint32)>=b.(uint32)
    case "uint64": return a.(uint64)>=b.(uint64)
    case "float32": return a.(float32)>=b.(float32)
    case "float64": return a.(float64)>=b.(float64)
    case "string": return a.(string)>=b.(string)
    default: fmt.Println("Invalid type")
    }
  return false
}

func (self *State)Lt(a,b any)bool{
  self.IncOperations(self.coeff["<"]+self.off["<"])
  var t string = fmt.Sprintf("%T", a)
  switch t {
    case "int": return a.(int)<b.(int)
    case "int8": return a.(int8)<b.(int8)
    case "int16": return a.(int16)<b.(int16)
    case "int32": return a.(int32)<b.(int32)
    case "int64": return a.(int64)<b.(int64)
    case "uint": return a.(uint)<b.(uint)
    case "uint8": return a.(uint8)<b.(uint8)
    case "uint16": return a.(uint16)<b.(uint16)
    case "uint32": return a.(uint32)<b.(uint32)
    case "uint64": return a.(uint64)<b.(uint64)
    case "float32": return a.(float32)<b.(float32)
    case "float64": return a.(float64)<b.(float64)
    case "string": return a.(string)<b.(string)
    default: fmt.Println("Invalid type")
    }
  return false
}

func (self *State)Gt(a,b any)bool{
  self.IncOperations(self.coeff[">"]+self.off[">"])
  var t string = fmt.Sprintf("%T", a)
  switch t {
    case "int": return a.(int)>b.(int)
    case "int8": return a.(int8)>b.(int8)
    case "int16": return a.(int16)>b.(int16)
    case "int32": return a.(int32)>b.(int32)
    case "int64": return a.(int64)>b.(int64)
    case "uint": return a.(uint)>b.(uint)
    case "uint8": return a.(uint8)>b.(uint8)
    case "uint16": return a.(uint16)>b.(uint16)
    case "uint32": return a.(uint32)>b.(uint32)
    case "uint64": return a.(uint64)>b.(uint64)
    case "float32": return a.(float32)>b.(float32)
    case "float64": return a.(float64)>b.(float64)
    case "string": return a.(string)>b.(string)
    default: fmt.Println("Invalid type")
    }
  return false
}

func (self *State)Neg(a any)any{
  self.IncOperations(self.coeff["neg"]+self.off["neg"])
  var t string = fmt.Sprintf("%T", a)
  switch t {
    case "int": return -a.(int)
    case "int8": return -a.(int8)
    case "int16": return -a.(int16)
    case "int32": return -a.(int32)
    case "int64": return -a.(int64)
    case "uint": return -a.(uint)
    case "uint8": return -a.(uint8)
    case "uint16": return -a.(uint16)
    case "uint32": return -a.(uint32)
    case "uint64": return -a.(uint64)
    case "float32": return -a.(float32)
    case "float64": return -a.(float64)
    case "complex64": return -a.(complex64)
    case "complex128": return -a.(complex128)
    case "bool": return !a.(bool)
    default: fmt.Println("Invalid type")
    }
  return nil
}

func (self *State)And(a,b any)bool{
  self.IncOperations(self.coeff["and"]+self.off["and"])
  var t string = fmt.Sprintf("%T", a)
  switch t {
    case "bool": return (a.(bool) && b.(bool))
    default: fmt.Println("Invalid type")
    }
  return false
}

func (self *State)Or(a,b any)bool{
  self.IncOperations(self.coeff["or"]+self.off["or"])
  var t string = fmt.Sprintf("%T", a)
  switch t {
    case "bool": return (a.(bool) || b.(bool))
    default: fmt.Println("Invalid type")
    }
  return false
}
