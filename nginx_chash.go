package main

import (
	"fmt"
	//"hash/crc32"
	"net"
	"sort"
	"strconv"
)

const (
	factor            = 160
	initCrc32  uint32 = 0xffffffff
	maxRetries        = 20
)

var (
	//crc32Table  = crc32.MakeTable(crc32.IEEE)
	littleEdian = 0 & int(1)
	crc32Table  = []uint32{
		0x00000000, 0x77073096, 0xee0e612c, 0x990951ba,
		0x076dc419, 0x706af48f, 0xe963a535, 0x9e6495a3,
		0x0edb8832, 0x79dcb8a4, 0xe0d5e91e, 0x97d2d988,
		0x09b64c2b, 0x7eb17cbd, 0xe7b82d07, 0x90bf1d91,
		0x1db71064, 0x6ab020f2, 0xf3b97148, 0x84be41de,
		0x1adad47d, 0x6ddde4eb, 0xf4d4b551, 0x83d385c7,
		0x136c9856, 0x646ba8c0, 0xfd62f97a, 0x8a65c9ec,
		0x14015c4f, 0x63066cd9, 0xfa0f3d63, 0x8d080df5,
		0x3b6e20c8, 0x4c69105e, 0xd56041e4, 0xa2677172,
		0x3c03e4d1, 0x4b04d447, 0xd20d85fd, 0xa50ab56b,
		0x35b5a8fa, 0x42b2986c, 0xdbbbc9d6, 0xacbcf940,
		0x32d86ce3, 0x45df5c75, 0xdcd60dcf, 0xabd13d59,
		0x26d930ac, 0x51de003a, 0xc8d75180, 0xbfd06116,
		0x21b4f4b5, 0x56b3c423, 0xcfba9599, 0xb8bda50f,
		0x2802b89e, 0x5f058808, 0xc60cd9b2, 0xb10be924,
		0x2f6f7c87, 0x58684c11, 0xc1611dab, 0xb6662d3d,
		0x76dc4190, 0x01db7106, 0x98d220bc, 0xefd5102a,
		0x71b18589, 0x06b6b51f, 0x9fbfe4a5, 0xe8b8d433,
		0x7807c9a2, 0x0f00f934, 0x9609a88e, 0xe10e9818,
		0x7f6a0dbb, 0x086d3d2d, 0x91646c97, 0xe6635c01,
		0x6b6b51f4, 0x1c6c6162, 0x856530d8, 0xf262004e,
		0x6c0695ed, 0x1b01a57b, 0x8208f4c1, 0xf50fc457,
		0x65b0d9c6, 0x12b7e950, 0x8bbeb8ea, 0xfcb9887c,
		0x62dd1ddf, 0x15da2d49, 0x8cd37cf3, 0xfbd44c65,
		0x4db26158, 0x3ab551ce, 0xa3bc0074, 0xd4bb30e2,
		0x4adfa541, 0x3dd895d7, 0xa4d1c46d, 0xd3d6f4fb,
		0x4369e96a, 0x346ed9fc, 0xad678846, 0xda60b8d0,
		0x44042d73, 0x33031de5, 0xaa0a4c5f, 0xdd0d7cc9,
		0x5005713c, 0x270241aa, 0xbe0b1010, 0xc90c2086,
		0x5768b525, 0x206f85b3, 0xb966d409, 0xce61e49f,
		0x5edef90e, 0x29d9c998, 0xb0d09822, 0xc7d7a8b4,
		0x59b33d17, 0x2eb40d81, 0xb7bd5c3b, 0xc0ba6cad,
		0xedb88320, 0x9abfb3b6, 0x03b6e20c, 0x74b1d29a,
		0xead54739, 0x9dd277af, 0x04db2615, 0x73dc1683,
		0xe3630b12, 0x94643b84, 0x0d6d6a3e, 0x7a6a5aa8,
		0xe40ecf0b, 0x9309ff9d, 0x0a00ae27, 0x7d079eb1,
		0xf00f9344, 0x8708a3d2, 0x1e01f268, 0x6906c2fe,
		0xf762575d, 0x806567cb, 0x196c3671, 0x6e6b06e7,
		0xfed41b76, 0x89d32be0, 0x10da7a5a, 0x67dd4acc,
		0xf9b9df6f, 0x8ebeeff9, 0x17b7be43, 0x60b08ed5,
		0xd6d6a3e8, 0xa1d1937e, 0x38d8c2c4, 0x4fdff252,
		0xd1bb67f1, 0xa6bc5767, 0x3fb506dd, 0x48b2364b,
		0xd80d2bda, 0xaf0a1b4c, 0x36034af6, 0x41047a60,
		0xdf60efc3, 0xa867df55, 0x316e8eef, 0x4669be79,
		0xcb61b38c, 0xbc66831a, 0x256fd2a0, 0x5268e236,
		0xcc0c7795, 0xbb0b4703, 0x220216b9, 0x5505262f,
		0xc5ba3bbe, 0xb2bd0b28, 0x2bb45a92, 0x5cb36a04,
		0xc2d7ffa7, 0xb5d0cf31, 0x2cd99e8b, 0x5bdeae1d,
		0x9b64c2b0, 0xec63f226, 0x756aa39c, 0x026d930a,
		0x9c0906a9, 0xeb0e363f, 0x72076785, 0x05005713,
		0x95bf4a82, 0xe2b87a14, 0x7bb12bae, 0x0cb61b38,
		0x92d28e9b, 0xe5d5be0d, 0x7cdcefb7, 0x0bdbdf21,
		0x86d3d2d4, 0xf1d4e242, 0x68ddb3f8, 0x1fda836e,
		0x81be16cd, 0xf6b9265b, 0x6fb077e1, 0x18b74777,
		0x88085ae6, 0xff0f6a70, 0x66063bca, 0x11010b5c,
		0x8f659eff, 0xf862ae69, 0x616bffd3, 0x166ccf45,
		0xa00ae278, 0xd70dd2ee, 0x4e048354, 0x3903b3c2,
		0xa7672661, 0xd06016f7, 0x4969474d, 0x3e6e77db,
		0xaed16a4a, 0xd9d65adc, 0x40df0b66, 0x37d83bf0,
		0xa9bcae53, 0xdebb9ec5, 0x47b2cf7f, 0x30b5ffe9,
		0xbdbdf21c, 0xcabac28a, 0x53b39330, 0x24b4a3a6,
		0xbad03605, 0xcdd70693, 0x54de5729, 0x23d967bf,
		0xb3667a2e, 0xc4614ab8, 0x5d681b02, 0x2a6f2b94,
		0xb40bbe37, 0xc30c8ea1, 0x5a05df1b, 0x2d02ef8d,
	}
)

