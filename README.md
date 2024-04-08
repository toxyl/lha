# LHA
... is a Go implementation of `ls -lha` with some additional features.

# Installation
```sh
CGO_ENABLED=0 go build .
cp lha /usr/local/bin/
chmod +x /usr/local/bin/lha
```

# Usage
```sh
lha --help
```

```sh
Command
Usage:   lha <flags> <paths>
Example: lha --sort time /tmp /home

Flags
Usage of /usr/local/bin/lha:
  -help
        Prints the help
    
  -monochrome
        Prints monochrome output
    
  -sort string
        Defines how to sort the output
        Options: name, name-desc, perm, perm-desc, user, user-desc, group, group-desc, size, size-desc, time, time-desc
         (default "name")
```
If no path is given the current directory will be inspected.

