/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-01
 * @Project: Proof of Evolution
 * @Filename: file_exists.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-01
 * @Copyright: 2020
 */
package utils

import(
  "os"
)

 func FileExists(filename string) bool {
     info, err := os.Stat(filename)
     if os.IsNotExist(err) {
         return false
     }
     return !info.IsDir()
 }
