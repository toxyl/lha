package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	_ "github.com/toxyl/termux-launch-fix"

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
		AddOptions("sort", &sortBy, []string{
			SORT_NAME, SORT_NAME_DESC,
			SORT_PERM, SORT_PERM_DESC,
			SORT_USER, SORT_USER_DESC,
			SORT_GROUP, SORT_GROUP_DESC,
			SORT_SIZE, SORT_SIZE_DESC,
			SORT_TIME, SORT_TIME_DESC,
		}, SORT_NAME, "Sets sort mode to use").
		Add("help", &showHelp, false, "Prints the help").
		Scan([]string{})
	if showHelp {
		fm.Help()
	}

	if len(os.Args) == 0 {
		os.Args = append(os.Args, ".")
	}

	config.ColorMode = !monochrome
	lower := strings.ToLower

	for _, path := range os.Args {
		pathAbs, maxLenUID, maxLenGID, totalSize, directories, files := getContents(path)

		if sortBy != SORT_NAME {
			sort.Slice(directories, func(i, j int) bool {
				d1 := directories[i]
				d2 := directories[j]
				switch sortBy {
				case SORT_NAME_DESC:
					return lower(d1.Path()) > lower(d2.Path())
				case SORT_PERM:
					return d1.Permissions().Uint() < d2.Permissions().Uint()
				case SORT_PERM_DESC:
					return d1.Permissions().Uint() > d2.Permissions().Uint()
				case SORT_USER:
					return lower(d1.Owner()) < lower(d2.Owner())
				case SORT_USER_DESC:
					return lower(d1.Owner()) > lower(d2.Owner())
				case SORT_GROUP:
					return lower(d1.Group()) < lower(d2.Group())
				case SORT_GROUP_DESC:
					return lower(d1.Group()) > lower(d2.Group())
				case SORT_SIZE:
					return d1.Size() < d2.Size()
				case SORT_SIZE_DESC:
					return d1.Size() > d2.Size()
				case SORT_TIME:
					return d1.LastModified().Unix() < d2.LastModified().Unix()
				case SORT_TIME_DESC:
					return d1.LastModified().Unix() > d2.LastModified().Unix()
				}
				return lower(d1.Path()) < lower(d2.Path())
			})
			sort.Slice(files, func(i, j int) bool {
				f1 := files[i]
				f2 := files[j]
				switch sortBy {
				case SORT_NAME_DESC:
					return lower(f1.Path()) > lower(f2.Path())
				case SORT_PERM:
					return f1.Permissions().Uint() < f2.Permissions().Uint()
				case SORT_PERM_DESC:
					return f1.Permissions().Uint() > f2.Permissions().Uint()
				case SORT_USER:
					return lower(f1.Owner()) < lower(f2.Owner())
				case SORT_USER_DESC:
					return lower(f1.Owner()) > lower(f2.Owner())
				case SORT_GROUP:
					return lower(f1.Group()) < lower(f2.Group())
				case SORT_GROUP_DESC:
					return lower(f1.Group()) > lower(f2.Group())
				case SORT_SIZE:
					return f1.Size() < f2.Size()
				case SORT_SIZE_DESC:
					return f1.Size() > f2.Size()
				case SORT_TIME:
					return f1.LastModified().Unix() < f2.LastModified().Unix()
				case SORT_TIME_DESC:
					return f1.LastModified().Unix() > f2.LastModified().Unix()
				}
				return lower(f1.Path()) < lower(f2.Path())
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
