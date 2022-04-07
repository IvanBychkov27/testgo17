package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"sync"
)

var (
	metricPool = &sync.Pool{
		New: func() interface{} {
			return &Metric{ID: make([]byte, 32)}
		},
	}
)

func AcquireMetric() *Metric {
	return metricPool.Get().(*Metric)
}

func ReleaseMetric(m *Metric) {
	m.reset()
	metricPool.Put(m)
}

func (m *Metric) reset() {
	m.ID = m.ID[:0]
}

type Metric struct {
	ID []byte
}

func main() {
	run()
	aaa()
	fmt.Println("done...")

}

func run() {
	n := 1000
	for n > 0 {
		n--
		d := convertIntInStringInByte(1234567890)
		//d, err := code(1234567890)
		//if err != nil {
		//	return
		//}
		m := Metric{ID: d}
		m.ID = m.ID[:0]
	}
}

func convertIntInStringInByte(i int) []byte {
	return []byte(strconv.Itoa(i))
}

func aaa() {
	var b []byte
	b = append(b, '1')
	b = append(b, '2')
	fmt.Println("b =", b)
}

func code(d int) ([]byte, error) {
	buf := new(bytes.Buffer)
	//err := binary.Write(buf, binary.LittleEndian, float64(d))
	err := binary.Write(buf, binary.BigEndian, uint32(d))
	if err != nil {
		return nil, fmt.Errorf("binary.Write failed: %w", err)
	}
	return buf.Bytes(), nil
}

func encode(d []byte) (int, error) {
	//var res float64
	var res uint32
	b := bytes.NewReader(d)
	//err := binary.Read(b, binary.LittleEndian, &res)
	err := binary.Read(b, binary.BigEndian, &res)
	if err != nil {
		return 0, fmt.Errorf("binary.Read failed: %w", err)
	}
	return int(res), nil
}
