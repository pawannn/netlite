package pkg

type NetliteErr struct {
	ClientMessage string
	Error         error
}

var NoErr NetliteErr = NetliteErr{}
