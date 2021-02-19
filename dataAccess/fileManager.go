package dataAccess

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	fp "path/filepath"
	"sort"
	c "speedСontrol/caches"
	m "speedСontrol/models"
	"strconv"
)

const db = "db"

var dbDir string

func init() {
	wd, _ := os.Getwd()
	dbDir = fp.Join(wd, db)
}

func SaveSpeedControlInfo(m m.SpeedControlMsg) {
	path := fp.Join(dbDir, m.Year, m.Month, m.Day)

	if err := ensureDir(path); err != nil {
		log.Fatal("Directory creation failed with error: ", err)
	}

	f, err := os.OpenFile(fp.Join(path, m.ParsedSpeedString), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if f != nil {
		defer f.Close()
	}

	if err != nil {
		log.Fatal("Unable to create file:", err)
	}

	if _, err = f.WriteString(m.StringFormat + "\n"); err != nil {
		log.Fatal("Unable to write to file:", err)
	}

	addValueToCache(path, m.ParsedSpeed)
}

func WriteInfoByDate(m m.ByDateMsg, wr io.Writer) {
	files := getMatchingFiles(fp.Join(dbDir, m.Year, m.Month, m.Day), m.ParsedSpeed)
	dirName := fp.Join(dbDir, m.Year, m.Month, m.Day)
	writFilesToResponseByFileNumbers(files, dirName, wr)
}

func WriteExtremesByDate(m m.ExtremesMsg, wr io.Writer) {
	files := getFilesWithExtremes(fp.Join(dbDir, m.Year, m.Month, m.Day))
	dirName := fp.Join(dbDir, m.Year, m.Month, m.Day)
	writFilesToResponseByFileNumbers(files, dirName, wr)

	if len(files) == 1 {
		_, _ = fmt.Fprintf(wr, "За %v.%v.%v указанная скорось является единственной зафиксированной (mix и max)",
			m.Day, m.Month, m.Year)
	}
}

func writFilesToResponseByFileNumbers(files []int, dirName string, wr io.Writer) {
	if files == nil || len(files) == 0 {
		_, _ = fmt.Fprintln(wr, "Данных с указанным параметрами нет")
	}

	w := bufio.NewWriter(wr)

	for _, fn := range files {

		fileName := strconv.Itoa(fn)
		f, err := os.Open(fp.Join(dirName, fileName))

		if err != nil {
			log.Fatal("Unable to open file:", err)
		}

		r := bufio.NewReader(f)
		buf := make([]byte, 1024)

		for {
			n, err := r.Read(buf)

			if err != nil && err != io.EOF {
				panic(err)
			}

			if n == 0 {
				break
			}

			if _, err := wr.Write(buf[:n]); err != nil {
				panic(err)
			}
		}

		if err = w.Flush(); err != nil {
			log.Fatal(err)
		}

		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

func addValueToCache(path string, value int) {
	if values, ok := c.Cache.Load(path); ok {
		values = insertSort(values, value)
		c.Cache.Store(path, values)
	} else {
		initializeFilesForFolder(path)
		addValueToCache(path, value)
	}
}

func insertSort(data []int, el int) []int {
	if i, ok := contains(data, el); !ok {
		data = append(data, 0)
		copy(data[i+1:], data[i:])
		data[i] = el
	}
	return data
}

func contains(s []int, si int) (i int, ok bool) {
	i = sort.SearchInts(s, si)
	ok = i < len(s) && s[i] == si
	return
}

func getMatchingFiles(path string, min int) []int {
	if values, ok := c.Cache.Load(path); ok {
		if len(values) > 0 {
			i, res := searchValueOrNearest(values, min, 0, len(values)-1)
			if res {
				return values[i:]
			}
		}
	} else {
		initializeFilesForFolder(path)
		return getMatchingFiles(path, min)
	}
	return nil
}

func getFilesWithExtremes(path string) []int {
	if values, ok := c.Cache.Load(path); ok {
		namesAsNumbers := make([]int, 0, 2)
		if len(values) > 1 {
			namesAsNumbers = append(namesAsNumbers, values[0], values[len(values)-1])
			return namesAsNumbers
		} else if len(values) == 1 {
			namesAsNumbers = append(namesAsNumbers, values[0])
		}
		return namesAsNumbers
	} else {
		initializeFilesForFolder(path)
		return getFilesWithExtremes(path)
	}
}

func initializeFilesForFolder(path string) {
	args := getDirFiles(path)
	sort.Ints(args)
	c.Cache.Store(path, args)
}

func searchValueOrNearest(data []int, target int, low int, high int) (index int, found bool) {
	mid := (high + low) / 2
	if low > high {
		if len(data) <= 2 {
			if data[0] >= target {
				index = 0
				found = true
				return
			} else if data[1] >= target {
				index = 1
				found = true
				return
			}
		}

		if mid != len(data)-1 {
			if target < data[0] {
				index = mid
				found = true
			} else {
				index = mid + 1
				found = true
			}
		} else {
			index = -1
			found = false
		}
	} else {
		if target < data[mid] {
			index, found = searchValueOrNearest(data, target, low, mid-1)
		} else if target > data[mid] {
			index, found = searchValueOrNearest(data, target, mid+1, high)
		} else if target == data[mid] {
			index = mid
			found = true
		}
	}
	return
}

func getDirFiles(path string) (namesAsNumbers []int) {
	_ = ensureDir(path)
	files, err := ioutil.ReadDir(path)

	if err != nil {
		log.Fatal(err)
	}

	namesAsNumbers = make([]int, 0, 8)
	for _, f := range files {
		val, err := strconv.Atoi(f.Name())
		if err == nil {
			namesAsNumbers = append(namesAsNumbers, val)
		}
	}
	return namesAsNumbers
}

func ensureDir(path string) (err error) {
	err = os.MkdirAll(path, os.ModeDir)
	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}
