#!/bin/bash

# Generate schema.sql from migration files
# This script consolidates all up migrations into a single schema file

SCHEMA_FILE="db/schema.sql"
MIGRATIONS_DIR="db/migrations"

echo "Generating schema from migrations..."

# Create schema file header
cat > "$SCHEMA_FILE" << 'EOF'
-- Schema generated from migration files
-- Run this script to regenerate: ./scripts/generate-schema.sh
-- DO NOT EDIT MANUALLY - Generated from migration files

EOF

# Concatenate all .up.sql migration files in order
for migration in "$MIGRATIONS_DIR"/*.up.sql; do
    if [ -f "$migration" ]; then
        echo "-- From: $(basename "$migration")" >> "$SCHEMA_FILE"
        cat "$migration" >> "$SCHEMA_FILE"
        echo "" >> "$SCHEMA_FILE"
    fi
done

echo "Schema generated successfully: $SCHEMA_FILE"