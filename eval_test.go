package eval_test

import (
	"fmt"
	"github.com/graphikDB/eval"
	"testing"
)

func Test(t *testing.T) {
	decision, err := eval.NewDecision(eval.AllTrue, "this.name == 'bob'")
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := decision.Eval(map[string]interface{}{
		"name":  "bob",
		"email": "bob@acme.com",
	}); err != nil {
		t.Fatal(err.Error())
	}
	if err := decision.Eval(map[string]interface{}{
		"name":  "bob3",
		"email": "bob@acme.com",
	}); err == nil {
		t.Fatal("expected an error since bob3 != bob")
	}
	trigg, err := eval.NewTrigger(decision, "{'name': 'coleman'}")
	if err != nil {
		t.Fatal(err.Error())
	}
	person := map[string]interface{}{
		"name":  "bob",
		"email": "bob@acme.com",
	}
	data, err := trigg.Trigger(person)
	if err != nil {
		t.Fatal(err.Error())
	}
	if data["name"] != "coleman" {
		t.Fatal("failed to trigger")
	}
	fmt.Println("trigger expressions: ", trigg.Expression())
}

func ExampleNewDecision() {
	decision, err := eval.NewDecision(eval.AllTrue, "this.email.endsWith('acme.com')")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := decision.Eval(map[string]interface{}{
		"name":  "bob",
		"email": "bob@acme.com",
	}); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(decision.Expression())
	// Output: this.email.endsWith('acme.com')
}

func ExampleNewTrigger() {
	decision, err := eval.NewDecision(eval.AllTrue, "this.email.endsWith('acme.com')")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	trigg, err := eval.NewTrigger(decision, `
	{
		'admin': true,
		'updated_at': now()
	}
`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	person := map[string]interface{}{
		"name":  "bob",
		"email": "bob@acme.com",
	}
	data, err := trigg.Trigger(person)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(data["admin"], data["updated_at"].(int64) > 0)
	// Output: true true
}
