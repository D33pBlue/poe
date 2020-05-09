/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-17
 * @Project: Proof of Evolution
 * @Filename: math_functions.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-09
 * @Copyright: 2020
 */



package op

import(
  "fmt"
  "math"
)

// Wrapper for a unary operation.
func wrap1(a any,f func(float64)float64)any{
  var t string = fmt.Sprintf("%T", a)
  switch t {
    case "int": return f(float64(a.(int)))
    case "int8": return f(float64(a.(int8)))
    case "int16": return f(float64(a.(int16)))
    case "int32": return f(float64(a.(int32)))
    case "int64": return f(float64(a.(int64)))
    case "uint": return f(float64(a.(uint)))
    case "uint8": return f(float64(a.(uint8)))
    case "uint16": return f(float64(a.(uint16)))
    case "uint32": return f(float64(a.(uint32)))
    case "uint64": return f(float64(a.(uint64)))
    case "float32": return f(float64(a.(float32)))
    case "float64": return f(a.(float64))
    default: fmt.Println("Invalid type")
    }
  return nil
}

// Wrapper for a binary operation.
func wrap2(a,b any,f func(float64,float64)float64)any{
  var t string = fmt.Sprintf("%T", a)
  var t2 string = fmt.Sprintf("%T", b)
  var b2 float64
  switch t2 {
    case "int": b2 = float64(b.(int))
    case "int8": b2 = float64(b.(int8))
    case "int16": b2 = float64(b.(int16))
    case "int32": b2 = float64(b.(int32))
    case "int64": b2 = float64(b.(int64))
    case "uint": b2 = float64(b.(uint))
    case "uint8": b2 = float64(b.(uint8))
    case "uint16": b2 = float64(b.(uint16))
    case "uint32": b2 = float64(b.(uint32))
    case "uint64": b2 = float64(b.(uint64))
    case "float32": b2 = float64(b.(float32))
    case "float64": b2 = b.(float64)
    default: fmt.Println("Invalid type")
  }
  switch t {
    case "int": return f(float64(a.(int)),b2)
    case "int8": return f(float64(a.(int8)),b2)
    case "int16": return f(float64(a.(int16)),b2)
    case "int32": return f(float64(a.(int32)),b2)
    case "int64": return f(float64(a.(int64)),b2)
    case "uint": return f(float64(a.(uint)),b2)
    case "uint8": return f(float64(a.(uint8)),b2)
    case "uint16": return f(float64(a.(uint16)),b2)
    case "uint32": return f(float64(a.(uint32)),b2)
    case "uint64": return f(float64(a.(uint64)),b2)
    case "float32": return f(float64(a.(float32)),b2)
    case "float64": return f(a.(float64),b2)
    default: fmt.Println("Invalid type")
    }
  return nil
}

// Returns the absolute value of a numner and updates the State.
func (self *State)Abs(a any)any{
  self.IncOperations(self.coeff["abs"]+self.off["abs"])
  return wrap1(a,math.Abs)
}

// Returns the power a^b and updates the State.
func (self *State)Pow(a,b any)any{
  self.IncOperations(self.coeff["pow"]+self.off["pow"])
  return wrap2(a,b,math.Pow)
}

// Returns the square root of a and updates the State.
func (self *State)Sqrt(a any)any{
  self.IncOperations(self.coeff["sqrt"]+self.off["sqrt"])
  return wrap1(a,math.Sqrt)
}

// Returns the smaller integer bigger than a and updates the State.
func (self *State)Ceil(a any)any{
  self.IncOperations(self.coeff["ceil"]+self.off["ceil"])
  return wrap1(a,math.Ceil)
}

// Returns the biggest integer smaller than a and updates the State.
func (self *State)Floor(a any)any{
  self.IncOperations(self.coeff["floor"]+self.off["floor"])
  return wrap1(a,math.Floor)
}

// Rounds a to the nearest integer and updates the State.
func (self *State)Round(a any)any{
  self.IncOperations(self.coeff["round"]+self.off["round"])
  return wrap1(a,math.Round)
}

