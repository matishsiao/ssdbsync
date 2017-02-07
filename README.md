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
      {"mode":"all","hash":"Hash-Test","key":""},
      {"mode":"","hash":"Hash-Key","key":"test"},
      {"mode":"scan","hash":"","key":"","start":"Key-A","end":"Key-B","bypass":["Bypass"]},
      {"mode":"","hash":"","key":"Single-Key"},
      {"mode":"hscan","start":"Hash-A","end":"Hash-B","bypass":["Bypass","bypassB"]}
    ]
  }

```

## SSDBProxy Configuration

| Config Fields  | Description | 
| ------------- | ------------- |
| debug  | debug mode:true / false  |
| host  | ssdb host  |
| port  | ssdb port  |
| password  | if you use auth params,you can use it to control connection |
| list | sync data list |

### List Configuration
| Node Fields | Description | Struct |
| ------------- | ------------- | ------------- |
| mode   | all / scan / hscan , default:empty string | string |
| hash | ssdb hash table name | string |
| key | ssdb key name | string |
| start | if you using scan or hscan, you will need set it | string |
| end | if you using scan or hscan, you will need set it | string |
| bypass | bypass sync data | string array |

#### List Mode Configuration
| Mode | Description |
| ------------- | ------------- |
| all | get all hash data with sync |
| scan | scan K/V data with sync |
| hscan | scan hash data with sync |
| "" | default, it will check hash or key has set. if yes, will automatic sync |


#How to build

```
 go get github.com/matishsiao/ssdbsync/
 cd $GOPATH/github.com/matishsiao/ssdbsync
 go build
```
