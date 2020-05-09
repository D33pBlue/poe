/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-01
 * @Project: Proof of Evolution
 * @Filename: file_exists.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-09
 * @Copyright: 2020
 */
package utils

import(
  "os"
)

// true <=> the file of [filename] path exists.
// That file must not be a folder.
func FileExists(filename string) bool {
  info, err := os.Stat(filename)
  if os.IsNotExist(err) {
     return false
  }
  return !info.IsDir()
}