func NgxCrc32Long(s string, length int) uint32 {
	bs := []byte(s)
	c := initCrc32
	c = NgxCrc32Update(c, bs, length)
	return NgxCrc32Final(c)
}

// init
func NgxCrc32Init(c uint32) uint32 {
	return initCrc32
}

// update Crc32
func NgxCrc32Update(c uint32, bs []byte, length int) uint32 {
	for i := 0; i < len(bs); i++ {
		c = crc32Table[(c^uint32(bs[i]))&0xff] ^ (c >> 8)
	}
	for i := len(bs); i < length; i++ {
		c = crc32Table[(c^0)&0xff] ^ (c >> 8)
	}
	return c
}

// final crc
func NgxCrc32Final(c uint32) uint32 {
	//
	return c ^ initCrc32
}

// prev_hash
func NgxCrc32Prev(c uint32, length int) []byte {
	//  if LITTLE_ENDIAN, update value, and return byte
	/*
	   union {
	       uint32_t                          value;
	       u_char                            byte[4];
	   } prev_hash;
	*/
	ph := make([]byte, length)
	for i := 0; i < length; i++ {
		ph[i] = byte((c >> (8 * uint32(i))) & 0xff)
	}
	return ph
}

type ServerConf struct {
	server   string
	weight   int
	minConns int
	maxConns int
}

// Nginx hash node
type NgxHashNode struct {
	server string // server name
	hash   uint32 // hash value
}

// Nginx hash node List
type NgxHashNodeList []NgxHashNode

// Nginx Upstream chash
type NgxUpstreamChash struct {
	ServerMap    map[string][]string // server to peers
	HashNodeList NgxHashNodeList     // consistent hash list/ring
}

// Implement sort.Sort
func (nl NgxHashNodeList) Len() int           { return len(nl) }
func (nl NgxHashNodeList) Less(i, j int) bool { return nl[i].hash < nl[j].hash }
func (nl NgxHashNodeList) Swap(i, j int)      { nl[i], nl[j] = nl[j], nl[i] }
func (nl NgxHashNodeList) Search(hashKey uint32) int {
	return sort.Search(len(nl), func(i int) bool {
		return nl[i].hash >= hashKey
	})
}

func (nl NgxHashNodeList) removeDuplicate() NgxHashNodeList {
	i, j := 0, 1
	for j < len(nl) {
		if nl[i] != nl[j] {
			i++
			nl[i] = nl[j]
		}
		j++
	}
	return nl[0 : i+1]
}

