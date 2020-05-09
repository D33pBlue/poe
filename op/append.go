/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-17
 * @Project: Proof of Evolution
 * @Filename: append.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-09
 * @Copyright: 2020
 */

package op

// Append an int to a []int and update the State.
func (self *State)AppInt(a *[]int,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(int))
}

// Append an int8 to a []int8 and update the State.
func (self *State)AppInt8(a *[]int8,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(int8))
}

// Append an int16 to a []int16 and update the State.
func (self *State)AppInt16(a *[]int16,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(int16))
}

// Append an int32 to a []int32 and update the State.
func (self *State)AppInt32(a *[]int32,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(int32))
}

// Append an int64 to a []int64 and update the State.
func (self *State)AppInt64(a *[]int64,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(int64))
}

// Append an uint to a []uint and update the State.
func (self *State)AppUint(a *[]uint,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(uint))
}

// Append an uint8 to a []uint8 and update the State.
func (self *State)AppUint8(a *[]uint8,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(uint8))
}

// Append an uint16 to a []uint16 and update the State.
func (self *State)AppUint16(a *[]uint16,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(uint16))
}

// Append an uint32 to a []uint32 and update the State.
func (self *State)AppUint32(a *[]uint32,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(uint32))
}

// Append an uint64 to a []uint64 and update the State.
func (self *State)AppUint64(a *[]uint64,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(uint64))
}

// Append a float32 to a []float32 and update the State.
func (self *State)AppFloat32(a *[]float32,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(float32))
}

// Append a float64 to a []float64 and update the State.
func (self *State)AppFloat64(a *[]float64,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(float64))
}

// Append a bool to a []bool and update the State.
func (self *State)AppBool(a *[]bool,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(bool))
}

// Append a string to a []string and update the State.
func (self *State)AppString(a *[]string,b any){
  self.IncOperations(self.coeff["append"]+self.off["append"])
  *a = append(*a,b.(string))
}
