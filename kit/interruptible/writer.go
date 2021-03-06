// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package interruptible

import (
	"io"
	"sync"

	"github.com/gocircuit/circuit/kit/fs/rh" // TODO: backwards dep
)

// Writer is an io.WriteCloser which also supports interruptible writes.
type Writer interface {
	io.WriteCloser
	WriteIntr([]byte, rh.Intr) (int, error)
}

//
type writer struct {
	w struct { // WriteIntr structure
		Mutex
		ch    chan<- []byte
		abort <-chan struct{}
	}
	a struct { // Close/abort structure
		sync.Mutex
		abort chan<- struct{}
	}
	s struct { // Stats 
		sync.Mutex
		closed bool
		n      int64
	}
}

func (w *writer) Write(p []byte) (int, error) {
	return w.WriteIntr(p, nil)
}

func (w *writer) WriteIntr(p []byte, intr rh.Intr) (n int, err error) {
	u := w.w.Lock(intr)
	if u == nil {
		return 0, rh.ErrIntr
	}
	defer u.Unlock()
	//
	if w.w.ch == nil {
		return 0, io.ErrClosedPipe
	}
	select {
	case w.w.ch <- p:
		w.s.Lock()
		defer w.s.Unlock()
		//
		w.s.n += int64(len(p))
		return len(p), nil
	case <-intr:
		return 0, io.ErrNoProgress
	case <-w.w.abort: // If we receive an abort during write, close out channel here
		close(w.w.ch)
		w.w.ch = nil
		w.close()
		return 0, io.ErrUnexpectedEOF
	}
}

func (w *writer) stop() {
	u := w.w.Lock(nil)
	if u == nil {
		panic(0)
	}
	defer u.Unlock()
	//
	if w.w.ch == nil {
		return
	}
	close(w.w.ch)
	w.w.ch = nil
	//
	w.close()
}

func (w *writer) close() {
	w.s.Lock()
	defer w.s.Unlock()
	//
	w.s.closed = true
}

func (w *writer) abort() {
	w.a.Lock()
	defer w.a.Unlock()
	//
	if w.a.abort == nil {
		return
	}
	close(w.a.abort)
	w.a.abort = nil
}

// Close will interrupt a pending write.
func (w *writer) Close() error {
	w.abort()
	w.stop()
	return nil
}

func (w *writer) Stat() (nsent int64, closed bool) {
	w.s.Lock()
	defer w.s.Unlock()
	//
	return w.s.n, w.s.closed
}
