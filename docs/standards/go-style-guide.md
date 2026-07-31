# Go Coding Standards

These conventions are observed and enforced across this repo and its sibling
services. Follow them for any new or modified code — don't introduce a
different pattern for the same problem just because it's also valid Go.

## Architecture: Handler → Service → Store

- The HTTP layer (`handler.go`) only parses/validates request *shape* (JSON
  decode errors → `400`). It must not contain business rules.
- Business logic and field validation live in the Service layer
  (`service.go`) — trim input, check required fields, and return a sentinel
  error (see below) on failure. Handlers never validate field content
  themselves.
- Persistence sits behind a Store interface (`internal/store`), with an
  in-memory implementation used directly in tests (no mocking framework).
- One package per domain under `internal/` (e.g. `internal/users`,
  `internal/inventory`), each containing `model(s).go`, `service.go`,
  `handler.go`, `errors.go`, and matching `_test.go` files.

## Errors

- Declare package-level sentinel errors with `errors.New(...)` in a
  dedicated `errors.go` per domain package, e.g.:
  ```go
  var (
      ErrNotFound   = errors.New("not found")
      ErrValidation = errors.New("validation error")
      ErrInvalidID  = errors.New("invalid id")
  )
  ```
- Handlers translate service errors to HTTP status through a shared helper
  (`httpapi.WriteErrorFromService`) — never a per-handler `switch` on error
  type/message.
- Don't return a bare `fmt.Errorf(...)` for an expected/business-rule
  failure. Use a sentinel so callers can `errors.Is(err, ErrValidation)`
  instead of string-matching.

## Validation

- `strings.TrimSpace` every string input before validating or storing it.
- Reject empty required fields and obviously malformed values (e.g.
  `strings.Contains(email, "@")` for an email field) at the Service layer,
  returning `ErrValidation` — not in the handler, not in the store.

## HTTP responses

- All JSON output goes through shared helpers (`httpapi.WriteJSON`,
  `httpapi.WriteError`). Never call `json.NewEncoder(w).Encode(...)`
  directly inside a handler.
- Status codes: `201 Created` on successful create, `204 No Content` on
  successful delete, `200 OK` on get/update. Not-found conditions are
  surfaced via the shared error translator, not a hardcoded status per
  handler.

## Testing

- Standard library `testing` + `net/http/httptest` only — no testify,
  ginkgo, or other assertion/mocking libraries.
- Assertions are plain `if` checks with `t.Fatalf("...: %v", err)` /
  `t.Fatalf("expected %d, got %d: %s", want, got, body)` — not a
  third-party assertion library.
- Tests build a real in-memory store + service + router and exercise the
  full HTTP path via `router.ServeHTTP(w, req)`, rather than mocking the
  service or store layer — see any `*_handler_test.go` for the pattern.

## Package layout

- Shared HTTP plumbing (JSON encode/decode helpers, router construction,
  error-to-status translation) lives in `internal/httpapi` and is imported
  by every domain package. Don't duplicate response-writing logic inside a
  domain package.
