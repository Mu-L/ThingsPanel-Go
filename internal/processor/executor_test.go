package processor

import (
	"context"
	"encoding/json"
	"testing"
)

func TestExecuteEncodeAllowsLegacyOSTime(t *testing.T) {
	const script = `
function encodeInp(msg,topic)
    local json = require("json")
    local originCmd = json.decode(msg)

    local output = string.format(
        '{"method":"%s","params":{"params":%s,"id":"%s","version":"1.0","method":"thing.service.property.set"}}',
        originCmd.method,
        json.encode(originCmd.params),
        tostring(os.time())
    )

    return output
end`

	result, err := NewLuaExecutor().ExecuteEncode(
		context.Background(),
		script,
		[]byte(`{"method":"CO2","params":{"CO2":2223}}`),
	)
	if err != nil {
		t.Fatalf("legacy script execution failed: %v", err)
	}

	var output struct {
		Method string `json:"method"`
		Params struct {
			ID string `json:"id"`
		} `json:"params"`
	}
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if output.Method != "CO2" {
		t.Fatalf("unexpected method: %q", output.Method)
	}
	if output.Params.ID == "" {
		t.Fatal("os.time() returned an empty id")
	}
}

func TestSandboxDoesNotExposeDangerousOSFunctions(t *testing.T) {
	const script = `
function encodeInp(msg,topic)
    return tostring(os.execute == nil and os.remove == nil and os.getenv == nil)
end`

	result, err := NewLuaExecutor().ExecuteEncode(context.Background(), script, []byte(`{}`))
	if err != nil {
		t.Fatalf("sandbox check failed: %v", err)
	}
	if result != "true" {
		t.Fatalf("dangerous os functions are exposed: %q", result)
	}
}
