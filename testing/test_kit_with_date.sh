#!/bin/bash

echo "Testing kit creation with date parsing fix..."

# Login first
echo "Logging in..."
curl -s -X POST http://localhost:2342/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "jcz", "password": "password"}' \
  -c test_cookies.txt

# Test creating a kit with date-only format
echo "Creating kit with date-only format..."
curl -s -X POST http://localhost:2342/api/kits \
  -H "Content-Type: application/json" \
  -d '{"name": "Taylor K-2006", "description": "Test kit for chlorine and pH", "purchased_date": "2025-07-01"}' \
  -b test_cookies.txt

# Clean up
rm -f test_cookies.txt

echo "Test completed."