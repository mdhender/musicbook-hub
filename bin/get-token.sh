#!/bin/bash

# Copyright (c) 2025 Michael D Henderson. All rights reserved.

set -e

PORT=8181
URL="http://localhost:$PORT"

UUID=$(jq -r '.[0]' magic-keys.json)

if [[ -z "$UUID" ]]; then
  echo "❌ Could not read UUID from magic-keys.json"
  exit 1
fi

echo "🔐 Requesting token for UUID: $UUID"
RESPONSE=$(curl -s "$URL/api/login/$UUID")

TOKEN=$(echo "$RESPONSE" | jq -r '.token')

if [[ -z "$TOKEN" || "$TOKEN" == "null" ]]; then
  echo "❌ Failed to retrieve token"
  echo "$RESPONSE"
  exit 1
fi

echo "✅ Token received:"
echo "$TOKEN"
echo

echo "🔍 Calling /api/me with token..."
curl -s -H "Authorization: Bearer $TOKEN" "$URL/api/me" | jq .


echo
echo "📚 Fetching all books (including private)..."
curl -s -H "Authorization: Bearer $TOKEN" "$URL/api/books" | jq .