// Input:
//  1. server: server, include port, aka www.example.com:80
//  2. weight: weight=<weight>
// Example:
//  * server www.example:80 weight=3
//      * server = www.example.com:80
//      * weight = 3
// Process:
//  1. base_hash = 0xffffffff
//  2. ngx_update_hash(base_hash, host)
//  3. ngx_update_hash(base_hash, '\0')
//  4. ngx_update_hash(base_hash, port)
func ngxChash(server string, weight int) NgxHashNodeList {
	host, port, err := net.SplitHostPort(server)
	if err != nil {
		panic(err)
	}

	// base_hash
	bh := initCrc32
	// host
	hostBytes := []byte(host)
	bh = NgxCrc32Update(bh, hostBytes, len(hostBytes))
	bh = NgxCrc32Update(bh, []byte{}, 1) // '\0'
	portBytes := []byte(port)
	bh = NgxCrc32Update(bh, portBytes, len(portBytes))

	ngxhnList := make(NgxHashNodeList, factor*weight)
	// prev_hash
	var ph []byte
	ph = []byte{}
	var h uint32

	for i := 0; i < factor*weight; i++ {
		// hash
		h = bh
		// update with ph
		h = NgxCrc32Update(h, ph, 4)
		// final
		h = NgxCrc32Final(h)
		ngxhnList[i] = NgxHashNode{
			hash:   h,
			server: server,
		}
		// prev_hash
		ph = NgxCrc32Prev(h, 4)
	}
	return ngxhnList
}

func NgxUpstreamServerToPeers(scList []ServerConf) map[string][]string {
	spm := make(map[string][]string)
	for _, sc := range scList {
		// if sc.server is domain, do dns resove
		if host, port, err := net.SplitHostPort(sc.server); err == nil {
			if ip := net.ParseIP(host); ip != nil {
				spm[sc.server] = []string{sc.server}
				continue
			}
			//
			if addrs, err := net.LookupHost(host); err == nil {
				peers := make([]string, len(addrs))
				for idx, addr := range addrs {
					peers[idx] = net.JoinHostPort(addr, port)
				}
				spm[sc.server] = peers
			} else {
				// no avaliable peer
				//panic(err)
				fmt.Printf("%v\n", err)
				spm[sc.server] = []string{}
			}
		} else {
			// invalid server name
			panic(err)
		}
		spm[sc.server] = spm[sc.server]
	}
	return spm
}

func NgxUpstreamInitChash(scList []ServerConf) NgxUpstreamChash {
	var ngxhnList NgxHashNodeList

	for _, sc := range scList {
		hnList := ngxChash(sc.server, sc.weight)
		ngxhnList = append(ngxhnList, hnList...)
	}
	//showNgxHashNodeList(ngxhnList)

	sort.Sort(ngxhnList)
	ngxhnList = ngxhnList.removeDuplicate()
	return NgxUpstreamChash{
		ServerMap:    NgxUpstreamServerToPeers(scList),
		HashNodeList: ngxhnList,
	}
}

// Find Server
//  1. get hash from hash List, binarySearch, last match/not [l, r), return l
//  2. hash % len(hashList)
//  3. find peer under this server
//      * consider weight/conns/fails, best current_weight win
//      * O(len(all_upstream_peers))
//  4. hash++, retry from step 2, if no peer avaliable
//  5. retry times reach, no server avaliable
//      * < 1.13.6, retry times == len(hashList)
//      * >= 1.13.6, retry times == 20, than, use rr
//
func (nuc NgxUpstreamChash) NgxUpstreamChashGetPeer(key string) string {
	hashKey := NgxCrc32Long(key, len(key))
	hashPoint := nuc.HashNodeList.Search(hashKey)
	retries := 0
	for {
		server := nuc.HashNodeList[hashPoint%len(nuc.HashNodeList)].server

		// should traversing all peers, then found best match peer
		// what's best (match all):
		//  1. server match
		//  2. conns <= max_conns && fails <= max_fails
		//  3. if fails, must pass fail_timeout
		//  4. consider weight: max current_weight
		if peers, found := nuc.ServerMap[server]; found && len(peers) > 0 {
			fmt.Printf("Input Key: %s, Match server: %s, Match peers: %v\n", key, server, peers)
			//  1. should base on conns/weight
			return peers[0]
		} else {
			retries++
			if retries > maxRetries {
				fmt.Printf("MaxRetiries: %d, fallback into rr algorithm.\n", maxRetries)

				// fallback into rr algorithm
				for _, peers = range nuc.ServerMap {
					// 1. should base on conns/weight
					if len(peers) > 0 {
						return peers[0]
					}
				}
				return ""
			}
			hashKey++
		}
	}
	return ""
}

func showNgxHashNodeList(nl NgxHashNodeList) {
	l := strconv.Itoa(len(strconv.Itoa(len(nl))))
	for idx, n := range nl {
		fmt.Printf("index: %"+l+"d\thash: %10d\tserver: %s\n", idx, n.hash, n.server)
	}
}

func main() {
	scList := []ServerConf{
		ServerConf{
			server: "127.0.0.1:8080",
			weight: 1,
		},
		ServerConf{
			server: "example.com:443",
			weight: 2,
		},
	}
	nuc := NgxUpstreamInitChash(scList)
	showNgxHashNodeList(nuc.HashNodeList)
	key := "192.168.0.1"
	fmt.Printf("Input key: %s, Match Peer: %s\n", key, nuc.NgxUpstreamChashGetPeer(key))
}
