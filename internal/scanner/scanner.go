package scanner

import (
	"net"
	"sort"
	"strconv"
	"time"

	"github.com/pawannn/netlite/pkg"
)

type ScanResult struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Open     bool   `json:"open"`
	Error    string `json:"error,omitempty"`
}

func ScanRange(host string, start int, end int, concurrency int) ([]ScanResult, pkg.NetliteErr) {
	if start < 1 {
		start = 1
	}

	if end > 65535 {
		end = 65535
	}
	if end < start {
		return nil, pkg.NetliteErr{
			ClientMessage: "end port must be >= start port",
			Error:         nil,
		}
	}

	if concurrency <= 0 {
		concurrency = 200
	}

	total := end - start + 1
	timeout := 250 * time.Millisecond

	jobs := make(chan int, total)
	results := make(chan ScanResult, total)

	worker := func() {
		for port := range jobs {
			address := net.JoinHostPort(host, strconv.Itoa(port))
			conn, err := net.DialTimeout("tcp", address, timeout)
			if err != nil {
				results <- ScanResult{Port: port, Protocol: "tcp", Open: false, Error: err.Error()}
			} else {
				conn.Close()
				results <- ScanResult{Port: port, Protocol: "tcp", Open: true}
			}
		}
	}

	if concurrency > total {
		concurrency = total
	}
	for i := 0; i < concurrency; i++ {
		go worker()
	}

	go func() {
		for p := start; p <= end; p++ {
			jobs <- p
		}
		close(jobs)
	}()

	out := make([]ScanResult, 0, total)
	for range total {
		out = append(out, <-results)
	}

	close(results)

	sort.Slice(out, func(i, j int) bool {
		return out[i].Port < out[j].Port
	})

	return out, pkg.NoErr
}
