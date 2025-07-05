package main

import (
	"fmt"
	"iter"
	"strings"

	"github.com/arielsrv/go-seq"
)

func main() {
	fmt.Println("=== Example of Where, SelectMany and Select (Functional API) ===")

	people := []string{
		"John,Engineer,25",
		"Mary,Doctor,30",
		"Charles,Engineer,28",
		"Anna,Teacher,35",
		"Louis,Doctor,32",
		"Helen,Engineer,27",
		"Robert,Teacher,40",
	}

	// 1. Filter only engineers
	// 2. Expand into skills
	// 3. Convert to uppercase
	result := seq.Select(
		seq.SelectMany(
			seq.Where(
				seq.Yield(people...),
				func(p string) bool {
					return strings.Contains(strings.ToLower(p), "engineer")
				},
			),
			func(p string) iter.Seq[string] {
				parts := strings.Split(p, ",")
				name := parts[0]
				profession := parts[1]
				age := parts[2]
				skills := []string{
					name + "-" + profession,
					name + "-" + age + " years",
					name + "-Expert in " + profession,
				}
				return seq.Yield(skills...)
			},
		),
		func(s string) string {
			return strings.ToUpper(s)
		},
	)

	fmt.Println("\n1. Result as slice:")
	slice := seq.ToSlice(result)
	for i, item := range slice {
		fmt.Printf("  %d: %s\n", i+1, item)
	}

	fmt.Println("\n2. Result as map (index -> value):")
	mapResult := make(map[int]string)
	for i, item := range slice {
		mapResult[i] = item
	}
	for k, v := range mapResult {
		fmt.Printf("  %d: %s\n", k, v)
	}

	fmt.Println("\n3. Result using SelectKeys (name -> skill):")
	resultWithKeys := seq.SelectKeys(
		seq.SelectMany(
			seq.Where(
				seq.Yield(people...),
				func(p string) bool {
					return strings.Contains(strings.ToLower(p), "engineer")
				},
			),
			func(p string) iter.Seq[string] {
				parts := strings.Split(p, ",")
				name := parts[0]
				profession := parts[1]
				age := parts[2]
				skills := []string{
					name + "-" + profession,
					name + "-" + age + " years",
					name + "-Expert in " + profession,
				}
				return seq.Yield(skills...)
			},
		),
		func(s string) string {
			parts := strings.Split(s, "-")
			return parts[0]
		},
	)
	mapWithKeys := seq.ToMap(resultWithKeys)
	for k, v := range mapWithKeys {
		fmt.Printf("  %s: %s\n", k, v)
	}

	fmt.Println("\n4. Advanced example with multiple transformations:")
	advancedData := []string{
		"ProjectA,Development,High",
		"ProjectB,Testing,Medium",
		"ProjectC,Development,High",
		"ProjectD,Design,Low",
	}

	advancedResult := seq.Select(
		seq.SelectMany(
			seq.Where(
				seq.Yield(advancedData...),
				func(p string) bool {
					return strings.Contains(p, "Development") && strings.Contains(p, "High")
				},
			),
			func(p string) iter.Seq[string] {
				parts := strings.Split(p, ",")
				project := parts[0]
				tasks := []string{
					project + "-Analysis",
					project + "-Coding",
					project + "-Testing",
					project + "-Documentation",
				}
				return seq.Yield(tasks...)
			},
		),
		func(t string) string {
			return "[HIGH] " + t
		},
	)

	sliceAdvanced := seq.ToSlice(advancedResult)
	mapAdvanced := make(map[string]string)
	for i, task := range sliceAdvanced {
		mapAdvanced[fmt.Sprintf("Task_%d", i+1)] = task
	}

	fmt.Println("  Resulting slice:")
	for i, item := range sliceAdvanced {
		fmt.Printf("    %d: %s\n", i+1, item)
	}

	fmt.Println("  Resulting map:")
	for k, v := range mapAdvanced {
		fmt.Printf("    %s: %s\n", k, v)
	}
}
