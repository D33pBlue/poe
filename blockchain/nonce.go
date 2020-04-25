/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: nonce.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-25
 * @Copyright: 2020
 */

package blockchain

import(
  "github.com/D33pBlue/poe/ga"
)

type Nonce struct{
  Solution ga.DNA
  Evaluation float64
  Complexity float64
}

// TODO: implenent the nonce with (sol,eval,complex)

type NonceNoJob struct{
  Value int
}
