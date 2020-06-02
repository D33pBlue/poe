/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-10
 * @Project: Proof of Evolution
 * @Filename: downloader.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-29
 * @Copyright: 2020
 */


package utils

import(
  "os"
  "io"
  "fmt"
  "net/http"
  "io/ioutil"
)

// Downloads a url to a local file, writing as it downloads and not
// loading the whole file into memory.
func DownloadFile(url string, filepath string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// Returns the data at a url as []byte
func FetchDataFromUrl(url string)[]byte{
	resp, err := http.Get(url)
	if err != nil {
    fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
  data, err2 := ioutil.ReadAll(resp.Body)
  if err2 != nil{
    fmt.Println(err2)
    return nil
  }
  return data
}
