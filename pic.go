package main

import (
    "fmt"
    "image"
    "image/color"
    "math"
    "bytes"
    "image/png"
    "io/ioutil"
    "os"
    "time"
    "strconv"
)

const IMG_WIDTH, IMG_HEIGHT := 1080, 1080

func check (e error) {
    if e != nil {
	panic(e)
    }
}

func WriteImageToFile (f string, b []byte) {
    e := ioutil.WriteFile(f, b, 0644)
    check (e)
}

type Image struct {
    x0, y0, xn, yn int
    depth uint
}

func (img Image) Bounds() image.Rectangle {
    var imrec = image.Rectangle {Min: image.Point {img.x0, img.y0}, Max: image.Point {img.xn, img.yn}}
    return imrec
}

func (image Image) ColorModel() color.Model {
    return color.RGBAModel
}

func (img Image) At (a int, b int) color.Color {
    x, y := float64(a), float64(b)
    depth := float64(img.depth)
    f, r := math.Sqrt((x*x + y*y) * depth), math.Atan(y/x)
    red := uint8(r+f)//x*(math.Sin((12*x-6*y)/2) * math.Exp((14*x-9*y)/4)))
    green := uint8(1)//uint8(128 + y*(math.Sin((x+y)/3) + math.Atan2(y/2, (x+y)/5)))
    blue := uint8(1)
    hue := uint8(255+math.Sin(360))
    return color.RGBA{red, green, blue, hue}
}

func RenderImage(m image.Image) []byte {
    var buf bytes.Buffer
    err := png.Encode(&buf, m)
    check(err)
    return buf.Bytes()
}

// func SaveImage()

func main() {
    start, e_start := strconv.Atoi(os.Args[1])
    depth, e_depth := strconv.Atoi(os.Args[2])
    delta, e_delta := strconv.Atoi(os.Args[3])
    check(e_start)
    check(e_depth)
    check(e_delta) // todo make interface for check

    out_dir := "./" + strconv.Itoa(int(time.Now().Unix()))
    e := os.MkdirAll(out_dir, 0755)
    check(e)
    for i := start; i < depth; i = i + delta {
        m := Image{1, 1, IMG_WIDTH + 1, IMG_HEIGHT + 1, uint(i)}
        ind := strconv.Itoa(i)
        f := out_dir + "/" + ind + ".png"
        WriteImageToFile(f, RenderImage(m))
        fmt.Printf("Written %s\n", f)
    }

}
