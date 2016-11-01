// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may not
// use this file except in compliance with the License. A copy of the
// License is located at
//
// http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Package parameterstore contains modules to resolve ssm parameters present in the document.
package parameterstore

import (
	"testing"

	"github.com/aws/amazon-ssm-agent/agent/contracts"
	"github.com/aws/amazon-ssm-agent/agent/log"
	"github.com/stretchr/testify/assert"
)

type StringTestCase struct {
	Input             string
	Output            string
	Parameters        []Parameter
	InvalidParameters []string
}

type StringListTestCase struct {
	Input             []string
	Output            []string
	Parameters        []Parameter
	InvalidParameters []string
}

var StringTestCases = []StringTestCase{
	StringTestCase{
		Input:             "This is a test string",
		Output:            "This is a test string",
		Parameters:        []Parameter{},
		InvalidParameters: []string{},
	},
	StringTestCase{
		Input:  "This is a {{ssm:test}} string",
		Output: "This is a testvalue string",
		Parameters: []Parameter{
			{
				Name:  "test",
				Type:  "String",
				Value: "testvalue",
			},
		},
		InvalidParameters: []string{},
	},
}

var StringListTestCases = []StringListTestCase{
	StringListTestCase{
		Input:             []string{"This is a test string", "Another test string"},
		Output:            []string{"This is a test string", "Another test string"},
		Parameters:        []Parameter{},
		InvalidParameters: []string{},
	},
	StringListTestCase{
		Input:  []string{"This is a {{ssm:test}} string", "Another parameter {{ ssm:foo }}"},
		Output: []string{"This is a testvalue string", "Another parameter randomvalue"},
		Parameters: []Parameter{
			{
				Name:  "test",
				Type:  "String",
				Value: "testvalue",
			},
			{
				Name:  "foo",
				Type:  "String",
				Value: "randomvalue",
			},
		},
		InvalidParameters: []string{},
	},
}

var logger = log.NewMockLog()

func TestResolve(t *testing.T) {
	testResolveMethod(t, StringTestCases[0])
	testResolveMethod(t, StringTestCases[1])
}

func testResolveMethod(t *testing.T, testCase StringTestCase) {
	callParameterService = func(
		log log.T,
		paramNames []string) (*GetParametersResponse, error) {
		result := GetParametersResponse{}
		result.Parameters = testCase.Parameters
		result.InvalidParameters = testCase.InvalidParameters
		return &result, nil
	}

	result, err := Resolve(logger, testCase.Input, false)

	assert.Equal(t, testCase.Output, result)
	assert.Nil(t, err)
}

func TestResolveSecureString(t *testing.T) {
	testResolveSecureStringMethod(t, StringTestCases[0])
	testResolveSecureStringMethod(t, StringTestCases[1])
}

func testResolveSecureStringMethod(t *testing.T, testCase StringTestCase) {
	callParameterService = func(
		log log.T,
		paramNames []string) (*GetParametersResponse, error) {
		result := GetParametersResponse{}
		result.Parameters = testCase.Parameters
		result.InvalidParameters = testCase.InvalidParameters
		return &result, nil
	}

	result, err := ResolveSecureString(logger, testCase.Input)

	assert.Equal(t, testCase.Output, result)
	assert.Nil(t, err)
}

func TestResolveSecureStringForStringList(t *testing.T) {
	testResolveSecureStringForStringListMethod(t, StringListTestCases[0])
	testResolveSecureStringForStringListMethod(t, StringListTestCases[1])
}

func testResolveSecureStringForStringListMethod(t *testing.T, testCase StringListTestCase) {
	callParameterService = func(
		log log.T,
		paramNames []string) (*GetParametersResponse, error) {
		result := GetParametersResponse{}
		result.Parameters = testCase.Parameters
		result.InvalidParameters = testCase.InvalidParameters
		return &result, nil
	}

	result, err := ResolveSecureStringForStringList(logger, testCase.Input)

	assert.Equal(t, testCase.Output, result)
	assert.Nil(t, err)
}

func TestValidateSSMParameters(t *testing.T) {
	var documentParameters = map[string]*contracts.Parameter{
		"commands": &contracts.Parameter{
			AllowedPattern: "^[a-zA-Z0-9]+$",
		},
		"workingDirectory": &contracts.Parameter{
			AllowedPattern: "",
		},
	}

	parameters := map[string]interface{}{
		"commands":         "test",
		"workingDirectory": "",
	}

	err := ValidateSSMParameters(logger, documentParameters, parameters)
	assert.Nil(t, err)
}
