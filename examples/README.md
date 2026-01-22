# ECS Log Examples

This directory contains example ECS formatted logs showing the actual format from Cloud Foundry.

## Actual CF Log Format

Cloud Foundry logs with ECS JSON have this structure:

```
TIMESTAMP [SOURCE] OUT/ERR JSON_PAYLOAD
```

## Example 1: Application Log with DEBUG Level

**Input:**
```
2026-01-22T17:21:36.27+0000 [APP/REV/47/PROC/WEB/0] OUT {"@timestamp":"2026-01-22T17:21:36.274136327Z","log":{"level":"DEBUG","logger":"org.springframework.security.web.FilterChainProxy"},"process":{"pid":8,"thread":{"name":"http-nio-8080-exec-9"}},"service":{"name":"mcp-gateway","version":"0.10.2"},"message":"Securing GET /actuator/health","tags":["COMMONS-LOGGING"],"ecs":{"version":"8.11"}}
```

**Output:**
```
2026-01-22T17:21:36.274Z [APP/REV/47/PROC/WEB/0] DEBUG [FilterChainProxy] {http-nio-8080-exec-9} Securing GET /actuator/health
```

**Explanation:**
- **Timestamp**: Taken from ECS `@timestamp` field
- **Source**: `[APP/REV/47/PROC/WEB/0]` - Application, revision 47, web process, instance 0
- **Level**: `DEBUG` from `log.level`
- **Logger**: `[FilterChainProxy]` - shortened from full class name
- **Thread**: `{http-nio-8080-exec-9}` - from `process.thread.name`
- **Message**: The actual log message

## Example 2: Router Log (Non-ECS)

**Input:**
```
2026-01-22T17:21:29.61+0000 [RTR/2] OUT mcp-gateway-jira-prod.apps.tanzu.broadcom.net - [2026-01-22T17:21:29.599673994Z] "GET /actuator/prometheus HTTP/1.1" 200 0 60432 "-" "Go-http-client/1.1" "100.64.0.13:19220" "192.168.112.49:61022" x_forwarded_for:"100.64.0.13"...
```

**Output:**
```
2026-01-22T17:21:29.61+0000 [RTR/2] OUT mcp-gateway-jira-prod.apps.tanzu.broadcom.net - [2026-01-22T17:21:29.599673994Z] "GET /actuator/prometheus HTTP/1.1" 200 0 60432...
```

**Explanation:**
Router logs don't contain ECS JSON, so they're passed through unchanged. The `[RTR/2]` indicates router instance 2.

## Example 3: Multiple Logs Showing Thread Activity

**Input (from example-logs.txt):**
```
2026-01-22T17:21:36.27+0000 [APP/REV/47/PROC/WEB/0] OUT {"@timestamp":"2026-01-22T17:21:36.274136327Z","log":{"level":"DEBUG","logger":"org.springframework.security.web.FilterChainProxy"},"process":{"pid":8,"thread":{"name":"http-nio-8080-exec-9"}},"service":{"name":"mcp-gateway","version":"0.10.2"},"message":"Securing GET /actuator/health","tags":["COMMONS-LOGGING"],"ecs":{"version":"8.11"}}
2026-01-22T17:21:36.27+0000 [APP/REV/47/PROC/WEB/0] OUT {"@timestamp":"2026-01-22T17:21:36.274832410Z","log":{"level":"DEBUG","logger":"org.springframework.security.web.authentication.AnonymousAuthenticationFilter"},"process":{"pid":8,"thread":{"name":"http-nio-8080-exec-9"}},"service":{"name":"mcp-gateway","version":"0.10.2"},"message":"Set SecurityContextHolder to anonymous SecurityContext","tags":["COMMONS-LOGGING"],"ecs":{"version":"8.11"}}
2026-01-22T17:21:36.27+0000 [APP/REV/47/PROC/WEB/0] OUT {"@timestamp":"2026-01-22T17:21:36.275341406Z","log":{"level":"DEBUG","logger":"org.springframework.security.web.FilterChainProxy"},"process":{"pid":8,"thread":{"name":"http-nio-8080-exec-9"}},"service":{"name":"mcp-gateway","version":"0.10.2"},"message":"Secured GET /actuator/health","tags":["COMMONS-LOGGING"],"ecs":{"version":"8.11"}}
```

**Output:**
```
2026-01-22T17:21:36.274Z [APP/REV/47/PROC/WEB/0] DEBUG [FilterChainProxy] {http-nio-8080-exec-9} Securing GET /actuator/health
2026-01-22T17:21:36.274Z [APP/REV/47/PROC/WEB/0] DEBUG [AnonymousAuthenticationFilter] {http-nio-8080-exec-9} Set SecurityContextHolder to anonymous SecurityContext
2026-01-22T17:21:36.275Z [APP/REV/47/PROC/WEB/0] DEBUG [FilterChainProxy] {http-nio-8080-exec-9} Secured GET /actuator/health
```

**Explanation:**
These logs show a single request flowing through Spring Security filters, all on the same thread. The thread name makes it easy to track related log entries.

## ECS JSON Structure

The plugin supports the following ECS fields:

### Core Fields
- `@timestamp` - ISO 8601 timestamp
- `message` - The log message text
- `tags` - Array of tag strings

### Log Object
- `log.level` - Log level (DEBUG, INFO, WARN, ERROR)
- `log.logger` - Fully qualified logger name

### Process Object
- `process.pid` - Process ID
- `process.thread.name` - Thread name

### Service Object
- `service.name` - Service name
- `service.version` - Service version

### ECS Metadata
- `ecs.version` - ECS version number

### Additional Supported Fields
- `error.*` - Error details and stack traces
- `labels.*` - Metadata (app, org, space)
- `cloud.*`, `host.*` - Infrastructure info

## Output Format Template

```
TIMESTAMP [SOURCE] LEVEL [Logger] {thread} message
   [Optional Stack Trace]
   [Optional Tags]
```

### Components:
- **TIMESTAMP**: ECS @timestamp in format `2006-01-02T15:04:05.000Z07:00`
- **[SOURCE]**: From CF log prefix (e.g., `APP/REV/47/PROC/WEB/0`, `RTR/2`)
- **LEVEL**: Uppercase log level from `log.level`
- **[Logger]**: Short class name from `log.logger` (last component)
- **{thread}**: Thread name from `process.thread.name` (if present)
- **message**: The actual log message

## Real-World Usage

### Debugging a Request
```bash
cf ecslogs my-app | grep "http-nio-8080-exec-1"
```
Follow all logs from a specific thread to trace a single request through your application.

### Finding Errors
```bash
cf ecslogs my-app --recent | grep ERROR
```
Quickly find error-level logs.

### Monitoring Specific Components
```bash
cf ecslogs my-app | grep "\[FilterChainProxy\]"
```
Watch logs from a specific logger/component.

## File: example-logs.txt

The `example-logs.txt` file contains real Cloud Foundry logs from a Spring Boot application using ECS logging format. It includes:

- Application logs from web process (APP/REV/47/PROC/WEB/0)
- Router logs (RTR/2)
- Various log levels (DEBUG, INFO)
- Different Spring Framework components
- Thread information showing concurrent request processing

Use this file to test the plugin:
```bash
./bin/test-parser example-logs.txt
```
