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

const (
    IMG_WIDTH int = 1080
    IMG_HEIGHT int = 1080
)

func check (e error) {
    if e != nil {
	panic(e)
    }
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
    red := uint8(f + r)//x*(math.Sin((12*x-6*y)/2) * math.Exp((14*x-9*y)/4)))
    green := uint8(10*r - 49032*f)//uint8(128 + y*(math.Sin((x+y)/3) + math.Atan2(y/2, (x+y)/5)))
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

func WriteImageToFile (f string, b []byte) int {
    e := ioutil.WriteFile(f, b, 0644)
    check (e)
    fmt.Printf("Written %s\n", f)
    return len(b)
}

func SaveImage(i int, dir string) int {
        m := Image{1, 1, IMG_WIDTH + 1, IMG_HEIGHT + 1, uint(i)}
        ind := strconv.Itoa(i)
        f := dir + "/" + ind + ".png"
        fmt.Println("Processing ", ind);
        go WriteImageToFile(f, RenderImage(m))
        return 0
}

func quit() {
    fmt.Println("QUITTING")
}
func main() {
    defer quit()
    start, e_start := strconv.Atoi(os.Args[1])
    depth, e_depth := strconv.Atoi(os.Args[2])
    delta, e_delta := strconv.Atoi(os.Args[3])
    check(e_start)
    check(e_depth)
    check(e_delta) // todo make interface for check
    out_dir := "./" + strconv.Itoa(int(time.Now().Unix()))
    // convert_cmd := "convert -delay 1 -loop 1 *.png 00000anim.gif"
    e := os.MkdirAll(out_dir, 0755)
    check(e)
    for i := start; i < depth; i = i + delta {
        go SaveImage(i, out_dir)
    }

}
