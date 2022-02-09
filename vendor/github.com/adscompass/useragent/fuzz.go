// +build gofuzz

package useragent

func Fuzz(data []byte) int {
	res := Parse(data, &Data{})
	if res {
		return 0
	}
	return 1
}
