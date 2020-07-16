package goson

import (
	"encoding/json"
	"testing"
)

func TestMakeZeroOmitempty(t *testing.T) {
	type Child struct {
		ChildA      string `json:"childA"`
		ChildB      string `json:"childB,omitempty"`
		ChildStruct struct {
			A int `json:"A,omitempty"`
			B int `json:"B"`
		} `json:"childStruct"`
	}
	type arr struct {
		Arr1 string `json:"arr1"`
		Arr2 string `json:"arr2,omitempty"`
	}

	type obj struct {
		Name       string `json:"name"`
		Additional string `json:"additional,omitempty"`
		AddStruct  struct {
			Field1 string `json:"field1,omitempty"`
			Field2 string `json:"field2"`
		} `json:"addStruct"`
		AddChildPtr *Child  `json:"child"`
		Slice       []arr   `json:"slice"`
		Nullable    *string `json:"nullable,omitempty"`
	}

	/*
		create example object with nested type & struct, json string from value in this object is:
		{
			"name": "Test",
			"additional": "test lagi (value in this field will removed because contains omitempty in json tag)",
			"addStruct": {
				"field1": "Field 1 (value in this field will removed because contains omitempty in json tag)",
				"field2": "Field 2"
			},
			"child": {
				"childA": "Child A",
				"childB": "Child B (value in this field will removed because contains omitempty in json tag)",
				"childStruct": {
					"A": 453, // value in this field will removed because contains omitempty in json tag
					"B": 567
				}
			},
			"slice": [
				{
					"arr1": "Test field slice 1",
					"arr2": "Test field slice 2 (value in this field will removed because contains omitempty in json tag)"
				}
			]
		}

		result after make zero field contains omitempty is:
		{
			"name": "Test",
			"addStruct": {
				"field2": "Field 2"
			},
			"child": {
				"childA": "Child A",
				"childStruct": {
					"B": 567
				}
			},
			"slice": [
				{
					"arr1": "Test field slice 1"
				}
			]
		}
	*/
	objExample := obj{
		Name:       "Test",
		Additional: "test lagi (value in this field will removed because contains omitempty in json tag)",
		AddStruct: struct {
			Field1 string `json:"field1,omitempty"`
			Field2 string `json:"field2"`
		}{
			Field1: "Field 1 (value in this field will removed because contains omitempty in json tag)",
			Field2: "Field 2",
		},
		AddChildPtr: &Child{
			ChildA: "Child A",
			ChildB: "Child B (value in this field will removed because contains omitempty in json tag)",
			ChildStruct: struct {
				A int `json:"A,omitempty"`
				B int `json:"B"`
			}{
				A: 453,
				B: 567,
			},
		},
		Slice: []arr{{
			Arr1: "Test field slice 1",
			Arr2: "Test field slice 2 (value in this field will removed because contains omitempty in json tag)",
		},
		},
	}

	tests := []struct {
		name       string
		args       interface{}
		wantResult string
		wantErr    bool
	}{
		{
			name:       "Testcase #1",
			args:       &objExample,
			wantResult: `{"name":"Test","addStruct":{"field2":"Field 2"},"child":{"childA":"Child A","childStruct":{"B":567}},"slice":[{"arr1":"Test field slice 1"}]}`,
		},
		{
			name:       "Testcase #2, argument is slice",
			args:       []*obj{&objExample},
			wantResult: `[{"name":"Test","addStruct":{"field2":"Field 2"},"child":{"childA":"Child A","childStruct":{"B":567}},"slice":[{"arr1":"Test field slice 1"}]}]`,
		},
		{
			name:       "Testcase #3, argument is map of slice",
			args:       map[string][]*obj{"aa": {&objExample}},
			wantResult: `{"aa":[{"name":"Test","addStruct":{"field2":"Field 2"},"child":{"childA":"Child A","childStruct":{"B":567}},"slice":[{"arr1":"Test field slice 1"}]}]}`,
		},
		{
			name:       "Testcase #4, argument is map of object include struct",
			args:       map[string]*obj{"aa": &objExample},
			wantResult: `{"aa":{"name":"Test","addStruct":{"field2":"Field 2"},"child":{"childA":"Child A","childStruct":{"B":567}},"slice":[{"arr1":"Test field slice 1"}]}}`,
		},
		{
			name:    "Testcase #6 want panic, invalid data type argument",
			args:    obj{Name: "Test", Additional: "test lagi"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MakeZeroOmitempty(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("error is unexpected, expect error is %v, error is: %v", tt.wantErr, err)
			}

			jsonAfter, _ := json.Marshal(tt.args)
			if tt.wantResult != "" && string(jsonAfter) != tt.wantResult {
				t.Errorf("expected %s, got %s", tt.wantResult, string(jsonAfter))
			}
		})
	}
}
