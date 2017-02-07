package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/matishsiao/gossdb/ssdb"
)

var (
	srcDBClient *ssdb.Client
	outDBClient *ssdb.Client
)

func main() {
	log.Println("Notice: Please using SSDB-Proxy to run this tools,becase we using a lot SSDB-Proxy functions.")
	configs, err := loadConfigs("sync.json")
	if err != nil {
		panic(err)
	}

	srcDBClient, err = ssdb.Connect(configs.SrcDB.Host, configs.SrcDB.Port, configs.SrcDB.Password)
	if err != nil {
		log.Fatal("src db connection error:", err)
	}
	srcDBClient.UseZip(true)
	srcDBClient.KeepAlive()
	if configs.SrcDB.Host != configs.OutDB.Host || configs.SrcDB.Port != configs.OutDB.Port {
		outDBClient, err = ssdb.Connect(configs.OutDB.Host, configs.OutDB.Port, configs.OutDB.Password)
		if err != nil {
			log.Fatal("output db connection error:", err)
		}
		outDBClient.UseZip(true)
		outDBClient.KeepAlive()
	} else {
		outDBClient = srcDBClient
	}

	DataSync(configs)

	if outDBClient != srcDBClient {
		outDBClient.Close()
	}
	srcDBClient.Close()

}

func DataSync(configs Configs) {
	for k, cmd := range configs.List {
		switch strings.ToLower(cmd.Mode) {
		case "all":
			if cmd.Hash != "" {
				list, err := srcDBClient.HashGetAll(cmd.Hash)
				if err != nil {
					log.Printf("DataSync[%d]:Sync %s get failed. error:%v\n", k, cmd.Hash, err)
					continue
				}
				for hk, hv := range list {
					outDBClient.BatchAppend("hset", cmd.Hash, hk, hv)
				}
				_, err = outDBClient.Exec()
				if err != nil {
					log.Printf("DataSync[%d]:Sync %s write failed. error:%v\n", k, cmd.Hash, err)
					continue
				}
				log.Printf("DataSync[%d]:Sync %s successful.\n", k, cmd.Hash)
			} else {
				log.Printf("DataSync[%d]:Sync format incorrect. hash:%v\n", k, cmd.Hash)
			}
		case "hscan":
			if cmd.Start != "" && cmd.End != "" {
				list, err := srcDBClient.Do("hlist", cmd.Start, cmd.End, -1)
				if err != nil {
					log.Printf("DataSync[%d]:Sync %s->%s get failed. error:%v\n", k, cmd.Start, cmd.End, err)
					continue
				}
				if len(list) > 1 {
					list = list[1:]
				}
				for _, hash := range list {
					bypass := false
					for _, substr := range cmd.Bypass {
						if strings.Contains(hash, substr) {
							bypass = true
							break
						}
					}
					if bypass {
						log.Printf("DataSync[%d]:Sync scan %s bypass this.\n", k, hash)
						continue
					}

					log.Printf("DataSync[%d]:Sync scan %s get data.\n", k, hash)
					clist, err := srcDBClient.HashGetAll(hash)
					if err != nil {
						log.Printf("DataSync[%d]:Sync %s get failed. error:%v\n", k, hash, err)
						continue
					}
					for hk, hv := range clist {
						outDBClient.BatchAppend("hset", hash, hk, hv)
					}
					_, err = outDBClient.Exec()
					if err != nil {
						log.Printf("DataSync[%d]:Sync %s write failed. error:%v\n", k, hash, err)
						continue
					}
					log.Printf("DataSync[%d]:Sync %s successful.\n", k, hash)
				}
				log.Printf("DataSync[%d]:Sync %s->%s total:%d successful.\n", k, cmd.Start, cmd.End, len(list))
			} else {
				log.Printf("DataSync[%d]:Sync format incorrect. start:%v end:%v\n", k, cmd.Start, cmd.End)
			}
		case "scan":
			if cmd.Start != "" && cmd.End != "" {
				list, err := srcDBClient.Do("scan", cmd.Start, cmd.End, -1)
				if err != nil {
					log.Printf("DataSync[%d]:Sync %s->%s get failed. error:%v\n", k, cmd.Start, cmd.End, err)
					continue
				}
				if len(list) > 1 {
					list = list[1:]
				}
				for _, key := range list {
					bypass := false
					for _, substr := range cmd.Bypass {
						if strings.Contains(key, substr) {
							bypass = true
							break
						}
					}
					if bypass {
						log.Printf("DataSync[%d]:Sync scan %s bypass this.\n", k, key)
						continue
					}
					result, err := srcDBClient.Do("get", key)
					if err != nil {
						log.Printf("DataSync[%d]:Sync key:%s get failed. error:%v\n", k, key, err)
						continue
					}
					if len(result) == 2 && result[0] == "ok" {
						_, err := outDBClient.Do("set", key, result[1])
						if err != nil {
							log.Printf("DataSync[%d]:Sync key:%s write failed. error:%v\n", k, key, err)
							continue
						}
						log.Printf("DataSync[%d]:Sync key:%s successful.\n", k, key)
					} else {
						log.Printf("DataSync[%d]:Sync key:%s get result failed. result:%v\n", k, key, result)
					}
				}
			} else {
				log.Printf("DataSync[%d]:Sync format incorrect. start:%v end:%v\n", k, cmd.Start, cmd.End)
			}
		default:
			if cmd.Hash != "" && cmd.Key != "" {
				result, err := srcDBClient.Do("hget", cmd.Hash, cmd.Key)
				if err != nil {
					log.Printf("DataSync[%d]:Sync %s:%s get failed. error:%v\n", k, cmd.Hash, cmd.Key, err)
					continue
				}
				if len(result) == 2 && result[0] == "ok" {
					_, err := outDBClient.Do("hset", cmd.Hash, cmd.Key, result[1])
					if err != nil {
						log.Printf("DataSync[%d]:Sync %s:%s write failed. error:%v\n", k, cmd.Hash, cmd.Key, err)
						continue
					}
					log.Printf("DataSync[%d]:Sync %s:%s successful.\n", k, cmd.Hash, cmd.Key)
				} else {
					log.Printf("DataSync[%d]:Sync %s:%s get result failed. result:%v\n", k, cmd.Hash, cmd.Key, result)
				}
			} else {
				if cmd.Key != "" {
					result, err := srcDBClient.Do("get", cmd.Key)
					if err != nil {
						log.Printf("DataSync[%d]:Sync key:%s get failed. error:%v\n", k, cmd.Key, err)
						continue
					}
					if len(result) == 2 && result[0] == "ok" {
						_, err := outDBClient.Do("set", cmd.Key, result[1])
						if err != nil {
							log.Printf("DataSync[%d]:Sync key:%s write failed. error:%v\n", k, cmd.Key, err)
							continue
						}
						log.Printf("DataSync[%d]:Sync key:%s successful.\n", k, cmd.Key)
					} else {
						log.Printf("DataSync[%d]:Sync key:%s get result failed. result:%v\n", k, cmd.Key, result)
					}
				} else {
					log.Printf("DataSync[%d]:Sync format incorrect. hash:%v key:%s\n", k, cmd.Hash, cmd.Key)
				}
			}
		}
	}
}

func loadConfigs(configName string) (Configs, error) {
	file, e := ioutil.ReadFile(configName)
	if e != nil {
		log.Printf("Load config error: %v\n", e)
		os.Exit(1)
	}

	var config Configs
	err := json.Unmarshal(file, &config)
	if err != nil {
		log.Printf("Config load error:%v \n", err)
		return config, err
	}
	return config, nil
}
