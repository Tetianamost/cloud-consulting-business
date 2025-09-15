#!/bin/bash

echo "Running Email Service Event Integration Test..."
cd "$(dirname "$0")"
go run test_email_service_event_integration.go