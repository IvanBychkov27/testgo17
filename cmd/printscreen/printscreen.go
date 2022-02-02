package main

import (
	"bytes"
	"fmt"
	"github.com/vova616/screenshot"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"os/exec"
)

func main() {
	//fileName := "cmd/printscreen/png/f01.png"
	//saveFilePS(fileName)

	//example()

	bufPS()

	fmt.Println("Done")
}

func saveFilePS(fileName string) {
	file, errC := os.Create(fileName)
	if errC != nil {
		log.Println(errC.Error())
		return
	} else {
		log.Printf("Создан файл %s\n", fileName)
	}
	defer file.Close()

	img, errPS := printScreen()
	if errPS != nil {
		fmt.Println("error print screen:", errPS.Error())
		return
	}

	err := png.Encode(file, img)
	if err != nil {
		fmt.Println("error write:", err)
		return
	}
	fmt.Println("saved file", fileName)
}

func printScreen() (image.Image, error) {
	img, err := screenshot.CaptureScreen() // *image.RGBA
	if err != nil {
		return nil, err
	}
	return image.Image(img), nil // can cast to image.Image, but not necessary
	//return img, nil // can cast to image.Image, but not necessary
}
func example() {
	const width, height = 256, 256
	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{R: uint8((x + y) & 255), G: uint8((x + y) << 1 & 255), B: uint8((x + y) << 2 & 255), A: 255})
		}
	}
	f, err := os.Create("cmd/printscreen/png/image.png")
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func bufPS() {
	var buf bytes.Buffer

	file := "import"
	path, err := exec.LookPath(file)
	if err != nil {
		log.Fatal("import not installed !")
	}
	fmt.Printf("import is available at %s\n", path)

	cmd := exec.Command(file, "-window", "root", "root2.png")

	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
}
