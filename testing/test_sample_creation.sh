#!/bin/bash

echo "Testing sample creation with datetime-local format..."

# Login first
echo "Logging in..."
curl -s -X POST http://localhost:2342/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "jcz", "password": "password"}' \
  -c sample_cookies.txt

# Test creating a sample with datetime-local format
echo "Creating sample with datetime-local format..."
curl -s -X POST http://localhost:2342/api/samples \
  -H "Content-Type: application/json" \
  -d '{
    "pool_id": 1,
    "sample_datetime": "2025-07-14T14:22",
    "kit_id": 1,
    "notes": "Test sample with datetime-local format",
    "measurements": {
      "ph": 7.2,
      "fc": 2.0,
      "tc": 2.5,
      "temperature": 80.0
    }
  }' \
  -b sample_cookies.txt

echo ""
echo "Retrieving samples to verify creation..."
curl -s -X GET http://localhost:2342/api/samples -b sample_cookies.txt

# Clean up
rm -f sample_cookies.txt

echo ""
echo "Test completed."