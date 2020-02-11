package log

// HTTP describes an HTTP-request-log
type HTTP struct {
	Method  string
	Route   string
	Code    uint32
	Size    uint64
	Version string
}
