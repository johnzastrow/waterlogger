#!/bin/bash

# Test script to verify kit creation functionality

echo "Testing kit creation functionality..."

# First, try to log in as the admin user
echo "Step 1: Logging in as admin user..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:2342/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "jcz", "password": "password"}' \
  -c cookies.txt)

echo "Login response: $LOGIN_RESPONSE"

# Test if login was successful by checking if cookies were set
if [ -f cookies.txt ]; then
    echo "Step 2: Testing kit creation..."
    
    # Create a test kit
    KIT_RESPONSE=$(curl -s -X POST http://localhost:2342/api/kits \
      -H "Content-Type: application/json" \
      -d '{"name": "Taylor K-2006 Test Kit", "description": "Complete test kit for pH, chlorine, and alkalinity", "purchased_date": "2025-07-14T00:00:00Z"}' \
      -b cookies.txt)
    
    echo "Kit creation response: $KIT_RESPONSE"
    
    # Test getting all kits
    echo "Step 3: Retrieving all kits..."
    KITS_RESPONSE=$(curl -s -X GET http://localhost:2342/api/kits -b cookies.txt)
    echo "Kits list response: $KITS_RESPONSE"
    
else
    echo "Login failed - no cookies set"
fi

# Clean up
rm -f cookies.txt

echo "Test completed."