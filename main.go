package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var (
	total  = 0
	count  = 0
	prefCd = flag.String("p", "", "PrefCd: select 01-47 number.")
)

const (
	LIMIT   = 3000 // nolint
	BASEURL = "http://www.kaigokensaku.mhlw.go.jp/index.php"
)

type Result struct {
	Total string      `json:"total"`
	List  interface{} `json:"list"`
}

func main() {
	flag.Parse()

	total, count := fetchFacilityCount()

	createTempfile(count)
	mergeList := mergeTempfile(count)
	result := Result{Total: strconv.Itoa(total), List: mergeList}
	createOutputJSON(result)

	removeTempfile(count)
}

func fetchFacilityCount() (total int, count int) {
	values := initUrlValues(*prefCd)

	resp, err := http.PostForm(BASEURL, values)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	body = removeBOM(body)

	var res Result
	err = json.Unmarshal(body, &res)
	if err != nil {
		panic(err)
	}

	total, err = strconv.Atoi(res.Total)
	if err != nil {
		panic(err)
	}
	count = int(math.Ceil(float64(total) / LIMIT))

	return
}

func createTempfile(count int) {
	values := initUrlValues(*prefCd)
	prefName := PrefCd2Str(*prefCd)

	for i := 0; i < count; i++ {
		values.Set("p_offset", strconv.Itoa(i*LIMIT))
		values.Set("p_count", strconv.Itoa(LIMIT))

		resp, err := http.PostForm(BASEURL, values)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		body = removeBOM(body)

		outfile := fmt.Sprintf("%s.%d.json", prefName, i+1)

		f, err := os.Create(outfile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		w := bufio.NewWriter(f)
		_, err = w.Write(body)
		if err != nil {
			panic(err)
		}
		w.Flush()
	}
}

func mergeTempfile(count int) []interface{} {
	var mergeList []interface{}
	var result Result
	name := PrefCd2Str(*prefCd)

	for i := 0; i < count; i++ {
		file := fmt.Sprintf("%s.%d.json", name, i+1)
		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		w, _ := ioutil.ReadAll(f)
		err = json.Unmarshal(w, &result)
		if err != nil {
			panic(err)
		}

		lists := result.List.([]interface{})
		mergeList = append(mergeList, lists...)
	}
	return mergeList
}

func createOutputJSON(result Result) {
	name := PrefCd2Str(*prefCd)

	b, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	outfile := fmt.Sprintf("%s-%s.json", *prefCd, name)
	f, err := os.Create(outfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.Write(b)
	if err != nil {
		panic(err)
	}
	w.Flush()

	fmt.Printf("Created %s.\n", outfile)
}

func removeBOM(b []byte) []byte {
	return bytes.Trim(b, "\xef\xbb\xbf")
}

func removeTempfile(count int) {
	name := PrefCd2Str(*prefCd)
	for i := 0; i < count; i++ {
		err := os.Remove(fmt.Sprintf("%s.%d.json", name, i+1))
		if err != nil {
			panic(err)
		}
	}
}

func initUrlValues(prefCd string) url.Values {
	values := url.Values{}
	values.Add("action_kouhyou_splist", "true")
	values.Add("PrefCd", prefCd)
	values.Add("OriPrefCd", prefCd)
	values.Add("p_count", "1")
	values.Add("p_offset", "0")
	values.Add("p_sort_name", "6")
	values.Add("p_order_name", "0")

	return values
}
