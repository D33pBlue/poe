/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-17
 * @Project: Proof of Evolution
 * @Filename: arithmetic_operators.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-09
 * @Copyright: 2020
 */



package op


import(
  "fmt"
)

// Performs an addition and updates the State.
func (self *State)Add(a,b any)any{
  self.IncOperations(self.coeff["+"]+self.off["+"])
  var t string = fmt.Sprintf("%T", a)
  switch t {
    case "int": return a.(int)+b.(int)
    case "int8": return a.(int8)+b.(int8)
    case "int16": return a.(int16)+b.(int16)
    case "int32": return a.(int32)+b.(int32)
    case "int64": return a.(int64)+b.(int64)
    case "uint": return a.(uint)+b.(uint)
    case "uint8": return a.(uint8)+b.(uint8)
    case "uint16": return a.(uint16)+b.(uint16)
    case "uint32": return a.(uint32)+b.(uint32)
    case "uint64": return a.(uint64)+b.(uint64)
    case "float32": return a.(float32)+b.(float32)
    case "float64": return a.(float64)+b.(float64)
    case "complex64": return a.(complex64)+b.(complex64)
    case "complex128": return a.(complex128)+b.(complex128)
    default: fmt.Println("Invalid type")
    }
  return nil
}

// Performs a subtraction and updates the State.
func (self *State)Sub(a,b any)any{
  self.IncOperations(self.coeff["-"]+self.off["-"])
  var t string = fmt.Sprintf("%T", a)
  switch t {
    case "int": return a.(int)-b.(int)
    case "int8": return a.(int8)-b.(int8)
    case "int16": return a.(int16)-b.(int16)
    case "int32": return a.(int32)-b.(int32)
    case "int64": return a.(int64)-b.(int64)
    case "uint": return a.(uint)-b.(uint)
    case "uint8": return a.(uint8)-b.(uint8)
    case "uint16": return a.(uint16)-b.(uint16)
    case "uint32": return a.(uint32)-b.(uint32)
    case "uint64": return a.(uint64)-b.(uint64)
    case "float32": return a.(float32)-b.(float32)
    case "float64": return a.(float64)-b.(float64)
    case "complex64": return a.(complex64)-b.(complex64)
    case "complex128": return a.(complex128)-b.(complex128)
    default: fmt.Println("Invalid type")
    }
  return nil
}

// Performs a multiplication and updates the State.
func (self *State)Mul(a,b any)any{
  self.IncOperations(self.coeff["*"]+self.off["*"])
  var t string = fmt.Sprintf("%T", a)
  switch t {
    case "int": return a.(int)*b.(int)
    case "int8": return a.(int8)*b.(int8)
    case "int16": return a.(int16)*b.(int16)
    case "int32": return a.(int32)*b.(int32)
    case "int64": return a.(int64)*b.(int64)
    case "uint": return a.(uint)*b.(uint)
    case "uint8": return a.(uint8)*b.(uint8)
    case "uint16": return a.(uint16)*b.(uint16)
    case "uint32": return a.(uint32)*b.(uint32)
    case "uint64": return a.(uint64)*b.(uint64)
    case "float32": return a.(float32)*b.(float32)
    case "float64": return a.(float64)*b.(float64)
    case "complex64": return a.(complex64)*b.(complex64)
    case "complex128": return a.(complex128)*b.(complex128)
    default: fmt.Println("Invalid type")
    }
  return nil
}

// Performs a division and updates the State.
func (self *State)Div(a,b any)any{
  self.IncOperations(self.coeff["/"]+self.off["/"])
  var t string = fmt.Sprintf("%T", a)
  switch t {
    case "int": return a.(int)/b.(int)
    case "int8": return a.(int8)/b.(int8)
    case "int16": return a.(int16)/b.(int16)
    case "int32": return a.(int32)/b.(int32)
    case "int64": return a.(int64)/b.(int64)
    case "uint": return a.(uint)/b.(uint)
    case "uint8": return a.(uint8)/b.(uint8)
    case "uint16": return a.(uint16)/b.(uint16)
    case "uint32": return a.(uint32)/b.(uint32)
    case "uint64": return a.(uint64)/b.(uint64)
    case "float32": return a.(float32)/b.(float32)
    case "float64": return a.(float64)/b.(float64)
    case "complex64": return a.(complex64)/b.(complex64)
    case "complex128": return a.(complex128)/b.(complex128)
    default: fmt.Println("Invalid type")
    }
  return nil
}

// Returns the remainder a % b and updates the State.
func (self *State)Mod(a,b any)any{
  self.IncOperations(self.coeff["%"]+self.off["%"])
  var t string = fmt.Sprintf("%T", a)
  switch t {
    case "int": return a.(int)%b.(int)
    case "int8": return a.(int8)%b.(int8)
    case "int16": return a.(int16)%b.(int16)
    case "int32": return a.(int32)%b.(int32)
    case "int64": return a.(int64)%b.(int64)
    case "uint": return a.(uint)%b.(uint)
    case "uint8": return a.(uint8)%b.(uint8)
    case "uint16": return a.(uint16)%b.(uint16)
    case "uint32": return a.(uint32)%b.(uint32)
    case "uint64": return a.(uint64)%b.(uint64)
    default: fmt.Println("Invalid type")
    }
  return nil
}

// Returns a number incremented by 1 and updates the State.
func (self *State)Succ(a any)any{
  self.IncOperations(self.coeff["++"]+self.off["++"])
  var t string = fmt.Sprintf("%T", a)
  switch t {
    case "int": return a.(int)+1
    case "int8": return a.(int8)+1
    case "int16": return a.(int16)+1
    case "int32": return a.(int32)+1
    case "int64": return a.(int64)+1
    case "uint": return a.(uint)+1
    case "uint8": return a.(uint8)+1
    case "uint16": return a.(uint16)+1
    case "uint32": return a.(uint32)+1
    case "uint64": return a.(uint64)+1
    case "float32": return a.(float32)+1
    case "float64": return a.(float64)+1
    case "complex64": return a.(complex64)+1
    case "complex128": return a.(complex128)+1
    default: fmt.Println("Invalid type")
    }
  return nil
}

// Returns a number decremented by 1 and updates the State.
func (self *State)Prec(a any)any{
  self.IncOperations(self.coeff["--"]+self.off["--"])
  var t string = fmt.Sprintf("%T", a)
  switch t {
    case "int": return a.(int)-1
    case "int8": return a.(int8)-1
    case "int16": return a.(int16)-1
    case "int32": return a.(int32)-1
    case "int64": return a.(int64)-1
    case "uint": return a.(uint)-1
    case "uint8": return a.(uint8)-1
    case "uint16": return a.(uint16)-1
    case "uint32": return a.(uint32)-1
    case "uint64": return a.(uint64)-1
    case "float32": return a.(float32)-1
    case "float64": return a.(float64)-1
    case "complex64": return a.(complex64)-1
    case "complex128": return a.(complex128)-1
    default: fmt.Println("Invalid type")
    }
  return nil
}
