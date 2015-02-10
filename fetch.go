package main

import "os"
import "log"
import "net/http"
import "io"

func fetch(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		return err
	}
	return resp.Body.Close()

}

func main() {
	for _, arg := range os.Args[1:] {
		err := fetch(arg)
		if err != nil {
			log.Fatal(err)
		}
	}
}
