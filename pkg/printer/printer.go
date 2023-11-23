package printer

import (
	"bytes"
	"celify/pkg/evaluator"
	"celify/pkg/helpers"
	"celify/pkg/models"
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/fatih/color"
)

type Printer struct {
	Evaluator *evaluator.Evaluator
}

func NewPrinter(evaluator *evaluator.Evaluator) *Printer {
	return &Printer{
		Evaluator: evaluator,
	}
}

func (p *Printer) PrintResults(results []models.EvaluationResult, supressObjects bool) {
	fmt.Println()
	for _, result := range results {
		color.New(color.Bold).Add(color.Underline).Printf("validation \"%s\":\n", result.Expression)
		if result.ValidationError != nil {
			fmt.Printf("%s %s\n", getErrorStr(), color.YellowString(result.ValidationError.Error()))
			fmt.Printf("%s\n", getErrorStr())
			if !supressObjects {
				printEvaluatedObjects(result.EvaluatedObjects, p.Evaluator.TargetData.Format)
			}
			continue
		}
		color.New(color.FgGreen).Println("Success: true")
		fmt.Println()
	}
}

func getErrorStr() string {
	return color.New(color.FgRed).Sprint("|")
}

func PrintMultilineError(input string, color *color.Color) {
	errLines := strings.Split(input, "\n")
	for _, line := range errLines {
		fmt.Printf("%s %s\n", getErrorStr(), color.Sprint(line))
	}
}

func printEvaluatedObjects(objects []models.EvaluatedObject, format string) {
	for _, obj := range objects {
		byteObj, err := helpers.MarshalData(obj.Object, format)
		if err != nil {
			fmt.Printf("%s %s\n", getErrorStr(), color.New(color.FgRed).Sprint("Error marshalling object"))
			return
		}
		strObj := string(byteObj)
		fmt.Printf("%s %s %s\n", getErrorStr(), color.New(color.Underline).Sprint("object:"), obj.Expression)
		if format == "yaml" {
			colorizedYaml := colorizeYaml(strObj)
			PrintMultilineError(colorizedYaml, color.New(color.Reset))
		} else if format == "json" {
			colorizedJson := colorizeJson(strObj)
			PrintMultilineError(colorizedJson, color.New(color.Reset))
		}
	}
}

func FmtError(err error) error {
	fmt.Println()
	summaryStr := color.New(color.FgRed).Add(color.Underline).Add(color.Bold).Sprint("Error Summary:")
	errStr := color.New(color.FgRed).Sprint(err)
	return fmt.Errorf("%s\n%s", summaryStr, errStr)
}

func colorizeJson(jsonStr string) string {
	// return string(pretty.Color(pretty.Pretty([]byte(jsonStr)), nil))
	lexer := lexers.Get("json")
	style := styles.Get("solarized-light")
	formatter := formatters.Get("terminal")

	iterator, err := lexer.Tokenise(nil, jsonStr)
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	writer := io.Writer(&b)
	formatter.Format(writer, style, iterator)
	return b.String()
}

func colorizeYaml(yamlStr string) string {
	lexer := lexers.Get("yaml")
	style := styles.Get("solarized-light")
	formatter := formatters.Get("terminal")

	iterator, err := lexer.Tokenise(nil, yamlStr)
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	writer := io.Writer(&b)
	formatter.Format(writer, style, iterator)
	return b.String()
}
