# Miro MCP: card / app_card coverage (fix report)

## P1 — `fields` rejected (jsonschema: string vs array)

**Diagnosis:** The failure happens inside `github.com/modelcontextprotocol/go-sdk` during `AddTool`’s `applySchema` step: request arguments are unmarshaled to `map[string]any` and validated with `github.com/google/jsonschema-go` against the schema derived from the Go `In` type (`jsonschema.ForType`).

If the **MCP client** sends `fields` (or `items`) as a **JSON string** (the literal characters `[{...}]` inside a quoted string), the map value has Go type `string`, so the validator reports: `has type "string", want one of "null, array"`. The struct tags and generated schema for `[]AppCardField` are correct; the problem is the wire format from certain clients, not a wrong schema for real JSON arrays.

**Where the fix lives:** In **this server**, not the SDK: a `ReceivingMiddleware` runs before `tools/call` and reparses `fields` and `items` when they appear as stringified JSON. Valid clients that send native JSON arrays are unchanged.

**If Inspector works but another bridge fails:** document the client as needing a fix or the workaround (proper JSON types). See `ERRORS.md`.

---

## P2 — Card description read, search, list

**Change:** `GetItemResult` adds `description`, `due_date`, and `assignee` (and maps them from Miro’s `data.*`). `SearchBoard` matches query against `content`, `title`, and `description`, and follows `/items` pagination (capped) so the first page is not the only one searched.

**Nice-to-have `miro_get_card`:** Not added; `miro_get_item` now carries the main card text fields for typical agent use.

---

## P3 — Native status / assignee / tags

**API:** Miro’s REST payload for cards includes `data.description`, `data.dueDate`, and `data.assignee` when present. “Status” as a separate UI pill may not always appear in REST; this is called out in `ERRORS.md`. Tags remain the supported structured workaround via the existing tag tools.

---

## P4 — Connectors unmarshal

**Change:** `ConnectorEndpoint.position` uses `ConnectorPoint` with custom JSON decoding so `x` / `y` can be numbers, numeric strings, or **percentage strings** (e.g. `50%` for a relative anchor on an item edge). Percentages are normalized to fractional `0..1` in `X`/`Y`, and the original token is kept in `RawX`/`RawY` (`x_raw`/`y_raw` in JSON).

---

## Validation

- `go test ./...` — pass  
- `go vet ./...` — pass  
- `make lint` — not run (golangci-lint not installed in this environment)

**Binary (darwin arm64):** build with  
`GOOS=darwin GOARCH=arm64 go build -o miro-mcp-server-darwin-arm64 .`  
(artifact can be attached to a release on your fork / PR branch.)

---

## Upstream

Suggested **four focused PRs** to `olgasafonova/miro-mcp-server` (or one branch with the four commits on this model): P1, P2, P3 (docs-only), P4.
