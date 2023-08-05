package downloader

import (
	"io"
	"net/http"
	"os"
	"strconv"
)

type downloader struct {
	total, current int64
	url            string
	progress       chan int
	errs           chan error
}

// Write implements io.Writer.
func (d *downloader) Write(p []byte) (int, error) {
	n := len(p)
	d.current += int64(n)

	d.notify()

	return n, nil
}

func (d *downloader) notify() {
	progress := (float64(d.current) / float64(d.total)) * 100
	d.progress <- int(progress)
}

func (d *downloader) download(file *os.File, dst string) {
	resp, err := http.Get(d.url)
	if err != nil {
		_ = file.Close()
		d.errs <- err
		return
	}

	defer resp.Body.Close()
	if _, err := io.Copy(file, io.TeeReader(resp.Body, d)); err != nil {
		_ = file.Close()
		d.errs <- err
		return
	}

	// Close file before rename
	if err := file.Close(); err != nil {
		d.errs <- err
		return
	}
	if err := os.Rename(file.Name(), dst); err != nil {
		d.errs <- err
	}

}

func Download(dstFilePath, srcUrl string) (progress chan int, errCh chan error, err error) {
	// Fetch head, so we can verify  if url is correct and get total size of file
	head, err := http.Head(srcUrl)
	if err != nil {
		return
	}

	d := &downloader{
		url:      srcUrl,
		progress: make(chan int, 100),
		errs:     make(chan error, 1),
	}

	d.total, err = strconv.ParseInt(head.Header.Get("Content-length"), 10, 64)
	if err != nil {
		return
	}

	file, err := os.CreateTemp(os.TempDir(), "*.mender")
	if err != nil {
		return
	}

	go d.download(file, dstFilePath)

	return d.progress, d.errs, nil
}
