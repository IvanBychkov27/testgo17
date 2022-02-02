// https://pkg.go.dev/embed
// Встраивание пакета обеспечивает доступ к файлам, встроенным в запущенную программу Go
// https://nikgalushko.com/post/go_1.16_embed/

package main

import (
	"embed"
	//_ "embed"
	"fmt"
)

//go:embed hello.txt
var s []byte

//go:embed www
var fs embed.FS

func main() {
	fmt.Println(string(s))
	fmt.Println()

	files, err := fs.ReadDir("www")
	if err != nil {
		fmt.Println("error read dir", err.Error())
	}

	for _, f := range files {
		fmt.Println(f.Name())
		data, err := fs.ReadFile("www/" + f.Name())
		if err != nil {
			fmt.Println("read file error", err.Error())
		}
		fmt.Println(string(data))
		fmt.Println()
	}

	fmt.Println("Done...")

}
