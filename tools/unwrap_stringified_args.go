package tools

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// UnwrapStringifiedToolArgumentsMiddleware runs before the MCP tool handler. Some MCP
// clients send array or object parameters as JSON strings. The go-sdk's jsonschema
// validation unmarshals top-level tool arguments into map[string]any, so string values
// fail validation for properties that are typed as arrays (e.g. app card fields, bulk items).
// This middleware parses those strings back into real JSON values before validation.
func UnwrapStringifiedToolArgumentsMiddleware(next mcp.MethodHandler) mcp.MethodHandler {
	return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
		if method == "tools/call" {
			if r, ok := req.(*mcp.ServerRequest[*mcp.CallToolParamsRaw]); ok && r != nil && r.Params != nil {
				fixed, err := unwrapStringifiedInRawJSON(r.Params.Arguments)
				if err == nil {
					r.Params.Arguments = fixed
				}
			}
		}
		return next(ctx, method, req)
	}
}

func unwrapStringifiedInRawJSON(raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 {
		return raw, nil
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return raw, err
	}
	unwrapArgMapInPlace(m)
	out, err := json.Marshal(m)
	if err != nil {
		return raw, err
	}
	return out, nil
}

// unwrapArgMapInPlace fixes stringified values for known keys and recurses into maps and slices.
func unwrapArgMapInPlace(m map[string]any) {
	for k, v := range m {
		m[k] = unwrapArgValue(k, v)
	}
}

func unwrapArgValue(key string, v any) any {
	if s, ok := v.(string); ok {
		if key == "fields" || key == "items" {
			if parsed, ok := tryParseJSONTop(s); ok {
				v = parsed
			}
		}
	}
	switch t := v.(type) {
	case map[string]any:
		unwrapArgMapInPlace(t)
		return t
	case []any:
		for i, el := range t {
			if sub, ok := el.(map[string]any); ok {
				unwrapArgMapInPlace(sub)
				t[i] = sub
			} else {
				t[i] = el
			}
		}
		return t
	default:
		return v
	}
}

func tryParseJSONTop(s string) (any, bool) {
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return nil, false
	}
	if s[0] != '[' && s[0] != '{' {
		return nil, false
	}
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return nil, false
	}
	return v, true
}
