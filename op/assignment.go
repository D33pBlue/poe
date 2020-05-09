/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-17
 * @Project: Proof of Evolution
 * @Filename: assignment.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-09
 * @Copyright: 2020
 */



package op

// Assigns an int to a variable and updates the State.
func (self *State)SetInt(a *int,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(int)
}

// Assigns an int8 to a variable and updates the State.
func (self *State)SetInt8(a *int8,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(int8)
}

// Assigns an int16 to a variable and updates the State.
func (self *State)SetInt16(a *int16,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(int16)
}

// Assigns an int32 to a variable and updates the State.
func (self *State)SetInt32(a *int32,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(int32)
}

// Assigns an int64 to a variable and updates the State.
func (self *State)SetInt64(a *int64,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(int64)
}

// Assigns an uint to a variable and updates the State.
func (self *State)SetUint(a *uint,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(uint)
}

// Assigns an uint8 to a variable and updates the State.
func (self *State)SetUint8(a *uint8,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(uint8)
}

// Assigns an uint16 to a variable and updates the State.
func (self *State)SetUint16(a *uint16,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(uint16)
}

// Assigns an uint32 to a variable and updates the State.
func (self *State)SetUint32(a *uint32,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(uint32)
}

// Assigns an uint64 to a variable and updates the State.
func (self *State)SetUint64(a *uint64,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(uint64)
}

// Assigns a float32 to a variable and updates the State.
func (self *State)SetFloat32(a *float32,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(float32)
}

// Assigns a float64 to a variable and updates the State.
func (self *State)SetFloat64(a *float64,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(float64)
}

// Assigns a bool to a variable and updates the State.
func (self *State)SetBool(a *bool,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(bool)
}

// Assigns a string to a variable and updates the State.
func (self *State)SetString(a *string,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(string)
}
