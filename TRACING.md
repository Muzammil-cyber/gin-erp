# Request Tracing Guide

## Overview

The ERP system now includes comprehensive request/response tracing with detailed logging to help debug issues and monitor system behavior.

## Features

### 1. **Automatic Log Files**

- Logs are automatically created in `logs/` directory
- One log file per day: `logs/app-2026-01-03.log`
- Logs contain both console output and detailed request/response data

### 2. **Trace ID**

Every request gets a unique Trace ID that appears in:

- HTTP Response Header: `X-Trace-ID`
- Error responses in JSON: `"trace_id": "..."`
- All log entries for that request

### 3. **Detailed Logging**

Each request logs:

- **Request Details**: Method, Path, Client IP, Timestamp
- **Request Headers**: Authorization, Content-Type, etc.
- **Request Body**: Up to 500 characters (for debugging)
- **Response Details**: Status Code, Duration, Size
- **Response Body**: Up to 500 characters (for debugging)
- **Errors**: Any errors that occurred

## Usage

### View Live Logs

```bash
make logs
```

### Trace a Specific Request

```bash
make trace ID=09e7f0fd-e523-4681-a743-f73a764b52ca
```

Or use the script directly:

```bash
./scripts/trace.sh 09e7f0fd-e523-4681-a743-f73a764b52ca
```

### View All Logs

```bash
make logs-all
```

### Search Logs Manually

```bash
# Search for a trace ID
grep "09e7f0fd-e523-4681-a743-f73a764b52ca" logs/*.log

# Search for errors
grep "Status: 500" logs/*.log

# Search for slow requests (>1s)
grep -E "Duration: [0-9]+\.[0-9]+s" logs/*.log

# Search for specific endpoint
grep "/api/v1/auth/profile" logs/*.log
```

## Log Format Example

```
═══════════════════════════════════════════════════════════════
▶ REQUEST | TraceID: 09e7f0fd-e523-4681-a743-f73a764b52ca
  GET /swagger/doc.json
  Client: 127.0.0.1
  Time: 2026-01-03 00:22:13
  Content-Type: application/json
  Body: {"email":"test@example.com","password":"***"}
◀ RESPONSE | TraceID: 09e7f0fd-e523-4681-a743-f73a764b52ca
  Status: 200
  Duration: 2.381898ms
  Size: 1024 bytes
  Body: {"success":true,"data":{...}}
═══════════════════════════════════════════════════════════════
```

## Debugging Workflow

1. **User Reports Issue** - Get the error response which includes `trace_id`

2. **Search Logs** - Use `make trace ID=<trace_id>` to find all logs for that request

3. **Analyze Flow**:
   - Check request headers (is auth token present?)
   - Check request body (is data valid?)
   - Check response status (what went wrong?)
   - Check response body (error details)
   - Check duration (performance issue?)

4. **Find Root Cause** - The detailed logs show exactly what happened at each step

## Best Practices

### For Developers

- Always include the trace ID when logging errors in your code
- Use structured logging for important operations
- Keep request/response bodies reasonable in size

### For Operations

- Rotate logs daily (already configured)
- Archive old logs after 30 days
- Monitor log file sizes
- Set up alerts for 500 errors

### For Support

- Always ask users for the trace ID from error responses
- Use trace ID to quickly find relevant logs
- Include trace ID when escalating issues

## Security Notes

⚠️ **Important**: Request/response bodies are logged (limited to 500 chars)

- Passwords are not hidden in logs (they should be in headers/tokens, not logged)
- Consider adding sensitive data filtering for production
- Restrict access to log files
- Do not share logs publicly without sanitizing

## Configuration

To disable response body logging in production, update `logger_middleware.go`:

```go
// Skip body logging in production
if os.Getenv("GIN_MODE") != "release" {
    // Log body code here
}
```

## Troubleshooting

### Logs not appearing

- Check if `logs/` directory exists
- Verify file permissions (logs need write access)
- Ensure `middleware.InitLogger()` is called in main.go

### Trace ID not found

- Verify the trace ID is correct
- Check if the request actually reached your server
- Ensure logs are being written (check file size)

### Log files too large

- Implement log rotation (see OPERATIONS.md)
- Reduce body logging size limit
- Filter out noisy endpoints (health checks, metrics)

## Integration with Monitoring Tools

### ELK Stack

```bash
# Ship logs to Logstash
filebeat -c filebeat.yml
```

### CloudWatch

```bash
# AWS CloudWatch Logs Agent
aws logs put-log-events --log-group-name erp-system
```

### Datadog

```bash
# Datadog Agent
datadog-agent logs tail logs/*.log
```

## Future Enhancements

- [ ] Add log rotation by size
- [ ] Add log compression for old files
- [ ] Filter sensitive data automatically
- [ ] Add performance metrics logging
- [ ] Add database query logging
- [ ] Export traces to distributed tracing systems (Jaeger, Zipkin)
