package main

import (
	"bytes"
	"net/http/httptest"
	"reflect"
	"testing"
)

type TestCase struct {
	URL         string
	Type        string
	Description string
	Payload     string
	Answer      string
	Status      int
}

func TestCases(t *testing.T) {

	cases := []TestCase{

		TestCase{
			URL:         "http://hostname/builds",
			Type:        "POST",
			Description: "Initial case",
			Payload: `{
		    ru {
		        "Build base AMI": {
		            "Builds": [{
		                "runtime_seconds": "1931",
		                "build_date": "1506741166",
		                "result": "SUCCESS",
		                "output": "base-ami us-west-2 ami-9f0ae4e5 d1541c88258ccb3ee565fa1d2322e04cdc5a1fda"
		            }, {
		                "runtime_seconds": "1825",
		                "build_date": "1506740166",
		                "result": "SUCCESS",
		                "output": "base-ami us-west-2 ami-d3b92a92 3dd2e093fc75f0e903a4fd25240c89dd17c75d66"
		            }, {
		                "runtime_seconds": "126",
		                "build_date": "1506240166",
		                "result": "FAILURE",
		                "output": "base-ami us-west-2 ami-38a2b9c1 936c7725e69855f3c259c117173782f8c1e42d9a"
		            }, {
		                "runtime_seconds": "1842",
		                "build_date": "1506240566",
		                "result": "SUCCESS",
		                "output": "base-ami us-west-2 ami-91a42ed5 936c7725e69855f3c259c117173782f8c1e42d9a"
		            }, {
		                "runtime_seconds": "5",
		                "build_date": "1506250561",
		            }, {
		                "runtime_seconds": "215",
		                "build_date": "1506250826",
		                "result": "FAILURE",
		                "output": "base-ami us-west-2 ami-34a42e15 936c7725e69855f3c259c117173782f8c1e42d9a"
		            }]
		        }
		    }
		}`,
			Answer: `{"latest":{"build_date":"\"1506741166\"","ami_id":"\"ami-9f0ae4e5\"","commit_hash":"\"d1541c88258ccb3ee565fa1d2322e04cdc5a1fda\""}}`,
			Status: 200,
		},

		TestCase{
			URL:         "http://hostname/builds",
			Type:        "POST",
			Description: "Valid JSON",
			Payload: `{
		        "Build base AMI": {
		            "Builds": [{
		                "runtime_seconds": "1931",
		                "build_date": "1506741166",
		                "result": "SUCCESS",
		                "output": "base-ami us-west-2 ami-9f0ae4e5 d1541c88258ccb3ee565fa1d2322e04cdc5a1fda"
		            }, {
		                "runtime_seconds": "1825",
		                "build_date": "1506740166",
		                "result": "SUCCESS",
		                "output": "base-ami us-west-2 ami-d3b92a92 3dd2e093fc75f0e903a4fd25240c89dd17c75d66"
		            }, {
		                "runtime_seconds": "126",
		                "build_date": "1506240166",
		                "result": "FAILURE",
		                "output": "base-ami us-west-2 ami-38a2b9c1 936c7725e69855f3c259c117173782f8c1e42d9a"
		            }, {
		                "runtime_seconds": "1842",
		                "build_date": "1506240566",
		                "result": "SUCCESS",
		                "output": "base-ami us-west-2 ami-91a42ed5 936c7725e69855f3c259c117173782f8c1e42d9a"
		            }, {
		                "runtime_seconds": "5",
		                "build_date": "1506250561"
		            }, {
		                "runtime_seconds": "215",
		                "build_date": "1506250826",
		                "result": "FAILURE",
		                "output": "base-ami us-west-2 ami-34a42e15 936c7725e69855f3c259c117173782f8c1e42d9a"
		            }]
		        }
		}`,
			Answer: `{"latest":{"build_date":"\"1506741166\"","ami_id":"\"ami-9f0ae4e5\"","commit_hash":"\"d1541c88258ccb3ee565fa1d2322e04cdc5a1fda\""}}`,
			Status: 200,
		},

		TestCase{
			URL:         "http://hostname/builds",
			Type:        "GET",
			Description: "Get instead of Post",
			Payload:     `{}`,
			Answer:      `Only POST is allowed`,
			Status:      400,
		},

		TestCase{
			URL:         "http://hostname/builds",
			Type:        "POST",
			Description: "Not correct payload",
			Payload:     `{123}`,
			Answer:      `Can't process body: invalid character '1' looking for beginning of object key string`,
			Status:      400,
		},

		TestCase{
			URL:         "http://hostname/builds",
			Type:        "POST",
			Description: "Not correct data",
			Payload: `{
		    ru {
		        "Build base AMI": {
		            "Builds": [{
		                "runtime_seconds": "1931",
		                "build_date": "test",
		                "result": "SUCCESS",
		                "output": "base-ami us-west-2 ami-9f0ae4e5 d1541c88258ccb3ee565fa1d2322e04cdc5a1fda"
		            }]
		        }
		    }
		}`,
			Answer: `Can't process body: strconv.ParseInt: parsing "test": invalid syntax`,
			Status: 400,
		},

		TestCase{
			URL:         "http://hostname/builds",
			Type:        "POST",
			Description: "Without ami",
			Payload: `{
		    ru {
		        "Build base AMI": {
		            "Builds": [{
		                "runtime_seconds": "1931",
		                "build_date": "1506741166",
		                "result": "SUCCESS",
		                "output": "base-ami d1541c88258ccb3ee565fa1d2322e04cdc5a1fda"
		            }]
		        }
		    }
		}`,
			Answer: `Can't process body: Output is not correct`,
			Status: 400,
		},

		TestCase{
			URL:         "http://hostname/builds",
			Type:        "POST",
			Description: "No output",
			Payload: `{ ru {
		        "Build base AMI": {
		            "Builds": [{
		                "runtime_seconds": "5",
		                "build_date": "1506250561",
		            }]}}}`,
			Answer: `Can't process body: unexpected end of JSON input`,
			Status: 400,
		},
	}

	for caseNum, item := range cases {

		req := httptest.NewRequest(item.Type, item.URL, bytes.NewReader([]byte(item.Payload)))
		w := httptest.NewRecorder()

		getBuilds(w, req)

		if !reflect.DeepEqual(item.Status, w.Code) || !reflect.DeepEqual(item.Answer, w.Body.String()) {
			t.Errorf("[%d] %s ;; wrong result, expected %#v, got %#v, with code %#v, expected %#v", caseNum, item.Description, item.Answer, w.Body.String(), w.Code, item.Status)
		}
	}

}
