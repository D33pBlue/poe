/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-16
 * @Project: Proof of Evolution
 * @Filename: packet.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-19
 * @Copyright: 2020
 */



package ga

// The type of the messages the miner exchanges with
// the GA thread.
type Packet struct{
  Solution Sol
  End bool
  Shared Population
}
