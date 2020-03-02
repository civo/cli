package cmd

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func outputTable(headers []string, data [][]string) {
	if OutputFormat == "custom" {
		for _, items := range data {
			output := OutputFields
			for index, name := range headers {
				if strings.Contains(output, name) {
					output = strings.Replace(output, name, items[index], 1)
				}
			}
			output = strings.Replace(output, "\\t", "\t", -1)
			output = strings.Replace(output, "\\n", "\n", -1)
			fmt.Println(output)
		}
	} else if OutputFormat == "table" || OutputFormat == "human" || OutputFormat == "" {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(headers)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoFormatHeaders(false)
		table.AppendBulk(data)
		table.Render()
	}
}

func outputKeyValue(data map[string]string) {
	if OutputFormat == "custom" {
		output := OutputFields
		for key, value := range data {
			if strings.Contains(output, key) {
				output = strings.Replace(output, key, value, 1)
			}
			output = strings.Replace(output, "\\t", "\t", -1)
			output = strings.Replace(output, "\\n", "\n", -1)
			fmt.Println(output)
		}
	} else if OutputFormat == "table" || OutputFormat == "human" || OutputFormat == "" {
		table := tablewriter.NewWriter(os.Stdout)
		for key, value := range data {
			table.Append([]string{key, value})
		}
		table.Render()
	}
}

func mapToStringKeys(data reflect.Value) ([]string, error) {
	var keys []string

	i := 0
	iter := data.MapRange()
	for iter.Next() {
		k := iter.Key().Interface().(string)
		keys = append(keys, k)
		i++
	}

	return keys, nil
}

func findPartialKey(search string, data interface{}) (string, error) {
	keys, err := mapToStringKeys(reflect.ValueOf(data))
	if err != nil {
		return "", err
	}

	var result string

	for _, k := range keys {
		if strings.Contains(k, search) {
			if result == "" {
				result = k
			} else {
				return "", fmt.Errorf("Unable to find %s because there were multiple matches\n", search)
			}
		}
	}

	if result == "" {
		return "", fmt.Errorf("Unable to find %s at all in the list\n", search)
	}

	return result, nil
}
