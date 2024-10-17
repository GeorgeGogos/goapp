# 2024/10/15

Initial version.

# 2024/10/16

## Problem #1: WebSocket Session Statistics Fix
- **File**: `internal/pkg/httpsrv/handler_websocket.go`, `internal/pkg/httpsrv/server.go`, `internal/pkg/httpsrv/stats.go`
  - Fixed the issue with counting only a single message per WebSocket session by properly incrementing counters for each received message.

## Problem #2: Minimize Memory Usage for WS Sessions
- **File**: `internal/pkg/httpsrv/routes.go`
  - Ensured proper cleanup of resources by closing channels and terminating goroutines cleanly to prevent memory leaks in WebSocket sessions.
 
## New Feature #1: Hexadecimal String Generator
- **Files**: `internal/pkg/strgen/strgen.go`, `internal/pkg/strgen/strgen_test.go`
  - Updated the random string generator to generate only hexadecimal values. Added unit tests (`TestRandHexString`) and a benchmark (`BenchmarkRandHexString`) to verify the accuracy and performance.

# 2024/10/17

## Problem #3: Add CSRF Token Verification
- **File**: `internal/pkg/httpsrv/watcher.go`, `internal/pkg/watcher/watcher.go`
  - Added a CSRF middleware (`CSRFTokenMiddleware`) to validate incoming requests for a CSRF token. Updated handler wrappers to enforce CSRF token verification for all routes.

## New Feature #2: Return Hex Value in WS Connection
- **File**: `internal/pkg/httpsrv/handler_websocket.go`, `pkg/util/string.go`
  - Modified WebSocket handler to return hexadecimal values along with iteration counts, enhancing WebSocket response content.
