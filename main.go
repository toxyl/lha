package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/toxyl/flo/config"
	"github.com/toxyl/flo/log"
	"github.com/toxyl/flo/utils"
	"github.com/toxyl/glog"
)

func init() {
	log.SetFns(nil, nil) // errors can happen when accessing unknown files, just ignore them
}

const (
	SORT_NAME       = "name"
	SORT_NAME_DESC  = "name-desc"
	SORT_PERM       = "perm"
	SORT_PERM_DESC  = "perm-desc"
	SORT_USER       = "user"
	SORT_USER_DESC  = "user-desc"
	SORT_GROUP      = "group"
	SORT_GROUP_DESC = "group-desc"
	SORT_SIZE       = "size"
	SORT_SIZE_DESC  = "size-desc"
	SORT_TIME       = "time"
	SORT_TIME_DESC  = "time-desc"
)

func main() {
	var monochrome, showHelp bool
	var sortBy string
	fm := NewFlagMap().
		Add("monochrome", &monochrome, false, "Prints monochrome output").
		AddOptions("sort", &sortBy, []string{SORT_NAME, SORT_NAME_DESC, SORT_PERM, SORT_PERM_DESC, SORT_USER, SORT_USER_DESC, SORT_GROUP, SORT_GROUP_DESC, SORT_SIZE, SORT_SIZE_DESC, SORT_TIME, SORT_TIME_DESC}, SORT_NAME, "Defines how to sort the output").
		Add("help", &showHelp, false, "Prints the help").
		Scan([]string{})
	if showHelp {
		fm.Help()
	}

	if len(os.Args) == 1 {
		os.Args = append(os.Args, ".")
	}

	config.ColorMode = !monochrome

	for _, path := range os.Args[1:] {
		pathAbs, maxLenUID, maxLenGID, totalSize, directories, files := getContents(path)

		if sortBy != SORT_NAME {
			// let's sort
			sort.Slice(directories, func(i, j int) bool {
				d1 := directories[i]
				d2 := directories[j]
				switch sortBy {
				case SORT_NAME_DESC:
					return strings.ToLower(d1.Path()) > strings.ToLower(d2.Path())
				case SORT_PERM:
					return d1.Permissions().Uint() < d2.Permissions().Uint()
				case SORT_PERM_DESC:
					return d1.Permissions().Uint() > d2.Permissions().Uint()
				case SORT_USER:
					return strings.ToLower(d1.Owner()) < strings.ToLower(d2.Owner())
				case SORT_USER_DESC:
					return strings.ToLower(d1.Owner()) > strings.ToLower(d2.Owner())
				case SORT_GROUP:
					return strings.ToLower(d1.Group()) < strings.ToLower(d2.Group())
				case SORT_GROUP_DESC:
					return strings.ToLower(d1.Group()) > strings.ToLower(d2.Group())
				case SORT_SIZE:
					return d1.Size() < d2.Size()
				case SORT_SIZE_DESC:
					return d1.Size() > d2.Size()
				case SORT_TIME:
					return d1.LastModified().Unix() < d2.LastModified().Unix()
				case SORT_TIME_DESC:
					return d1.LastModified().Unix() > d2.LastModified().Unix()
				}
				return strings.ToLower(d1.Path()) < strings.ToLower(d2.Path())
			})
			sort.Slice(files, func(i, j int) bool {
				f1 := files[i]
				f2 := files[j]
				switch sortBy {
				case SORT_NAME_DESC:
					return strings.ToLower(f1.Path()) > strings.ToLower(f2.Path())
				case SORT_PERM:
					return f1.Permissions().Uint() < f2.Permissions().Uint()
				case SORT_PERM_DESC:
					return f1.Permissions().Uint() > f2.Permissions().Uint()
				case SORT_USER:
					return strings.ToLower(f1.Owner()) < strings.ToLower(f2.Owner())
				case SORT_USER_DESC:
					return strings.ToLower(f1.Owner()) > strings.ToLower(f2.Owner())
				case SORT_GROUP:
					return strings.ToLower(f1.Group()) < strings.ToLower(f2.Group())
				case SORT_GROUP_DESC:
					return strings.ToLower(f1.Group()) > strings.ToLower(f2.Group())
				case SORT_SIZE:
					return f1.Size() < f2.Size()
				case SORT_SIZE_DESC:
					return f1.Size() > f2.Size()
				case SORT_TIME:
					return f1.LastModified().Unix() < f2.LastModified().Unix()
				case SORT_TIME_DESC:
					return f1.LastModified().Unix() > f2.LastModified().Unix()
				}
				return strings.ToLower(f1.Path()) < strings.ToLower(f2.Path())
			})
		}

		str := utils.NewString().LF().StrUnderline(pathAbs).LF().LF()

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
}
