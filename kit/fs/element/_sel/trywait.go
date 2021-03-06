// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package sel

import (
	"bytes"

	"github.com/gocircuit/circuit/kit/iomisc"
	"github.com/gocircuit/circuit/kit/fs/rh"
	"github.com/gocircuit/circuit/kit/fs/namespace/file"
)

type TryWaitFile struct {
	s *Select
}

func NewTryWaitFile(s *Select) file.File {
	return &TryWaitFile{s: s}
}

func (f *TryWaitFile) Perm() rh.Perm {
	return 0444 // r--r--r--
}

func (f *TryWaitFile) Open(flag rh.Flag, _ rh.Intr) (_ rh.FID, err error) {
	f.s.ErrorFile.Clear() // clear error file
	if flag.Attr != rh.ReadOnly {
		return nil, rh.ErrPerm
	}
	var u Unblock
	u.Clause, u.Commit, u.Error = f.s.TryWait()
	return file.NewOpenReaderFile(iomisc.ReaderNopCloser(bytes.NewBufferString(marshal(u)))), nil
}

func (f *TryWaitFile) Remove() error {
	return rh.ErrPerm
}
