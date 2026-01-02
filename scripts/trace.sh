#!/bin/bash

# Trace script for finding request/response flow by TraceID

if [ -z "$1" ]; then
    echo "Usage: ./trace.sh <trace-id>"
    echo "Example: ./trace.sh 09e7f0fd-e523-4681-a743-f73a764b52ca"
    exit 1
fi

TRACE_ID=$1
LOG_DIR="logs"

echo "════════════════════════════════════════════════════════════════"
echo "🔍 Searching for TraceID: $TRACE_ID"
echo "════════════════════════════════════════════════════════════════"
echo ""

# Check if logs directory exists
if [ ! -d "$LOG_DIR" ]; then
    echo "❌ Logs directory not found: $LOG_DIR"
    exit 1
fi

# Find all log files and search for the trace ID
found=0
for logfile in "$LOG_DIR"/*.log; do
    if [ -f "$logfile" ]; then
        matches=$(grep -c "$TRACE_ID" "$logfile" 2>/dev/null || echo "0")
        if [ "$matches" -gt 0 ]; then
            echo "📁 Found $matches occurrences in: $(basename $logfile)"
            echo "────────────────────────────────────────────────────────────────"
            # Show 15 lines before and after to capture full request/response block
            grep -B 15 -A 15 "$TRACE_ID" "$logfile" --color=always
            echo ""
            found=1
        fi
    fi
done

if [ $found -eq 0 ]; then
    echo "❌ No logs found for TraceID: $TRACE_ID"
    echo ""
    echo "💡 Tips:"
    echo "   - Make sure the server is running and has processed requests"
    echo "   - Check if the TraceID is correct"
    echo "   - Logs are stored in: $LOG_DIR/"
else
    echo "════════════════════════════════════════════════════════════════"
    echo "✅ Trace complete"
    echo "════════════════════════════════════════════════════════════════"
fi
