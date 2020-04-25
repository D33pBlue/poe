/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: state_test.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-25
 * @Copyright: 2020
 */



package op

import(
  "testing"
  "github.com/D33pBlue/poe/utils"
)

func TestOps(t *testing.T){
  hb := new(utils.HashBuilder)
  hb.Add("data..")
  st := MakeState(hb.GetHash())
  var a int = 5
  var b float64 = 10.5
  var c float64 = 57.66
  if st.Add(b,c).(float64)!=68.16{
    t.Errorf("%f != %f\n",st.Add(b,c).(float64),68.16)
  }
  if st.Add(a,4)!=9 {t.Errorf("%v != %v\n",st.Add(a,4)!=9,9)}
  if st.Eq(9,9)!=true {t.Errorf("%v != %v\n",st.Eq(9,9),true)}
  if st.Eq(9.1,9.10)!=true {t.Errorf("%v != %v\n",st.Eq(9.1,9.10),true)}
  if st.Eq(9.1,9.11)!=false {t.Errorf("%v != %v\n",st.Eq(9.1,9.11),false)}
  if st.Eq(9,-9)!=false {t.Errorf("%v != %v\n",st.Eq(9,-9),false)}
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
  if b!=57.66 {t.Errorf("%v != %v\n",b,57.66)}
  if c!=60.66 {t.Errorf("%v != %v\n",c,60.66)}
  ltrue := []int{2,5,15,25,35}
  if len(lk)!=len(ltrue) {
    t.Errorf("%v != %v\n",lk,ltrue)
  }else{
    diff := false
    for i:=0;!diff && i<len(lk);i++{
      if lk[i]!=ltrue[i]{
        diff = true
        t.Errorf("%v != %v\n",lk,ltrue)
      }
    }
  }
}
