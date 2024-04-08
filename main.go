package main

import (
	"fmt"
	"os"

	"github.com/toxyl/flo"
	"github.com/toxyl/flo/config"
	"github.com/toxyl/flo/log"
	"github.com/toxyl/flo/utils"
	"github.com/toxyl/glog"
)

func init() {
	log.SetFns(nil, nil) // errors can happen when accessing unknown files, just ignore them
}

func main() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, ".")
	}
	d := flo.Dir(os.Args[1])
	maxLenUID, maxLenGID, totalSize := 0, 0, 0
	directories := []*flo.DirObj{}
	files := []*flo.FileObj{}
	if !d.Info().Permissions.IsDir() {
		f := flo.File(d.Path())
		files = append(files, f)
		maxLenUID, maxLenGID = len(f.Owner()), len(f.Group())
	}

	str := utils.NewString().LF().StrUnderline(d.Path()).LF().LF()

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

	for _, d := range directories {
		str = str.Pad(1).Str(d.String(maxLenUID, maxLenGID)).LF()
	}

	for _, f := range files {
		totalSize += int(f.Size())
		str = str.Pad(1).Str(f.String(maxLenUID, maxLenGID)).LF()
	}
	str.LF().
		StrClean(!config.ColorMode,
			fmt.Sprintf("%s, %s (%s)\n",
				glog.IntAmount(len(directories), "directory", "directories"),
				glog.IntAmount(len(files), "file", "files"),
				glog.HumanReadableBytesIEC(totalSize))).
		LF().
		Print()
}
