#!/bin/bash

echo "Testing measurements creation fix..."

# Login
curl -s -X POST http://localhost:2342/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "jcz", "password": "password"}' \
  -c test_cookies.txt

# Create a sample with measurements
echo "Creating sample with measurements..."
RESULT=$(curl -s -X POST http://localhost:2342/api/samples \
  -H "Content-Type: application/json" \
  -d '{
    "pool_id": 1,
    "sample_datetime": "2025-07-14T19:30",
    "kit_id": 1,
    "user_id": 1,
    "notes": "Test sample for measurements fix",
    "measurements": {
      "ph": 7.4,
      "fc": 1.5,
      "tc": 1.5,
      "ta": 110.0,
      "ch": 180.0,
      "temperature": 78.0
    }
  }' \
  -b test_cookies.txt)

echo "Result: $RESULT"

# Check if sample was created successfully
if echo "$RESULT" | grep -q "error"; then
    echo "❌ Sample creation failed"
else
    echo "✅ Sample creation succeeded"
    
    # Check if measurements and indices were created
    echo "Checking sample data..."
    curl -s -X GET http://localhost:2342/api/samples -b test_cookies.txt | python3 -c "
import json
import sys
data = json.load(sys.stdin)
if data:
    sample = data[-1]  # Get the last sample
    print(f'Sample ID: {sample[\"id\"]}')
    print(f'Notes: {sample[\"notes\"]}')
    print(f'Has measurements: {\"measurements\" in sample}')
    if 'measurements' in sample:
        print(f'  pH: {sample[\"measurements\"][\"ph\"]}')
        print(f'  FC: {sample[\"measurements\"][\"fc\"]}')
        print(f'  Temperature: {sample[\"measurements\"][\"temperature\"]}')
    print(f'Has indices: {\"indices\" in sample}')
    if 'indices' in sample and sample['indices']:
        print(f'  LSI: {sample[\"indices\"].get(\"lsi\", \"None\")}')
        print(f'  RSI: {sample[\"indices\"].get(\"rsi\", \"None\")}')
    else:
        print('  No indices calculated')
else:
    print('No samples found')
"
fi

# Clean up
rm -f test_cookies.txt
echo "Test completed."