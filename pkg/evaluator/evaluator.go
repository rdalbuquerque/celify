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
	targetData map[string]interface{}
	env        *cel.Env
}

func NewEvaluator(targetInput map[string]interface{}) (*Evaluator, error) {
	// Setup CEL environment
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("object", decls.NewMapType(decls.String, decls.Dyn)),
		),
	)
	if err != nil {
		return nil, err
	}
	targetData := map[string]interface{}{
		"object": targetInput,
	}
	return &Evaluator{
		targetData: targetData,
		env:        env,
	}, nil
}

func (ev *Evaluator) EvaluateSingleExpression(expression string) (any, error) {
	pgr, err := ev.getProgram(expression)
	if err != nil {
		return nil, errors.Errorf("Error generating program: %v", err)
	}
	out, _, err := pgr.Eval(ev.targetData)
	if err != nil {
		return nil, errors.Errorf("Error evaluating program: %v", err)
	}
	return out.Value(), nil
}

func (ev *Evaluator) Evaluate(validations []models.ValidationRule) (any, error) {
	var result any
	var err error
	for _, validation := range validations {
		result, err = ev.EvaluateSingleExpression(validation.Expression)
		if err != nil {
			return false, errors.Errorf("Error evaluating expression '%s': %v", validation.Expression, err)
		}
		if !result.(bool) {
			errMsg, err := ev.EvaluateSingleExpression(validation.ErrorMessage)
			if err != nil {
				return false, errors.Errorf("Error evaluating error message expression '%s': %v", validation.ErrorMessage, err)
			}
			fmt.Println(helpers.GetErrorStr(errMsg.(string)))
			ev.printEvaluatedObject(validation.Expression)
		}
	}
	return result, nil
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

func (ev *Evaluator) printEvaluatedObject(expression string) error {
	obj, err := ev.getEvaluatedObject(expression)
	if err != nil {
		return errors.Errorf("Error evaluating object expression '%s': %v", expression, err)
	}
	objStr, err := helpers.MarshalData(obj)
	if err != nil {
		return errors.Errorf("Error marshalling object: %v", err)
	}
	helpers.PrintEvaluatedObject(string(objStr))
	return nil
}

func (ev *Evaluator) getEvaluatedObject(expression string) (any, error) {
	objExpr := helpers.ExtractObject(expression)
	obj, err := ev.EvaluateSingleExpression(objExpr)
	if err != nil {
		return nil, errors.Errorf("Error evaluating object expression '%s': %v", objExpr, err)
	}
	return obj, nil
}
