package printer

import (
	"celify/pkg/evaluator"
	"celify/pkg/helpers"
	"celify/pkg/models"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/tidwall/pretty"
)

type Printer struct {
	Evaluator *evaluator.Evaluator
}

func NewPrinter(evaluator *evaluator.Evaluator) *Printer {
	return &Printer{
		Evaluator: evaluator,
	}
}

func (p *Printer) PrintResults(results []models.EvaluationResult) {
	fmt.Println()
	for _, result := range results {
		color.New(color.Bold).Add(color.Underline).Printf("validation \"%s\":\n", result.Expression)
		if result.ValidationError != nil {
			fmt.Printf("%s %s\n", getErrorStr(), color.YellowString(result.ValidationError.Error()))
			printEvaluatedObject(result.EvaluatedObject, p.Evaluator.TargetData.Format)
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

func printEvaluatedObject(obj interface{}, format string) {
	byteObj, err := helpers.MarshalData(obj, format)
	if err != nil {
		fmt.Printf("%s %s\n", getErrorStr(), color.New(color.FgRed).Sprint("Error marshalling object"))
		return
	}
	strObj := string(byteObj)
	fmt.Printf("%s %s", getErrorStr(), color.New(color.Underline).Sprintln("Evaluated object:"))
	if format == "yaml" {
		c := color.New(color.FgBlue)
		PrintMultilineError(strObj, c)
	} else if format == "json" {
		PrintMultilineError(string(pretty.Color(pretty.Pretty(byteObj), nil)), color.New(color.Reset))
	}
}

func FmtError(err error) error {
	fmt.Println()
	summaryStr := color.New(color.FgRed).Add(color.Underline).Add(color.Bold).Sprint("Error Summary:")
	errStr := color.New(color.FgRed).Sprint(err)
	return fmt.Errorf("%s\n%s", summaryStr, errStr)
}
