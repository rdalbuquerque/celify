package evaluator

import (
	"celify/pkg/helpers"
	"celify/pkg/models"
	"fmt"
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/pkg/errors"
)

var (
	StringType = reflect.TypeOf("")
	IntType    = reflect.TypeOf(0)
	BoolType   = reflect.TypeOf(true)
	AnyType    = reflect.TypeOf(new(interface{})).Elem()
)

type Evaluator struct {
	TargetData *models.TargetData
	env        *cel.Env
}

func NewEvaluator(targetInput *models.TargetData) (*Evaluator, error) {
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("object", decls.NewMapType(decls.String, decls.Dyn)),
		),
	)
	if err != nil {
		return nil, err
	}
	return &Evaluator{
		TargetData: targetInput,
		env:        env,
	}, nil
}

func (ev *Evaluator) executeEvaluation(expression string, expectedReturnType reflect.Type) (interface{}, error) {
	pgr, err := ev.getProgram(expression)
	if err != nil {
		return nil, fmt.Errorf("error getting program: %v", err)
	}
	out, _, err := pgr.Eval(ev.TargetData.Data)
	if err != nil {
		return nil, fmt.Errorf("error evaluating expression: %v", err)
	}
	return out.ConvertToNative(expectedReturnType)
}

func (ev *Evaluator) EvaluateRule(rule models.ValidationRule) models.EvaluationResult {
	result, err := ev.executeEvaluation(rule.Expression, BoolType)
	if err != nil {
		return models.EvaluationResult{
			Expression:      rule.Expression,
			ValidationError: fmt.Errorf("error evaluating expression: %w", err),
		}
	}
	boolResult := result.(bool)

	objExpr := helpers.ExtractObject(rule.Expression)
	evalObj, err := ev.executeEvaluation(objExpr, AnyType)
	if err != nil {
		evalObj = "unable to evaluate object"
	}

	if !boolResult && rule.MessageExpression != "" {
		msgExpr, err := ev.executeEvaluation(rule.MessageExpression, StringType)
		if err != nil {
			return models.EvaluationResult{
				Expression:      rule.Expression,
				ValidationError: fmt.Errorf("error evaluating message expression: %w", err),
			}
		}
		return models.EvaluationResult{
			Expression:        rule.Expression,
			ValidationResult:  &boolResult,
			EvaluatedObject:   evalObj,
			FailedRule:        rule.Expression,
			MessageExpression: msgExpr.(string),
		}
	}
	return models.EvaluationResult{
		Expression:       rule.Expression,
		ValidationResult: &boolResult,
		EvaluatedObject:  evalObj,
	}
}

func (ev *Evaluator) Evaluate(validations models.ValidationConfig) []models.EvaluationResult {
	var evalResults []models.EvaluationResult
	for _, validation := range validations.Validations {
		evalResults = append(evalResults, ev.EvaluateRule(validation))
	}
	return evalResults
}

func (ev *Evaluator) getProgram(expression string) (cel.Program, error) {
	ast, issues := ev.env.Compile(expression)
	if issues != nil && issues.Err() != nil {
		return nil, errors.Errorf("Failed to compile expression '%s': %v", expression, issues.Err())
	}
	pgr, err := ev.env.Program(ast)
	if err != nil {
		return nil, errors.Errorf("Failed to generate program for expression '%s': %v", expression, err)
	}
	return pgr, nil
}
