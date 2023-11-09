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
	for i, result := range results {
		color.New(color.Bold).Add(color.Underline).Printf("validation %d:\n", i+1)
		if result.ValidationError != nil {
			color.New(color.FgRed).Printf("Error: %v\n", result.ValidationError)
			continue
		}
		if result.ValidationResult == nil {
			color.New(color.FgHiYellow).Printf("%s result is not true or false\n", getErrorStr())
			continue
		}
		success := *result.ValidationResult
		if !success {
			color.New(color.FgHiYellow).Printf("%s %s\n", getErrorStr(), result.MessageExpression)
			printEvaluatedObject(result.EvaluatedObject, p.Evaluator.TargetData.Format)
		}
		color.New(color.FgGreen).Printf("Success: %v\n", success)
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
