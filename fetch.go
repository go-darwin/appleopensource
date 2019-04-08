// Copyright 2017 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package appleopensource

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/pkg/errors"
	"github.com/vbauerster/mpb"

	"github.com/zchee/appleopensource/pkg/fs"
)

var p = mpb.New()

func fetch(wg *sync.WaitGroup, dst, uri string) error {
	defer wg.Done()

	res, err := http.Head(uri)
	if err != nil {
		return err
	}

	header := res.Header
	// Get the content length from the header request
	length, err := strconv.Atoi(header["Content-Length"][0])
	if err != nil {
		return err
	}

	filename := path.Base(uri)
	pb := p.AddBar(int64(length)).PrependName(filename, 30, 0).PrependCounters("%3s / %3s", mpb.UnitBytes, 18, mpb.DwidthSync|mpb.DextraSpace).AppendETA(5, mpb.DwidthSync)

	limit := 5                     // goroutine limit
	subLength := length / limit    // Bytes for each goroutine
	diff := length % limit         // Get the remaining for the last request
	buf := make([][]byte, limit+1) // Make up a temporary array to hold the data to be written to the file

	var (
		mu sync.Mutex
		w  sync.WaitGroup
	)
	for i := 0; i < limit; i++ {
		w.Add(1)

		min := subLength * i
		max := subLength * (i + 1)

		var rangeHeader string
		if i == limit-1 {
			max += diff // Add the remaining bytes in the last request
			rangeHeader = "bytes=" + strconv.Itoa(min) + "-" + strconv.Itoa(max)
		} else {
			rangeHeader = "bytes=" + strconv.Itoa(min) + "-" + strconv.Itoa(max-1)
		}

		go func(min int, max int, i int) {
			defer func() {
				w.Done()
			}()

			client := &http.Client{}
			req, err := http.NewRequest("GET", uri, nil)
			if err != nil {
				log.Fatal(err)
			}

			req.Header.Add("Range", rangeHeader)
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			pb.Incr(subLength)

			mu.Lock()
			buf[i], _ = ioutil.ReadAll(resp.Body)
			mu.Unlock()
		}(min, max, i)
	}
	w.Wait()
	pb.Completed()

	ioutil.WriteFile(filepath.Join(dst, filename), bytes.Join(buf[:], nil), 0644)

	return nil
}

// Fetch fetchs the uri file to dst with multiple progress bars.
func Fetch(dst string, uri ...string) (err error) {
	if !fs.IsDirExist(dst) {
		return errors.Wrapf(err, "no such %s dist directory", dst)
	}

	var wg sync.WaitGroup
	for _, u := range uri {
		wg.Add(1)
		go func(u string) {
			if e := fetch(&wg, dst, u); e != nil {
				err = e
			}
		}(u)
	}
	wg.Wait()
	if err != nil {
		return err
	}
	p.Stop()

	return nil
}