// Returns the min of a and b and updates the State.
func (self *State)Min(a,b any)any{
  self.IncOperations(self.coeff["min"]+self.off["min"])
  return wrap2(a,b,math.Min)
}

// Returns the max of a and b and updates the State.
func (self *State)Max(a,b any)any{
  self.IncOperations(self.coeff["max"]+self.off["max"])
  return wrap2(a,b,math.Max)
}

// Returns the sin of a and updates the State.
func (self *State)Sin(a any)any{
  self.IncOperations(self.coeff["sin"]+self.off["sin"])
  return wrap1(a,math.Sin)
}

// Returns the cos of a and updates the State.
func (self *State)Cos(a any)any{
  self.IncOperations(self.coeff["cos"]+self.off["cos"])
  return wrap1(a,math.Cos)
}

// Returns the asin of a and updates the State.
func (self *State)Asin(a any)any{
  self.IncOperations(self.coeff["asin"]+self.off["asin"])
  return wrap1(a,math.Asin)
}

// Returns the acos of a and updates the State.
func (self *State)Acos(a any)any{
  self.IncOperations(self.coeff["acos"]+self.off["acos"])
  return wrap1(a,math.Acos)
}

// Returns the tan of a and updates the State.
func (self *State)Tan(a any)any{
  self.IncOperations(self.coeff["tan"]+self.off["tan"])
  return wrap1(a,math.Tan)
}

// Returns the atan of a and updates the State.
func (self *State)Atan(a any)any{
  self.IncOperations(self.coeff["atan"]+self.off["atan"])
  return wrap1(a,math.Atan)
}

// Returns the sinh of a and updates the State.
func (self *State)Sinh(a any)any{
  self.IncOperations(self.coeff["sinh"]+self.off["sinh"])
  return wrap1(a,math.Sinh)
}

// Returns the cosh of a and updates the State.
func (self *State)Cosh(a any)any{
  self.IncOperations(self.coeff["cosh"]+self.off["cosh"])
  return wrap1(a,math.Cosh)
}

// Returns the asinh of a and updates the State.
func (self *State)Asinh(a any)any{
  self.IncOperations(self.coeff["asinh"]+self.off["asinh"])
  return wrap1(a,math.Asinh)
}

// Returns the acosh of a and updates the State.
func (self *State)Acosh(a any)any{
  self.IncOperations(self.coeff["acosh"]+self.off["acosh"])
  return wrap1(a,math.Acosh)
}

// Returns the tanh of a and updates the State.
func (self *State)Tanh(a any)any{
  self.IncOperations(self.coeff["tanh"]+self.off["tanh"])
  return wrap1(a,math.Tanh)
}

// Returns the atanh of a and updates the State.
func (self *State)Atanh(a any)any{
  self.IncOperations(self.coeff["atanh"]+self.off["atanh"])
  return wrap1(a,math.Atanh)
}

// Returns the natural logarithm of a and updates the State.
func (self *State)Log(a any)any{
  self.IncOperations(self.coeff["log"]+self.off["log"])
  return wrap1(a,math.Log)
}

// Returns the logarithm on base 2 of a and updates the State.
func (self *State)Log2(a any)any{
  self.IncOperations(self.coeff["log2"]+self.off["log2"])
  return wrap1(a,math.Log2)
}

// Returns the logarithm on base 10 of a and updates the State.
func (self *State)Log10(a any)any{
  self.IncOperations(self.coeff["log10"]+self.off["log10"])
  return wrap1(a,math.Log10)
}

// Returns exp(a) and updates the State.
func (self *State)Exp(a any)any{
  self.IncOperations(self.coeff["exp"]+self.off["exp"])
  return wrap1(a,math.Exp)
}

// Returns the power 2^a and updates the State.
func (self *State)Exp2(a any)any{
  self.IncOperations(self.coeff["exp2"]+self.off["exp2"])
  return wrap1(a,math.Exp2)
}
