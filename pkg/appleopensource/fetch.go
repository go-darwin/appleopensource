// Copyright 2017 The appleopensource Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package appleopensource

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	progressbar "github.com/schollz/progressbar/v3"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
)

const (
	hdrContentLength = "Content-Length"
)

// Fetch fetchs the uri file to dst with multiple progress bars.
func Fetch(ctx context.Context, dst string, uris ...string) (err error) {
	if _, err := os.Stat(dst); err != nil && os.IsNotExist(err) {
		return fmt.Errorf("no such %s dist directory: %w", dst, err)
	}

	for _, uri := range uris {
		if err := fetch(ctx, dst, uri); err != nil {
			return err
		}
	}

	return nil
}

func fetch(ctx context.Context, dst, uri string) error {
	resp, err := http.Head(uri) // 187 MB file of random numbers per line
	if err != nil {
		return err
	}

	sz := resp.Header.Get(hdrContentLength)
	var length int64
	if sz != "" {
		length, err = strconv.ParseInt(sz, 10, 64) // Get the content length from the header request
		if err != nil {
			return err
		}
	}

	filename := path.Base(uri)
	pb := progressbar.NewOptions(int(length), progressbar.OptionShowBytes(true), progressbar.OptionSetDescription(filename))

	const limit = int64(10)    // 10 Go-routines for the process so each downloads 18.7MB
	lenSub := length / limit   // Bytes for each Go-routine
	diff := length % limit     // Get the remaining for the last request
	body := make([][]byte, 11) // Make up a temporary array to hold the data to be written to the file

	eg, ctx := errgroup.WithContext(ctx)
	for i := int64(0); i < limit; i++ {
		min := lenSub * i       // Min range
		max := lenSub * (i + 1) // Max range

		if i == limit-1 {
			max += diff // Add the remaining bytes in the last request
		}

		eg.Go(func() error {
			req, err := http.NewRequest(http.MethodGet, uri, nil)
			if err != nil {
				return err
			}

			rangeHdr := "bytes=" + strconv.FormatInt(min, 10) + "-" + strconv.FormatInt(max-1, 10) // Add the data for the Range header of the form "bytes=0-100"
			req.Header.Add("Range", rangeHdr)

			hc := http.DefaultClient
			resp, err := hc.Do(req.WithContext(ctx))
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			var buf bytes.Buffer
			out := io.MultiWriter(&buf, pb)
			if _, err := io.Copy(out, resp.Body); err != nil {
				return err
			}

			body[i] = buf.Bytes()

			return nil
		})
	}

	if err := multierr.Combine(eg.Wait(), pb.Finish()); err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(dst, filename), bytes.Join(body, nil), 0644)
}
