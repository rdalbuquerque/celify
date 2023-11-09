package evaluator

import (
	"celify/pkg/helpers"
	"celify/pkg/models"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/pkg/errors"
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

func (ev *Evaluator) executeEvaluation(expression string) (interface{}, error) {
	pgr, err := ev.getProgram(expression)
	if err != nil {
		return nil, errors.Errorf("Error getting program: %v", err)
	}
	out, _, err := pgr.Eval(ev.TargetData.Data)
	if err != nil {
		return nil, errors.Errorf("Error evaluating expression: %v", err)
	}
	return out.Value(), nil
}

func (ev *Evaluator) EvaluateRule(rule models.ValidationRule) models.EvaluationResult {
	result, err := ev.executeEvaluation(rule.Expression)
	if err != nil {
		return models.EvaluationResult{
			ValidationError: errors.Errorf("Error evaluating expression: %v", err),
		}
	}
	objStr := helpers.ExtractObject(rule.Expression)
	evalObj, err := ev.executeEvaluation(objStr)
	if err != nil {
		return models.EvaluationResult{
			ValidationError: errors.Errorf("Error evaluating object expression '%s': %v", objStr, err),
		}
	}
	var boolResult *bool
	var ok bool
	*boolResult, ok = result.(bool)
	if !ok {
		boolResult = nil
	}
	if !*boolResult {
		msgExpr, err := ev.executeEvaluation(rule.ErrorMessage)
		if err != nil {
			return models.EvaluationResult{
				ValidationError: errors.Errorf("Error evaluating message expression '%s': %v", rule.ErrorMessage, err),
			}
		}
		return models.EvaluationResult{
			ValidationResult:  boolResult,
			EvaluatedObject:   evalObj,
			FailedRule:        rule.Expression,
			MessageExpression: msgExpr.(string),
		}
	}
	return models.EvaluationResult{
		ValidationResult: boolResult,
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
