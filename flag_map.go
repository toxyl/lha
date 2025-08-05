package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/toxyl/glog"
)

type FlagMap struct {
	data    map[string]interface{}
	targets map[string]interface{}
	options map[string][]string
}

func (fm *FlagMap) checkForDuplicateFlag(name string) {
	if _, ok := fm.data[name]; ok {
		panic(fmt.Sprintf("The flag '%s' has already been defined.", name))
	}
}

func (fm *FlagMap) Add(name string, target, def interface{}, desc string) *FlagMap {
	fm.checkForDuplicateFlag(name)
	desc += "\n"
	switch t := def.(type) {
	case string:
		fm.data[name] = flag.String(name, t, desc)
	case bool:
		fm.data[name] = flag.Bool(name, t, desc)
	default:
		panic("Only string and bool flags are implemented")
	}
	fm.targets[name] = target
	return fm
}

func (fm *FlagMap) AddOptions(name string, target *string, options []string, def string, desc string) *FlagMap {
	fm.checkForDuplicateFlag(name)
	desc += "\n"
	desc += "Options: " + strings.Join(options, ", ")

	fm.data[name] = flag.String(name, def, desc)
	fm.targets[name] = target
	fm.options[name] = options
	return fm
}

func (fm *FlagMap) Get(name string) interface{} {
	if f, ok := fm.data[name]; ok {
		switch t := f.(type) {
		case *string:
			return *t
		case *bool:
			return *t
		default:
			panic("Only string and bool flags are implemented")
		}
	}
	return nil
}

func (fm *FlagMap) Help() {
	fmt.Printf("%s\nUsage:   %s <flags> <paths>\n", glog.Bold()+"Command"+glog.Reset(), filepath.Base(os.Args[0]))
	fmt.Printf("Example: %s --sort time /tmp /home\n\n%s\n", filepath.Base(os.Args[0]), glog.Bold()+"Flags"+glog.Reset())
	flag.Usage()
	fmt.Println()
	os.Exit(0)
}

func (fm *FlagMap) String(name string) string {
	return fm.Get(name).(string)
}

func (fm *FlagMap) Bool(name string) bool {
	return fm.Get(name).(bool)
}

func (fm *FlagMap) Scan(requiredStringFlags []string) *FlagMap {
	flag.Parse()
	for k, f := range fm.data {
		target := fm.targets[k]
		switch t := target.(type) {
		case *string:
			v := *f.(*string)
			if opt, ok := fm.options[k]; ok {
				if len(opt) > 0 {
					isValid := false
					for _, o := range opt {
						if o == v {
							isValid = true
							break
						}
					}
					if !isValid {
						fmt.Printf("\nInvalid '--%s' option, use one of these: %s\n\n", k, strings.Join(opt, ", "))
						os.Exit(1)
					}
				}
			}
			(*t) = v
		case *bool:
			(*t) = *f.(*bool)
		default:
			panic("Only string and bool flags are implemented")
		}
	}
	os.Args = flag.Args()
	return fm
}

func NewFlagMap() *FlagMap {
	fm := &FlagMap{
		data:    map[string]interface{}{},
		targets: map[string]interface{}{},
		options: map[string][]string{},
	}
	return fm
}
