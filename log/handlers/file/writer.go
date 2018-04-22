package file

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"aahframework.org/essentials.v0"
)

type writer struct {
	*sync.Mutex
	messages []string

	out      io.Writer
	filename string

	openDay  time.Time
	stats    *stats
	isClosed bool

	fileFormat string
	folder     string
	maxSize    int64
	maxLines   int64
}

func (w *writer) Init() error {
	w.filename = w.getFile()

	return w.openFile()
}

func (w *writer) Write(message []byte) {
	w.Lock()
	defer w.Unlock()

	if w.isRotate() {
		_ = w.rotateFile()

		// reset rotation values
		w.openDay = time.Now()
		w.stats.lines = 0
		w.stats.bytes = 0
	}

	size, _ := w.out.Write(message)

	// calculate receiver stats
	w.stats.bytes += int64(size)
	w.stats.lines += int64(bytes.Count(message, []byte("\n")))
}

func (w *writer) open(file string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(file, flag, perm)
}

func (w *writer) isRotate() bool {
	if w.maxLines > 0 && w.stats.Lines() >= w.maxLines {
		return true
	}

	if w.maxSize > 0 && w.stats.Bytes() >= w.maxSize {
		return true
	}

	if !w.openDay.IsZero() && !w.sameDay() {
		return true
	}

	return false
}

func (w *writer) rotateFile() error {
	if _, err := os.Lstat(filepath.Join(w.folder, w.filename)); err == nil {
		w.close()
		if err = os.Rename(filepath.Join(w.folder, w.filename), w.getFile()); err != nil {
			return err
		}
	}

	w.filename = w.getFile()

	if err := w.openFile(); err != nil {
		return err
	}

	return nil
}

func (w *writer) openFile() error {
	err := ess.MkDirAll(w.folder, 0755)
	if err != nil {
		return nil
	}

	file, err := os.OpenFile(w.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	w.out = file
	w.isClosed = false
	w.stats = &stats{}
	w.stats.bytes = fileStat.Size()
	w.stats.lines = int64(ess.LineCntr(file))

	return nil
}

func (w *writer) sameDay() bool {
	y1, m1, d1 := w.openDay.Date()
	y2, m2, d2 := time.Now().Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (w *writer) close() {
	if !w.isClosed {
		ess.CloseQuietly(w.out)
		w.isClosed = false
	}
}

func (w *writer) getFile() string {
	file := w.fileFormat
	date := time.Now().Format("2006-01-02")
	file = strings.Replace(file, "%date%", date, -1)
	return filepath.Join(w.folder, file)
}
