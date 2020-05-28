/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: transaction.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-25
 * @Copyright: 2020
 */

package blockchain

import(
  "time"
  "github.com/D33pBlue/poe/utils"
)


// Transaction is the iterface for all the transactions. Assumptions:
// - all inputs from the same address (that signs the transactions)
// - only one output for each address that receives money
//
// Check has to check the validity of the Transaction inside the chain.
// GetHash calculates and returns the hash of the Transaction.
// GetHashCached returns the cached hash of the Transaction.
// GetSignature returns the signature of the Transaction.
// GetCreator returns the public key of the creator of the Transaction.
// GetTimestamp returns the time of creation of the Transaction.
// GetType returns the string type of the Transaction.
// GetOutputAt returns the TrOutput at an index, inside the Transaction.
type Transaction interface{
  Check(block *Block,trChanges *map[string]string)bool
  GetHash()string
  GetHashCached()string
  GetCreator()utils.Addr
  GetTimestamp()time.Time
  GetType()string
  GetOutputAt(int)*TrOutput
  Serialize()[]byte
}

// The possible types of a Transaction
const (
  TrStd = "StdTransaction"
  TrCoin = "CoinTransaction"
  TrJob = "JobTransaction"
  TrSol = "SolTransaction"
  TrRes = "ResTransaction"
  TrPrize = "PrizeTransaction"
)
