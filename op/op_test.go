package op

import(
  "fmt"
  "testing"
)


func TestOps(t *testing.T){
  st := MakeState("")
  var a int = 5
  var b float64 = 10.5
  var c float64 = 57.66
  if st.Add(b,c).(float64)!=68.16{
    t.Errorf("%f != %f\n",st.Add(b,c).(float64),68.16)
  }
  fmt.Println(st.Add(a,4))
  fmt.Println(st.Eq(9,9))
  fmt.Println(st.Eq(9.1,9.10))
  fmt.Println(st.Eq(9.1,9.11))
  fmt.Println(st.Eq(9,-9))
  fmt.Println(st.Eq("9","9"))
  var res bool = st.And(false,st.Eq(2.0,st.Div(12.0,st.Add(1.0,5.0))))
  if res != false{
    t.Errorf("%v != %v\n",res,false)
  }
  st.SetFloat64(&b,c)
  var lk []int
  c += 3
  lk = append(lk,2)
  for i:=0;i<4;i++{
    st.AppInt(&lk,a)
    st.SetInt(&a,st.Add(a,10))
  }
  fmt.Println(b,c,st.NumOperations())
  fmt.Println(lk)
}
