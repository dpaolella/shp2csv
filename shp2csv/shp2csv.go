/*
This program saves data stored in a shapefile to CSV (one row per polygon)
*/

package shp2csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ctessum/geom/encoding/shp"
)

type polygonData struct {
	fields map[string]float64
}

func Run(dirFlag bool, path string) {

	if dirFlag == true {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if filepath.Ext(path+"/"+file.Name()) == ".shp" {

				fileBase := file.Name()[:len(file.Name())-4]

				fmt.Println("Loading " + file.Name())
				records, fieldNames := loadShp(path + "/" + file.Name())
				csvData := make([][]string, len(records)+1)

				for _, f := range fieldNames {
					csvData[0] = append(csvData[0], f)
				}

				for i, rec := range records {
					for _, f := range fieldNames {
						csvData[i+1] = append(csvData[i+1], strconv.FormatFloat(rec.fields[f], 'f', -1, 64))
					}
				}

				outfile, err := os.Create(path + "/" + fileBase + ".csv")
				checkError("Cannot create file", err)
				defer outfile.Close()

				fmt.Println("Exporting " + fileBase + ".csv\n")
				writer := csv.NewWriter(outfile)
				defer writer.Flush()

				for _, value := range csvData {
					err := writer.Write(value)
					checkError("Cannot write to file", err)
				}
			}
		}
	} else {
		fileBase := path[:len(path)-4]

		fmt.Println("Loading " + path)
		records, fieldNames := loadShp(path)
		csvData := make([][]string, len(records)+1)

		for _, f := range fieldNames {
			csvData[0] = append(csvData[0], f)
		}

		for i, rec := range records {
			for _, v := range rec.fields {
				csvData[i+1] = append(csvData[i+1], strconv.FormatFloat(v, 'f', -1, 64))
			}
		}

		outfile, err := os.Create(fileBase + ".csv")
		checkError("Cannot create file", err)
		defer outfile.Close()

		fmt.Println("Exporting " + fileBase + ".csv\n")
		writer := csv.NewWriter(outfile)
		defer writer.Flush()

		for _, value := range csvData {
			err := writer.Write(value)
			checkError("Cannot write to file", err)
		}
	}
}

func loadShp(path string) (records []polygonData, fieldNames []string) {

	// Load geometries.
	d, err := shp.NewDecoder(path)
	if err != nil {
		panic(err)
	}

	for _, f := range d.Fields() {
		fieldNames = append(fieldNames, shpFieldName2String(f.Name))
	}

	// Decode a record from the input file.
	for {

		var rec polygonData
		_, f, more := d.DecodeRowFields(fieldNames...)
		if !more {
			break
		}

		rec.fields = make(map[string]float64, len(f))
		for k, v := range f {
			rec.fields[k], err = s2f(v)
			if err != nil {
				panic(err)
			}
		}
		records = append(records, rec)
	}

	// Check to see if any errors occured during decoding.
	if err = d.Error(); err != nil {
		panic(err)
	}

	return records, fieldNames

}

func shpFieldName2String(name [11]byte) string {
	b := bytes.Trim(name[:], "\x00")
	n := bytes.Index(b, []byte{0})
	if n == -1 {
		n = len(b)
	}
	return strings.TrimSpace(string(b[0:n]))
}

func s2f(s string) (float64, error) {
	if removeNull(s) == "************************" { // Null value
		return 0., nil
	}
	f, err := strconv.ParseFloat(removeNull(s), 64)
	return f, err
}

func removeNull(s string) string {
	s = s[0 : len(s)-strings.Count(s, "\x00")]
	return s
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
