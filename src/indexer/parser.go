package indexer

import (
	"bufio"
	"db"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	eGEONAMEID = iota
	eNAME
	eASCIINAME
	eALTERNATENAMES
	eLATITUDE
	eLONGITUDE
	eFEATURE_CLASS
	eFEATURE_CODE
	eCOUNTRY_CODE
	eCC2
	eADMIN1_CODE
	eADMIN2_CODE
	eADMIN3_CODE
	eADMIN4_CODE
	ePOPULATION
	eELEVATION
	eDEM
	eTIMEZONE
	eMODIFICATION_DATE
)

func ReadFile(path string, d *db.Data) bool {
	start := time.Now()
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scan := bufio.NewScanner(file)

	d.Cities4Search = make(map[int64]db.City)

	//parsing by line
	var s string
	var count int
	for scan.Scan() {
		s = scan.Text()
		list := strings.Split(s, "\t")

		var c db.City

		for i := 0; i < len(list); i++ {
			switch i {
			case eGEONAMEID:
				c.Geonameid, _ = strconv.ParseInt(list[i], 10, 64)
			case eNAME:
				c.Name = list[i]
			case eASCIINAME:
				c.Asciiname = list[i]
			case eALTERNATENAMES:
				c.Alternatenames = list[i]
			case eLATITUDE:
				c.Latitude = list[i]
			case eLONGITUDE:
				c.Longitude = list[i]
			case eFEATURE_CLASS:
				c.Feature_class = list[i]
			case eFEATURE_CODE:
				c.Feature_code = list[i]
			case eCOUNTRY_CODE:
				c.Country_code = list[i]
			case eCC2:
				c.Cc2 = list[i]
			case eADMIN1_CODE:
				c.Admin1_code = list[i]
			case eADMIN2_CODE:
				c.Admin2_code = list[i]
			case eADMIN3_CODE:
				c.Admin3_code = list[i]
			case eADMIN4_CODE:
				c.Admin4_code = list[i]
			case ePOPULATION:
				c.Population = list[i]
			case eELEVATION:
				c.Elevation = list[i]
			case eDEM:
				c.Dem = list[i]
			case eTIMEZONE:
				c.Timezone = list[i]
			case eMODIFICATION_DATE:
				c.Modification_date = list[i]
			}
		}

		if c.Geonameid != 0 {
			d.Cities4Search[c.Geonameid] = c
			d.Cities4Indexer = append(d.Cities4Indexer, &c)
		}

		if len(s) > 0 {
			count++
		}
	}

	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	elapsed := t.Sub(start)

	fmt.Println(count, " in ", elapsed)
	fmt.Println("Map size ", len(d.Cities4Search))
	//fmt.Println("Size of Map ", len(d.Cities4Search), " Size of city: ", unsafe.Sizeof(c))

	return false
}
