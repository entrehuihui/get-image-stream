package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gocv.io/x/gocv"
)

func main() {
	resp, err := http.Get("http://192.168.191.2:8000/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	window := gocv.NewWindow("-----")
	defer window.Close()
	for {
		_, err := reader.ReadSlice(10)
		if err != nil {
			fmt.Println("1", err)
			continue
		}
		line, err := reader.ReadSlice(10)
		if err != nil {
			fmt.Println("2", err)
			continue
		}
		if string(line) != "--MJPEGBOUNDARY\r\n" {
			// fmt.Println(line)
			continue
		}
		_, err = reader.ReadSlice(10)
		if err != nil {
			fmt.Println("3", err)
			continue
		}
		line, err = reader.ReadSlice(10)
		if err != nil {
			fmt.Println(err)
			continue
		}
		lens := strings.Split(string(line), ":")
		if len(lens) != 2 {
			continue
		}
		lens[1] = strings.Trim(lens[1], " ")
		length, err := strconv.Atoi(strings.Replace(lens[1], "\r\n", "", -1))
		if err != nil {
			fmt.Println(len(lens[1]), err)
			continue
		}
		fmt.Println(length)
		// strings.Repeat()
		_, err = reader.ReadSlice(10)
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, err = reader.ReadSlice(10)
		if err != nil {
			fmt.Println(err)
			continue
		}
		buf := make([]byte, length)
		// n, err := reader.Read(buf)
		// if err != nil {
		// 	fmt.Println(n, err)
		// 	continue
		// }
		// fmt.Println(n, "---------")

		n := 0
		for {
			n1, err := reader.Read(buf[n:])
			if err != nil {
				break
			}
			n += n1
			if n >= length {
				break
			}
		}
		if n != length {
			continue
		}
		img, err := gocv.IMDecode(buf, -1)
		if err != nil {
			continue
		}
		window.IMShow(img)
		if window.WaitKey(1) > 0 {
			break
		}
	}
}
