package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func readFromFile(path string) ([]float64, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")

	var data []float64
	for _, line := range lines {
		if line == "" {
			continue
		}
		num, err := strconv.ParseFloat(strings.TrimSpace(line), 64)
		if err == nil {
			data = append(data, num)
		}
	}
	return data, nil
}

func average(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func median(data []float64) float64 {

	mid := len(data) / 2
	if len(data)%2 == 0 {
		return (data[mid-1] + data[mid]) / 2
	}
	return data[mid]
}

func variance(data []float64, avg float64) float64 {
	var sum float64
	for _, v := range data {
		diff := v - avg
		sum += diff * diff
	}
	return sum / float64(len(data))
}

func stdDev(variance float64) float64 {
	return math.Sqrt(variance)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go math.txt")
		return
	}

	data, err := readFromFile(os.Args[1])
	if err != nil || len(data) == 0 {
		fmt.Println("Error reading data")
		return
	}

	sort.Float64s(data)

	avg := average(data)
	med := median(data)
	vari := variance(data, avg)
	std := stdDev(vari)

	fmt.Printf("Average: %d\n", int(math.Round(avg)))
	fmt.Printf("Median: %d\n", int(math.Round(med)))
	fmt.Printf("Variance: %d\n", int(math.Round(vari)))
	fmt.Printf("Standard Deviation: %d\n", int(math.Round(std)))
}
