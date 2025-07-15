#!/bin/bash

echo "Testing complete sample creation with measurements and indices..."

# Login first
echo "Logging in..."
curl -s -X POST http://localhost:2342/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "jcz", "password": "password"}' \
  -c complete_cookies.txt

# Test creating a sample with complete measurements for index calculation
echo "Creating sample with complete measurements..."
curl -s -X POST http://localhost:2342/api/samples \
  -H "Content-Type: application/json" \
  -d '{
    "pool_id": 1,
    "sample_datetime": "2025-07-14T15:00",
    "kit_id": 1,
    "notes": "Complete sample with all measurements for index calculation",
    "measurements": {
      "ph": 7.2,
      "fc": 2.0,
      "tc": 2.5,
      "ta": 120.0,
      "ch": 200.0,
      "temperature": 80.0,
      "cya": 30.0,
      "salinity": 3200.0
    }
  }' \
  -b complete_cookies.txt

echo ""
echo "Retrieving samples to verify creation with measurements and indices..."
curl -s -X GET http://localhost:2342/api/samples -b complete_cookies.txt | python3 -m json.tool

# Clean up
rm -f complete_cookies.txt

echo ""
echo "Test completed."