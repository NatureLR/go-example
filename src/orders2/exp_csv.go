package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
)

func expToCSV(root, sid string) {
	dirs, err := filepath.Glob(path.Join(root, sid, "*"))
	assert(err)
	for _, dir := range dirs {
		func() {
			f, err := os.Create(dir + ".csv")
			assert(err)
			defer func() {
				f.Close()
				os.Remove(dir)
			}()
			cw := csv.NewWriter(f)
			cw.UseCRLF = true
			fs, err := filepath.Glob(path.Join(dir, "*.json"))
			assert(err)
			sort.Strings(fs)
			var headers []string
			for i, fn := range fs {
				func() {
					f, err := os.Open(fn)
					if err != nil {
						fmt.Printf("ERROR open %s: %v\n", fn, err)
						return
					}
					defer func() {
						f.Close()
						os.Remove(fn)
					}()
					dec := json.NewDecoder(f)
					var obj map[string]string
					err = dec.Decode(&obj)
					if err != nil {
						fmt.Printf("ERROR decode %s: %v\n", fn, err)
						return
					}
					if i == 0 {
						for k := range obj {
							headers = append(headers, k)
						}
						cw.Write(headers)
					}
					var row []string
					for _, h := range headers {
						row = append(row, obj[h])
					}
					cw.Write(row)
				}()
			}
			cw.Flush()
		}()
	}
}
