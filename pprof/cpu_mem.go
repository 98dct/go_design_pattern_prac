package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
	"time"
)

type record struct {
	lon, lat float64
	temps    [12]float64
}

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to this file")
	memprofile = flag.String("memprofile", "", "write memory profile to this file")
)

func main() {

	go func() {
		http.ListenAndServe(":6060", nil)
	}()
	//dir, _ := os.Getwd()
	//fmt.Println(dir)
	//flag.Parse()
	//fmt.Println(*cpuprofile)
	//fmt.Println(*memprofile)
	//if *cpuprofile != "" {
	//	cpuf, err := os.Create(*cpuprofile)
	//	if err != nil{
	//		fmt.Println(err)
	//	}
	//	pprof.StartCPUProfile(cpuf)
	//	defer cpuf.Close()
	//	defer pprof.StopCPUProfile()
	//}
	//// ... some code
	//if *memprofile != "" {
	//	memf, err := os.Create(*memprofile)
	//	if err != nil{
	//		fmt.Println(err)
	//	}
	//	pprof.WriteHeapProfile(memf)
	//	memf.Close()
	//}

	st := time.Now()
	files, err := os.ReadDir("../../air_data")
	if err != nil {
		panic(err)
	}
	rawData := make(map[string][]byte)
	for _, file := range files {
		f, err := os.Open("../../air_data/" + file.Name())
		if err != nil {
			panic(err)
		}
		rawData[file.Name()], _ = io.ReadAll(f)
		f.Close()
	}
	parsed := parseData(rawData)
	output := make(map[string][]float64)
	for _, v := range parsed {
		for _, r := range v {
			key := fmt.Sprintf("%.3f:%.3f", r.lat, r.lon)
			output[key] = append(output[key], r.temps[:]...)
		}
	}
	outputFile, _ := os.Create("output.csv")
	w := csv.NewWriter(outputFile)
	defer w.Flush()
	defer outputFile.Close()
	for coordinates, yearlyTemps := range output {
		row := []string{coordinates, fmt.Sprintf("%.2f", average(yearlyTemps))}
		w.Write(row)
	}
	fmt.Println(time.Since(st))
}

func parseData(input map[string][]byte) map[string][]record {
	m := make(map[string][]record)
	for filename, v := range input {
		lines := strings.Split(string(v), "\n")
		for _, line := range lines {
			seg := strings.Fields(line)
			if len(seg) != 14 {
				continue
			}
			lon, _ := strconv.ParseFloat(seg[0], 64)
			lat, _ := strconv.ParseFloat(seg[1], 64)
			temps := [12]float64{}
			for i := 2; i < 14; i++ {
				t, _ := strconv.ParseFloat(seg[i], 64)
				temps[i-2] = t
			}
			m[filename] = append(m[filename], record{lon, lat, temps})
		}
	}
	return m
}

func average(input []float64) float64 {
	var sum float64
	for _, v := range input {
		sum += v
	}
	return sum / float64(len(input))
}
