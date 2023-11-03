package evaluator

import (
	"celify/pkg/helpers"
	"celify/pkg/models"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/pkg/errors"
)

type EvaluatorInterface interface {
	Evaluate(targetData map[string]interface{}, validations []models.ValidationRule) ([]any, error)
	EvaluateSingleExpression() (any, error)
}

type Evaluator struct {
	targetData *models.TargetData
	env        *cel.Env
}

func NewEvaluator(targetInput *models.TargetData) (*Evaluator, error) {
	// Setup CEL environment
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("object", decls.NewMapType(decls.String, decls.Dyn)),
		),
	)
	if err != nil {
		return nil, err
	}
	return &Evaluator{
		targetData: targetInput,
		env:        env,
	}, nil
}

func (ev *Evaluator) EvaluateSingleExpression(expression string) (any, error) {
	pgr, err := ev.getProgram(expression)
	if err != nil {
		return nil, errors.Errorf("Error generating program: %v\n", err)
	}
	out, _, err := pgr.Eval(ev.targetData.Data)
	if err != nil {
		return nil, errors.Errorf("Error evaluating program: %v\n", err)
	}
	return out.Value(), nil
}

func (ev *Evaluator) Evaluate(validations []models.ValidationRule) (any, error) {
	var result any
	var err error
	for _, validation := range validations {
		result, err = ev.EvaluateSingleExpression(validation.Expression)
		if err != nil {
			return false, errors.Errorf("Error evaluating expression '%s': %v\n", validation.Expression, err)
		}
		if !result.(bool) {
			errMsg, err := ev.EvaluateSingleExpression(validation.ErrorMessage)
			if err != nil {
				return false, errors.Errorf("Error evaluating error message expression '%s': %v\n", validation.ErrorMessage, err)
			}
			fmt.Println(helpers.GetErrorStr(errMsg.(string)))
			ev.printEvaluatedObject(validation.Expression, ev.targetData.Format)
		}
	}
	return result, nil
}

func (ev *Evaluator) getProgram(expression string) (cel.Program, error) {
	ast, issues := ev.env.Compile(expression)
	if issues != nil && issues.Err() != nil {
		return nil, errors.Errorf("Failed to compile expression '%s': %v\n", expression, issues.Err())
	}
	pgr, err := ev.env.Program(ast)
	if err != nil {
		return nil, errors.Errorf("Failed to generate program for expression '%s': %v\n", expression, err)
	}
	return pgr, nil
}

func (ev *Evaluator) printEvaluatedObject(expression, format string) error {
	obj, err := ev.getEvaluatedObject(expression)
	if err != nil {
		return errors.Errorf("Error evaluating object expression '%s': %v\n", expression, err)
	}
	objStr, err := helpers.MarshalData(obj, format)
	if err != nil {
		return errors.Errorf("Error marshalling object: %v\n", err)
	}
	helpers.PrintEvaluatedObject(string(objStr), format)
	return nil
}

func (ev *Evaluator) getEvaluatedObject(expression string) (any, error) {
	objExpr := helpers.ExtractObject(expression)
	obj, err := ev.EvaluateSingleExpression(objExpr)
	if err != nil {
		return nil, errors.Errorf("Error evaluating object expression '%s': %v\n", objExpr, err)
	}
	return obj, nil
}
