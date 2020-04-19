/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: currency.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-19
 * @Copyright: 2020
 */



// Package coin provides the type definition of the currency.
package coin

type any interface{}

// Currency is the interface to implement
// to define a currency. This interface is
// defined in order to make easier to change
// how the currency is defined (e.g. with float or with int).
type Currency interface{
  New(any)Currency
  String() string
  Check()bool
  Compare(b Currency) int
  Add(b Currency)Currency
  Sub(b Currency)Currency
  Mul(t float64)Currency
  Div(t float64)Currency
}

// coin is used to choose the default Currency implementation.
var coin SpC

// New function returns a Currency object with v as value
func New(v any)Currency{
  return coin.New(v)
}
