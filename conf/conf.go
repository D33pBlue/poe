/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-10
 * @Project: Proof of Evolution
 * @Filename: conf.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-28
 * @Copyright: 2020
 */

package conf

import(
  // "fmt"
  "errors"
  "io/ioutil"
  "strconv"
  "encoding/json"
  "github.com/D33pBlue/poe/utils"
)

type Config struct{
  MainDataFolder string
  KeyFolder string
  JobFolder string
  ChainFolder string
  Key string
  Port string
  Miners []string
  publicKey utils.Addr
}


func LoadConfiguration(file string)(*Config,error) {
  // Load and parse config file
  config := new(Config)
  dat, err := ioutil.ReadFile(file)
  if err != nil { return nil,err }
  err = json.Unmarshal([]byte(dat), config)
  if err != nil { return nil,err }
  // check and load public key
  err2 := config.loadPublicKey(config.MainDataFolder+config.KeyFolder+config.Key)
  if err2 != nil { return nil,err2 }
  // check numeric port
  _,err3 := strconv.Atoi(config.Port)
  if err3!=nil{ return nil,err3 }
  return config,nil
}

func (self *Config)GetKeyPath()string{
  path := self.MainDataFolder+self.KeyFolder+self.Key
  if !utils.FileExists(path){
    return ""
  }
  if !utils.FileExists(path+".priv"){
    return ""
  }
  return path
}

// Returns two paths, one to store the job and one to store its data.
func (self *Config)GetSuitablePathForJob(hash string)(string,string){
  jobPath := self.MainDataFolder+self.JobFolder+hash+".go"
  dataPath := self.MainDataFolder+self.JobFolder+hash+"_data.json"
  return jobPath,dataPath
}


func (self *Config)GetSuitablePathForResults(jobind string)string{
  return self.MainDataFolder+self.JobFolder+"Results_"+jobind+".json"
}


func (self *Config)GetPrivateKey()utils.Key{
  path := self.GetKeyPath()
  data, err := ioutil.ReadFile(path)
  if err!=nil{ return nil }
  pub,err2 := utils.LoadPublicKeyFromPemStr(data)
  if err2!=nil { return nil}
  data, err = ioutil.ReadFile(path+".priv")
  if err!=nil { return nil}
  priv,err3 := utils.LoadPrivateKeyFromPemStr(data)
  if err3!=nil { return nil }
  priv.PublicKey = *pub
  return priv
}

func (self *Config)GetPublicKey()utils.Addr{
  return self.publicKey
}

func (self *Config)GetPort()string{
  return self.Port
}

func (self *Config)GetChainFolder()string{
  return self.MainDataFolder+self.ChainFolder
}

func (self *Config)GetLinkedMinersIp()[]string{
  return self.Miners
}

func (self *Config)loadPublicKey(keypath string)error{
  if utils.FileExists(keypath){
    data, err := ioutil.ReadFile(keypath)
    if err!=nil{ return err }
    pub,err2 := utils.LoadPublicKeyFromPemStr(data)
    if err2!=nil { return err2 }
    self.publicKey = utils.GetAddr2(pub)
  }else{
    return errors.New(keypath+" does not exist\nYou need to link a valid public key file to start mining.\nYou can generate it with a client.")
  }
  return nil
}
