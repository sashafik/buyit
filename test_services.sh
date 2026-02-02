#!/bin/bash

BASE_URL="http://localhost:8080"

echo "1. Registering User..."
curl -s -X POST $BASE_URL/auth/register -d '{"username":"testuser", "password":"password123"}' | json_pp

echo -e "\n\n2. Logging in..."
TOKEN=$(curl -s -X POST $BASE_URL/auth/login -d '{"username":"testuser", "password":"password123"}' | jq -r .token)
echo "Token: $TOKEN"

echo -e "\n3. Listing Products..."
curl -s $BASE_URL/products | json_pp

echo -e "\n4. Creating Order (Buying 1 Laptop)..."
curl -s -X POST $BASE_URL/orders \
  -H "Authorization: $TOKEN" \
  -d '[{"productId":"1", "quantity":1}]' | json_pp

echo -e "\n5. Checking Orders..."
curl -s -X GET $BASE_URL/orders \
  -H "Authorization: $TOKEN" | json_pp

echo -e "\n6. Verifying Inventory (Laptop stock should be 9)..."
curl -s $BASE_URL/products/1 | json_pp
