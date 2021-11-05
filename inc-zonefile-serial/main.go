package main

import (
	"bytes"
	"fmt"
	"github.com/palladiate/go-zonefile"
	"io/ioutil"
	"os"
	"strconv"
)

// Increments the serial of a zonefile
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "<path to zonefile>")
		os.Exit(1)
	}

	// Load zonefile
	data, ioerr := ioutil.ReadFile(os.Args[1])
	if ioerr != nil {
		fmt.Println(os.Args[1], ioerr)
		os.Exit(2)
	}

	zf, perr := zonefile.Load(data)
	if perr != nil {
		fmt.Println(os.Args[1], perr.LineNo(), perr)
		os.Exit(3)
	}

	// Find SOA entry
	ok := false
	for _, e := range zf.Entries() {
		if !bytes.Equal(e.Type(), []byte("SOA")) {
			continue
		}
		vs := e.Values()
		if len(vs) != 7 {
			fmt.Println("Wrong number of parameters to SOA line")
			os.Exit(4)
		}
		serial, err := strconv.Atoi(string(vs[2]))
		if err != nil {
			fmt.Println("Could not parse serial:", err)
			os.Exit(5)
		}
		e.SetValue(2, []byte(strconv.Itoa(serial+1)))
		ok = true
		break
	}
	if !ok {
		fmt.Println("Could not find SOA entry")
		os.Exit(6)
	}

	fh, err := os.OpenFile(os.Args[1], os.O_WRONLY, 0)
	if err != nil {
		fmt.Println(os.Args[1], err)
		os.Exit(7)
	}

	_, err = fh.Write(zf.Save())
	if err != nil {
		fmt.Println(os.Args[1], err)
		os.Exit(8)
	}
}
