/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-17
 * @Project: Proof of Evolution
 * @Filename: assignment.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-19
 * @Copyright: 2020
 */



package op


func (self *State)SetInt(a *int,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(int)
}

func (self *State)SetInt8(a *int8,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(int8)
}

func (self *State)SetInt16(a *int16,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(int16)
}

func (self *State)SetInt32(a *int32,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(int32)
}

func (self *State)SetInt64(a *int64,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(int64)
}

func (self *State)SetUint(a *uint,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(uint)
}

func (self *State)SetUint8(a *uint8,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(uint8)
}

func (self *State)SetUint16(a *uint16,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(uint16)
}

func (self *State)SetUint32(a *uint32,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(uint32)
}

func (self *State)SetUint64(a *uint64,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(uint64)
}

func (self *State)SetFloat32(a *float32,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(float32)
}

func (self *State)SetFloat64(a *float64,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(float64)
}

func (self *State)SetBool(a *bool,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(bool)
}

func (self *State)SetString(a *string,b any){
  self.IncOperations(self.coeff["assign"]+self.off["assign"])
  *a = b.(string)
}
