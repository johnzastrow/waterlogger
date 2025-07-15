#!/bin/bash

echo "Testing indices calculation and display..."

# Start a fresh server
pkill -f waterlogger
sleep 2

# Start server in background
./waterlogger > indices_test.log 2>&1 &
sleep 3

# Login
curl -s -X POST http://localhost:2342/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "jcz", "password": "password"}' \
  -c indices_cookies.txt

# Create a sample with complete measurements that should calculate indices
echo "Creating sample with measurements for indices calculation..."
curl -s -X POST http://localhost:2342/api/samples \
  -H "Content-Type: application/json" \
  -d '{
    "pool_id": 1,
    "sample_datetime": "2025-07-14T19:00",
    "kit_id": 1,
    "user_id": 1,
    "notes": "Test sample for indices calculation",
    "measurements": {
      "ph": 7.2,
      "fc": 2.0,
      "tc": 2.0,
      "ta": 120.0,
      "ch": 200.0,
      "temperature": 80.0
    }
  }' \
  -b indices_cookies.txt

echo ""
echo "Checking if indices were calculated..."
curl -s -X GET http://localhost:2342/api/samples -b indices_cookies.txt | python3 -c "
import json
import sys
data = json.load(sys.stdin)
if data:
    sample = data[-1]  # Get the last sample
    print(f'Sample ID: {sample[\"id\"]}')
    print(f'Has measurements: {\"measurements\" in sample}')
    if 'measurements' in sample:
        print(f'pH: {sample[\"measurements\"][\"ph\"]}')
    print(f'Has indices: {\"indices\" in sample}')
    if 'indices' in sample:
        print(f'Indices: {sample[\"indices\"]}')
else:
    print('No samples found')
"

# Clean up
rm -f indices_cookies.txt
echo "Test completed. Check indices_test.log for debug output."