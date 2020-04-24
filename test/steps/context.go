// Copyright 2019 Red Hat, Inc. and/or its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package steps

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/cucumber/godog"
)

// registerContextSteps register all existing Context steps
func registerContextSteps(s *godog.Suite, data *Data) {
	s.Step(`^Check variable "([^"]*)" is an array of size (\d+)$`, data.checkVariableIsAnArrayOfSize)
	s.Step(`^Parse variable "([^"]*)" as a map and set key which value matches "([^"]*)" into variable "([^"]*)"$`, data.parseVariableAsMapAndGetKeyFromValue)
	s.Step(`^Parse variable "([^"]*)" with expression "([^"]*)" into variable "([^"]*)"$`, data.parseVariableWithExpression)
}

func (data *Data) checkVariableIsAnArrayOfSize(contextKey string, size int) error {
	var array []interface{}
	if err := data.parseScenarioVariableAsJSON(contextKey, &array); err != nil {
		return err
	}

	if len(array) != size {
		return fmt.Errorf("Array is not of expected size. Actual: %d, Expected: %d", len(array), size)
	}

	return nil
}

func (data *Data) parseVariableAsMapAndGetKeyFromValue(contextKey, value, outContextKey string) error {
	var collection map[string]interface{}
	if err := data.parseScenarioVariableAsJSON(contextKey, &collection); err != nil {
		return err
	}

	for actualKey, actualValue := range collection {
		if actualValue == value {
			data.ScenarioContext[outContextKey] = actualKey
			return nil
		}
	}

	return fmt.Errorf("Value '%s' not found in map", value)
}

func (data *Data) parseVariableWithExpression(contextKey, expression, outContextKey string) error {
	eval, _ := jsonpath.Language().NewEvaluable(expression)

	var jsonValue interface{}
	if err := data.parseScenarioVariableAsJSON(contextKey, &jsonValue); err != nil {
		return err
	}

	val, err := eval(context.Background(), jsonValue)
	if err != nil {
		return err
	}

	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("Error parsing value to string. Value: %v", val)
	}

	data.ScenarioContext[outContextKey] = str

	return nil
}

func (data *Data) parseScenarioVariableAsJSON(contextKey string, v interface{}) error {
	contextValue := data.ScenarioContext[contextKey]
	return json.Unmarshal([]byte(contextValue), &v)
}
