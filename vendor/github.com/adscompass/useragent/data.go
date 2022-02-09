package useragent

import "sync"

// Data struct represent parsed data
type Data struct {
	BrowserName     []byte
	BrowserVersion  []byte
	PlatformName    []byte
	PlatformVersion []byte
	IsMobile        bool
	IsCrawler       bool
	IsWebView       bool
}

var (
	pool = sync.Pool{
		New: func() interface{} {
			return newData()
		},
	}
)

func newData() *Data {
	return &Data{
		BrowserName:     make([]byte, 0, 64),
		BrowserVersion:  make([]byte, 0, 64),
		PlatformName:    make([]byte, 0, 64),
		PlatformVersion: make([]byte, 0, 64),
		IsMobile:        false,
		IsCrawler:       false,
		IsWebView:       false,
	}
}

// AcquireData returns a Data struct from the pool
func AcquireData() *Data {
	return pool.Get().(*Data)
}

// ReleaseData receive a data struct, reset it and put to the pool
func ReleaseData(d *Data) {
	d.reset()
	pool.Put(d)
}

func (d *Data) reset() {
	d.BrowserName = d.BrowserName[:0]
	d.BrowserVersion = d.BrowserVersion[:0]
	d.PlatformName = d.PlatformName[:0]
	d.PlatformVersion = d.PlatformVersion[:0]
	d.IsMobile = false
	d.IsCrawler = false
	d.IsWebView = false
}
