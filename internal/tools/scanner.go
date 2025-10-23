package tools

import (
	"math/rand"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/pawannn/netlite/pkg"
)

type ScanResult struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Open     bool   `json:"open"`
	Error    string `json:"error,omitempty"`
}

func isLocalHostIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	if ip.IsLoopback() {
		return true
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return false
	}
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			var addrIP net.IP
			switch v := a.(type) {
			case *net.IPNet:
				addrIP = v.IP
			case *net.IPAddr:
				addrIP = v.IP
			}
			if addrIP != nil && addrIP.Equal(ip) {
				return true
			}
		}
	}
	return false
}

func ScanRange(host string, start int, end int, concurrency int, progress chan<- int) ([]ScanResult, pkg.NetliteErr) {
	if start < pkg.START_RANGE {
		start = pkg.START_RANGE
	}
	if end > pkg.END_RANGE {
		end = pkg.END_RANGE
	}
	if end < start {
		return nil, pkg.NetliteErr{
			ClientMessage: "end port must be >= start port",
			Error:         nil,
		}
	}
	if concurrency <= 0 {
		concurrency = pkg.CONCURRENCY
	}

	total := end - start + 1
	if total <= 0 {
		return nil, pkg.NetliteErr{
			ClientMessage: "no ports to scan",
			Error:         nil,
		}
	}

	local := isLocalHostIP(host)

	var timeout time.Duration
	var retries int
	var perJobDelay time.Duration
	const maxConcurrencyGlobal = 1200
	const maxConcurrencyLocal = 400

	if local {
		timeout = 120 * time.Millisecond
		retries = 0
		perJobDelay = 0
		if concurrency > maxConcurrencyLocal {
			concurrency = maxConcurrencyLocal
		}
	} else {
		timeout = 900 * time.Millisecond
		retries = 1
		perJobDelay = 2 * time.Millisecond
		if concurrency > maxConcurrencyGlobal {
			concurrency = maxConcurrencyGlobal
		}
	}

	if concurrency > total {
		concurrency = total
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	jobs := make(chan int, total)
	results := make(chan ScanResult, total)

	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		for port := range jobs {
			if perJobDelay > 0 {
				time.Sleep(perJobDelay + time.Duration(rng.Intn(6))*time.Millisecond)
			}

			address := net.JoinHostPort(host, strconv.Itoa(port))
			var conn net.Conn
			var err error

			for attempt := 0; attempt <= retries; attempt++ {
				conn, err = net.DialTimeout("tcp", address, timeout)
				if err == nil {
					conn.Close()
					res := ScanResult{Port: port, Protocol: "tcp", Open: true}
					results <- res
					if progress != nil {
						select {
						case progress <- port:
						default:
						}
					}
					break
				}

				if attempt == retries {
					res := ScanResult{Port: port, Protocol: "tcp", Open: false, Error: err.Error()}
					results <- res
					if progress != nil {
						select {
						case progress <- port:
						default:
						}
					}
					break
				}
				backoff := 80*time.Millisecond*(time.Duration(attempt+1)) + time.Duration(rng.Intn(120))*time.Millisecond
				time.Sleep(backoff)
			}
		}
	}

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go worker()
	}

	go func() {
		for p := start; p <= end; p++ {
			jobs <- p
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	out := make([]ScanResult, 0, total)
	for r := range results {
		out = append(out, r)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Port < out[j].Port
	})
	return out, pkg.NoErr
}
