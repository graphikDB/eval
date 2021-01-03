package trigger

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
	"github.com/pkg/errors"
	"strings"
)

// Trigger creates values as map[string]interface{} if it's decisider returns no errors against a Mapper
type Trigger struct {
	decision   *Decision
	program    cel.Program
	expression string
}

// NewTrigger creates a new trigger instance from the decision & trigger expressions
func NewTrigger(decision *Decision, triggerExpression string) (*Trigger, error) {
	if triggerExpression == "" {
		return nil, ErrEmptyExpressions
	}
	program, err := globalEnv.Program(triggerExpression)
	if err != nil {
		return nil, err
	}
	return &Trigger{
		decision:   decision,
		program:    program,
		expression: triggerExpression,
	}, nil
}

// Trigger executes it's decision against the Mapper and then overwrites the
func (t *Trigger) Trigger(data map[string]interface{}) (map[string]interface{}, error) {
	if err := t.decision.Eval(data); err == nil {
		out, _, err := t.program.Eval(map[string]interface{}{
			"this": data,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "trigger: failed to evaluate trigger (%s)", t.expression)
		}
		if patchFields, ok := out.Value().(map[ref.Val]ref.Val); ok {
			newData := map[string]interface{}{}
			for k, v := range patchFields {
				newData[k.Value().(string)] = v.Value()
			}
			return newData, nil
		}
		if patchFields, ok := out.Value().(map[string]interface{}); ok {
			return patchFields, nil
		}
		if patchFields, ok := out.Value().(map[string]string); ok {
			newData := map[string]interface{}{}
			for k, v := range patchFields {
				newData[k] = v
			}
			return newData, nil
		}
		return map[string]interface{}{
			"value": out.Value(),
		}, nil
	}
	return map[string]interface{}{}, nil
}

// Expression returns the triggers raw CEL expressions
func (e *Trigger) Expression() string {
	return e.expression
}

const ArrowOperator = "=>"

var ErrArrowOperator = errors.Errorf("arrow operator: expecting syntax ${decision} %s ${mutation}", ArrowOperator)

// NewArrowTrigger creates a trigger from arrow syntax  ${decision} => ${mutation}
func NewArrowTrigger(arrowExpression string) (*Trigger, error) {
	split := strings.Split(arrowExpression, ArrowOperator)
	if len(split) != 2 {
		return nil, ErrArrowOperator
	}
	decisionExp := split[0]
	decision, err := NewDecision(decisionExp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create trigger from arrow expression")
	}
	triggerExp := split[1]
	t, err := NewTrigger(decision, triggerExp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create trigger from arrow expression")
	}
	return t, nil
}
