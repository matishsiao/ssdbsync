# ssdbsync

ssdb-sync is ssdb dump tool for real time sync data.

# Warning

Please make sure you using ssdbproxy for your ssdb.

https://github.com/matishsiao/ssdbproxy

# Version

version: 0.0.1

# Futures
    support functions:
	single key sync
	single hash key sync
	key scan sync
	hash scan sync

# Configuration

use json format to configuration ssdb-sync setting.

## Configuration Example

```
	{
    "srcdb":{// source ssdb
      "host":"192.168.0.1",
      "port":4002,
      "password":""
    },
    "outdb":{ //output ssdb
      "host":"192.168.1.1",
      "port":4002,
      "password":""
    },
    "list":[ //need sync data list
      {"type":"all","mode":"diff","hash":"Hash-Test","key":""},
      {"type":"","mode":"","hash":"Hash-Key","key":"test"},
      {"type":"scan","mode":"diff","hash":"","key":"","start":"Key-A","end":"Key-B","bypass":["Bypass"]},
      {"type":"","mode":"","hash":"","key":"Single-Key"},
      {"type":"hscan","mode":"diff","start":"Hash-A","end":"Hash-B","bypass":["Bypass","bypassB"]}
    ]
  }

```

## SSDB-Sync Configuration

| Config Fields  | Description |
| ------------- | ------------- |
| srcdb  | source ssdb  |
| outdb  | output ssdb |
| list | sync data list |

### SSDB Configuration
| Node Fields | Description | Struct |
| ------------- | ------------- | ------------- |
| host   | ssdb host | string |
| port   | ssdb port | int |
| password   | ssdb auth string | string |

### List Configuration
| Node Fields | Description | Struct |
| ------------- | ------------- | ------------- |
| type   | all / scan / hscan , default:empty string | string |
| mode   | diff / "" , default:empty string | string |
| hash | ssdb hash table name | string |
| key | ssdb key name | string |
| start | if you using scan or hscan, you will need set it | string |
| end | if you using scan or hscan, you will need set it | string |
| bypass | bypass sync data | string array |

#### List Type Configuration
| Type | Description |
| ------------- | ------------- |
| all | get all hash data with sync |
| scan | scan K/V data with sync |
| hscan | scan hash data with sync |
| "" | default, it will check hash or key has set. if yes, will automatic sync |

#### List Mode Configuration
| Mode | Description |
| ------------- | ------------- |
| diff | diff all K/V, if value had changed,will automatic sync the value to output SSDB |
| "" | default, it will sync all data to output SSDB |


#How to build

```
 go get github.com/matishsiao/ssdbsync/
 cd $GOPATH/github.com/matishsiao/ssdbsync
 go build
```
