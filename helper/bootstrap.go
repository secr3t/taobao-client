package helper

import (
	"encoding/json"
	"fmt"
	"github.com/secr3t/rakuten-taobao-client/client"
	"github.com/thoas/go-funk"
	"image"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

var (
	separator  = regexp.MustCompile(`[\r\n]+`)
	httpClient = http.DefaultClient
)

func readString(path string) string {
	bytes, _ := ioutil.ReadFile(path)
	return string(bytes)
}

func GetInterval() int {
	interval, _ := strconv.ParseInt(readString("interval.txt"), 10, 32)
	return int(interval)
}

func GetUrls() []string {
	urls := separator.Split(readString("urls.txt"), -1)
	return funk.Filter(urls, func(s string) bool {
		return strings.Trim(s, " ") != ""
	}).([]string)
}

func ParseUri(uri string) url.Values {
	parse, _ := url.Parse(uri)
	return parse.Query()
}

func GetStartEndPrice(filter string) (startPrice int, endPrice int) {
	reg, _ := regexp.Compile(`reserve_price\[(\d+)?,(\d+)?\]`)

	matched := reg.FindStringSubmatch(filter)

	if len(matched) > 2 {
		startPrice, _ = strconv.Atoi(matched[1])
		endPrice, _ = strconv.Atoi(matched[2])
	}

	if len(matched) == 2 {
		startPrice, _ = strconv.Atoi(matched[1])
	}

	return
}

func GetCurrentDir() string {
	return time.Now().Format("2006-01-02T150405")
}

func GetImgNum() int {
	imgNum, _ := strconv.ParseInt(readString("imgNum.txt"), 10, 64)
	return int(imgNum)
}

func SaveImgNum(imgNum int) {
	_ = ioutil.WriteFile("imgNum.txt", []byte(fmt.Sprint(imgNum)), os.ModePerm)
}

func MakeDir(path string) {
	os.Mkdir(path, os.ModePerm)
	dirs := strings.Split(path, "/")
	os.Mkdir(filepath.Join(dirs[0],"diff-"+dirs[1]), os.ModePerm)
}

func SaveJson(path string, datum interface{}, num int) {
	realPath := filepath.Join(path, fmt.Sprint(num)+".json")
	marshalJson, _ := json.Marshal(datum)
	_ = ioutil.WriteFile(realPath, marshalJson, os.ModePerm)
}

func SaveImage(url, path string, imgNum int) bool {
	img := readImageFromUrl(url)
	if img == nil {
		return false
	}

	path = filepath.Join(path, fmt.Sprint(imgNum)+".jpg")

	saveImage(img, path)
	saveDiffImage(img, path)

	return true
}

func readImageFromUrl(url string) image.Image {
	res, err := httpClient.Get(url)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)

	if err != nil {
		return nil
	}

	return img
}

func saveDiffImage(img image.Image, path string) {
	diffImg := imaging.Resize(img, 200, 200, imaging.Lanczos)
	dir, fileName := filepath.Split(path)
	dirs := strings.Split(dir, "/")

	saveImage(diffImg, filepath.Join(dirs[0], "diff-"+dirs[1], fileName))
}

func saveImage(img image.Image, path string) {
	_ = imaging.Save(img, path)
}

func GetRakutenClientParamsFromUri(uri string) client.SearchParam {
	values := ParseUri(uri)

	catId, _ := strconv.Atoi(values.Get("cat"))
	startPrice, endPrice := GetStartEndPrice(values.Get("filter"))

	return client.NewSearchParam(
		values.Get("q"),
		values.Get("sort"),
		1,
		100,
		startPrice,
		endPrice,
		catId,
	)
}

func PriceAsFloat(price string) float64 {
	p, _ := strconv.ParseFloat(price, 64)
	return p
}
