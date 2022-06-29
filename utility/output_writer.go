// Package utility provides utility functionality that is used throughout
// the Civo CLI.
package utility

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type byLen []string

func (a byLen) Len() int {
	return len(a)
}
func (a byLen) Less(i, j int) bool {
	return len(a[i]) > len(a[j])
}
func (a byLen) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// OutputWriter is for printing structured data in various
// tabular formats
//
//   ow := utility.NewOutputWriter()
//   ow.StartLine()
//   ow.AppendData("ID", instance.ID)
//
//   # Then one of:
//   ow.WriteSingleObjectJSON()
//   ow.WriteMultipleObjectsJSON()
//   ow.WriteCustomOutput(common.OutputFields)
//   ow.WriteKeyValues()
//   ow.WriteTable()
type OutputWriter struct {
	Keys       []string
	Labels     []string
	Values     [][]string
	TempValues []string
}

// NewOutputWriter builds a new OutputWriter
func NewOutputWriter() *OutputWriter {
	ret := OutputWriter{}
	return &ret
}

// NewOutputWriterWithMap builds a new OutputWriter and automatically
// inserts the supplied map as a single line
func NewOutputWriterWithMap(data map[string]string) *OutputWriter {
	ow := OutputWriter{}
	ow.StartLine()

	for k, v := range data {
		ow.AppendData(k, v)
	}

	return &ow
}

// ToJSON is a function to show the output in json format
func (ow *OutputWriter) ToJSON(v interface{}, pretty bool) {
	value, _ := json.Marshal(v)

	switch pretty {
	case true:
		result, _ := prettyprint(value)
		fmt.Println(string(result))
	default:
		fmt.Println(string(value))
	}
}

// StartLine starts a new line of output
func (ow *OutputWriter) StartLine() {
	ow.finishExistingLine()
	ow.TempValues = make([]string, len(ow.Keys))
}

func (ow *OutputWriter) finishExistingLine() {
	if len(ow.TempValues) > 0 {
		ow.Values = append(ow.Values, ow.TempValues)
	}
}

// AppendDataWithLabel adds a line of data to the output writer
func (ow *OutputWriter) AppendDataWithLabel(key, value, label string) {
	found := -1
	for i, v := range ow.Keys {
		if v == key {
			found = i
		}
	}

	if found == -1 {
		ow.Keys = append(ow.Keys, key)
		ow.Labels = append(ow.Labels, label)
		ow.TempValues = append(ow.TempValues, value)
	} else {
		ow.TempValues[found] = value
	}
}

// AppendData adds a line of data to the output writer
func (ow *OutputWriter) AppendData(key, value string) {
	ow.AppendDataWithLabel(key, value, key)
}

// WriteSingleObjectJSON writes the JSON for a single object to STDOUT
func (ow *OutputWriter) WriteSingleObjectJSON(pretty bool) {
	ow.finishExistingLine()

	data := map[string]string{}

	for i, k := range ow.Keys {
		if ow.Values[0][i] != "" {
			data[k] = ow.Values[0][i]
		}
	}

	jsonString, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	switch pretty {
	case true:
		result, _ := prettyprint(jsonString)
		fmt.Println(string(result))
	default:
		fmt.Println(string(jsonString))
	}

}

// WriteMultipleObjectsJSON writes the JSON for multiple objects to STDOUT
func (ow *OutputWriter) WriteMultipleObjectsJSON(pretty bool) {
	ow.finishExistingLine()

	data := make([]map[string]string, len(ow.Values))
	for i, row := range ow.Values {
		dataRow := map[string]string{}
		for col, k := range ow.Keys {
			if row[col] != "" {
				dataRow[k] = row[col]
			}
		}
		data[i] = dataRow
	}

	jsonString, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	switch pretty {
	case true:
		result, _ := prettyprint(jsonString)
		fmt.Println(string(result))
	default:
		fmt.Println(string(jsonString))
	}

}

// WriteKeyValues prints a single object stored in the OutputWriter
// in key: value format
func (ow *OutputWriter) WriteKeyValues() {
	ow.finishExistingLine()

	longestLabelLength := 0
	for _, label := range ow.Labels {
		if len(label) > longestLabelLength {
			longestLabelLength = len(label)
		}
	}

	for i := range ow.Keys {
		value := ow.Values[0][i]
		label := ow.Labels[i]
		fmt.Printf("%"+strconv.Itoa(longestLabelLength)+"s : %s\n", label, value)
	}
}

// WriteTable prints multiple objects stored in the OutputWriter
// in tabular format
func (ow *OutputWriter) WriteTable() {
	ow.finishExistingLine()

	table := tablewriter.NewWriter(os.Stdout)
	if len(ow.Keys) > 0 {
		table.SetHeader(ow.Labels)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(false)
	} else {
		table.SetBorder(false)
	}

	table.AppendBulk(ow.Values)
	table.Render()
}

// Replace the nth occurrence of old in s by new.
func replaceNth(s, old, new string, n int) string {
	i := 0
	for m := 1; m <= n; m++ {
		x := strings.Index(s[i:], old)
		if x < 0 {
			break
		}
		i += x
		if m == n {
			return s[:i] + new + s[i+len(old):]
		}
		i += len(old)
	}
	return s
}

// WriteCustomOutput prints one or multiple objects using custom formatting
func (ow *OutputWriter) WriteCustomOutput(fields string) {
	ow.finishExistingLine()
	defaultKeys := make([]string, len(ow.Keys))
	copy(defaultKeys, ow.Keys)
	sort.Sort(byLen(ow.Keys))

	// Transcribe Data
	customMap := make(map[string][]string)
	for i, key := range defaultKeys {
		for _, row := range ow.Values {
			customMap[key] = append(customMap[key], row[i])
		}
	}

	for i := range ow.Values {
		output := fields
		for key, name := range ow.Keys {
			var re = regexp.MustCompile(name)
			if len(re.FindStringIndex(output)) > 0 {
				output = replaceNth(output, name, fmt.Sprintf("$%v$", key), 1)
			}
		}

		for index, name := range ow.Keys {
			if strings.Contains(output, fmt.Sprintf("$%v$", index)) {
				output = strings.Replace(output, fmt.Sprintf("$%v$", index), customMap[name][i], 1)
			}
		}
		output = strings.Replace(output, "\\t", "\t", -1)
		output = strings.Replace(output, "\\n", "\n", -1)
		fmt.Println(output)
	}
}

// WriteSubheader writes a centred heading line in to output
func (ow *OutputWriter) WriteSubheader(label string) {
	count := (72 - len(label) + 2) / 2
	fmt.Println(strings.Repeat("-", count) + " " + label + " " + strings.Repeat("-", count))
}

// WriteHeader WriteSubheader writes a centred heading line in to output
func (ow *OutputWriter) WriteHeader(label string) {
	fmt.Printf("%s:\n", label)
}

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
