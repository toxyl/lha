package main

import (
	"github.com/toxyl/flo"
	"github.com/toxyl/glog"
)

func getContents(path string) (pathAbs string, maxLenUID, maxLenGID, totalSize int, directories []*flo.DirObj, files []*flo.FileObj) {
	d := flo.Dir(path)
	pathAbs = d.Path()
	maxLenUID, maxLenGID, totalSize = 0, 0, 0
	directories, files = []*flo.DirObj{}, []*flo.FileObj{}
	if !d.Info().Permissions.IsDir() {
		f := flo.File(d.Path())
		files = append(files, f)
		maxLenUID, maxLenGID = len(f.Owner()), len(f.Group())
	}
	d.EachLimit(
		func(f *flo.FileObj) {
			files = append(files, f)
			maxLenUID = glog.Max(maxLenUID, len(f.Owner()))
			maxLenGID = glog.Max(maxLenGID, len(f.Group()))
		},
		func(d *flo.DirObj) {
			directories = append(directories, d)
			maxLenUID = glog.Max(maxLenUID, len(d.Owner()))
			maxLenGID = glog.Max(maxLenGID, len(d.Group()))
		},
		0,
	)
	return
}
