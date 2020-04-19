/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-17
 * @Project: Proof of Evolution
 * @Filename: append.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-19
 * @Copyright: 2020
 */



package op


func (self *State)AppInt(a *[]int,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(int))
}

func (self *State)AppInt8(a *[]int8,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(int8))
}

func (self *State)AppInt16(a *[]int16,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(int16))
}

func (self *State)AppInt32(a *[]int32,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(int32))
}

func (self *State)AppInt64(a *[]int64,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(int64))
}

func (self *State)AppUint(a *[]uint,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(uint))
}

func (self *State)AppUint8(a *[]uint8,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(uint8))
}

func (self *State)AppUint16(a *[]uint16,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(uint16))
}

func (self *State)AppUint32(a *[]uint32,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(uint32))
}

func (self *State)AppUint64(a *[]uint64,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(uint64))
}

func (self *State)AppFloat32(a *[]float32,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(float32))
}

func (self *State)AppFloat64(a *[]float64,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(float64))
}

func (self *State)AppBool(a *[]bool,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(bool))
}

func (self *State)AppString(a *[]string,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(string))
}
