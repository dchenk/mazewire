package main

import (
	"bytes"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	options()
	blobs()
}

// Generate random data for a "blobs" table.
func blobs() {
	buf := new(bytes.Buffer)
	buf.WriteString("INSERT INTO blobs2 (role,k,v) VALUES ")
	ln := false
	const maxLenRole = 255
	const numRows = 150
	rolesWritten := make([][]rune, 0, numRows)
	roleBuf := make([]rune, 0, maxLenRole)
	for i := 0; i < numRows; i++ {
		buf.WriteString("('")
		buf.WriteString(string(genUniqueStr(maxLenRole, roleBuf, rolesWritten)))
		buf.WriteString("',")
		buf.WriteString(strconv.FormatUint(uint64(rand.Uint32()), 10)) // the "k" column
		buf.WriteString(",'")
		vChars := rand.Intn(3000)
		for j := 0; j < vChars; j++ {
			buf.WriteRune(randomRune()) // The "role" column
		}
		buf.WriteString("')")
		if i < numRows-1 {
			buf.WriteByte(',')
		}
		if ln { // Two rows to insert per line.
			buf.WriteByte('\n')
		}
		ln = !ln
	}
	buf.WriteByte(';')
	writeOut("blobs.sql", buf)
}

// Generate random data for an "options" table.
func options() {
	buf := new(bytes.Buffer)
	buf.WriteString("INSERT INTO options_b (site,k,v) VALUES ")
	ln := false
	const maxLenK = 255
	const numRows = 100
	ks := make([][]rune, 0, numRows)
	kBuf := make([]rune, 0, maxLenK)
	for i := 0; i < numRows; i++ {
		buf.WriteString("(" + strconv.FormatInt(1, 10) + ",'") // the "site" ID
		buf.WriteString(string(genUniqueStr(maxLenK, kBuf, ks)))
		buf.WriteString("','")
		vChars := rand.Intn(4500) // Can be up to 7680
		for j := 0; j < vChars; j++ {
			buf.WriteRune(randomRune()) // The "v" column
		}
		buf.WriteString("')")
		if i < numRows-1 {
			buf.WriteByte(',')
		}
		if ln { // Two rows to insert per line.
			buf.WriteByte('\n')
		}
		ln = !ln
	}
	buf.WriteByte(';')
	writeOut("options.sql", buf)
}

// maxLen is the longest any of the strings will get.
func genUniqueStr(maxLen int, buf []rune, existing [][]rune) []rune {
	theLen := rand.Intn(maxLen + 1)
	buf = buf[:theLen] // reset the buffer
	rand.Seed(time.Now().UnixNano())
	for j := 0; j < theLen; j++ {
		buf[j] = randomRune()
	}
	for str := range existing { // check if this string already exists
		if string(buf) == string(str) {
			log.Println("found matching string in buffer")
			genUniqueStr(maxLen, buf, existing) // generate another string
		}
	}
	existing = append(existing, buf)
	return buf
}

func randomRune() rune {
	r := rune(rand.Intn(91) + 35)
	switch r {
	case '\'': // single quote
		return 'A'
	case '\\': // backslash
		return 'a'
	}
	return r
}

// writeOut writes the file out into the directory of this source code file (package main) with "gen_" prepended to localName.
// The function panics if an error occurs.
func writeOut(localName string, buf *bytes.Buffer) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(dir + "/gen_" + localName)
	if err != nil {
		log.Fatal(f)
	}
	defer f.Close()
	_, err = buf.WriteTo(f)
	if err != nil {
		log.Fatal("doing buffer to file:", err)
	}
}
