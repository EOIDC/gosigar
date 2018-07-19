// Copyright (c) 2012 VMware, Inc.

package main

import (
	"fmt"
	"os"

	"github.com/eoidc/gosigar"
)

const output_format = "%-15s %4s %4s %5s %4s %-15s\n"

func main() {
	fslist := gosigar.FileSystemList{}
	fslist.Get()

	fmt.Fprintf(os.Stdout, output_format,
		"Filesystem", "Size", "Used", "Avail", "Use%", "Mounted on")

	for _, fs := range fslist.List {
		dir_name := fs.DirName
		usage := gosigar.FileSystemUsage{}

		err := gosigar.FsPing(fs)
		if  err == nil {
			usage.Get(dir_name)
		} else {
			fmt.Println(err)
		}

		fmt.Fprintf(os.Stdout, output_format,
			fs.DevName,
			gosigar.FormatSize(usage.Total),
			gosigar.FormatSize(usage.Used),
			gosigar.FormatSize(usage.Avail),
			gosigar.FormatPercent(usage.UsePercent()),
			dir_name)
	}
}
