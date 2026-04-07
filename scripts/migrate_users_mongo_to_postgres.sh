#!/bin/bash

# ============================================
# Migrate Users from MongoDB to PostgreSQL
# ============================================

echo "Starting MongoDB to PostgreSQL user migration..."

# Get MongoDB connection details
MONGO_HOST=${MONGO_HOST:-localhost}
MONGO_PORT=${MONGO_PORT:-27017}
MONGO_DB=${MONGO_DATABASE:-fiber_app}

# Get PostgreSQL connection details
PG_HOST=${POSTGRES_HOST:-localhost}
PG_PORT=${POSTGRES_PORT:-5432}
PG_USER=${POSTGRES_USER:-postgres}
PG_DB=${POSTGRES_DATABASE:-fiber_app}
PG_PASSWORD=${POSTGRES_PASSWORD:-postgres}

# Export users from MongoDB to JSON file
echo "Exporting users from MongoDB..."
mongosh --quiet --host "$MONGO_HOST:$MONGO_PORT" --db "$MONGO_DB" --eval "
  db.users.find().forEach(function(user) {
    print(JSON.stringify(user));
  });
" > /tmp/mongodb_users.json

# Check if we got any users
if [ ! -s /tmp/mongodb_users.json ]; then
  echo "No users found in MongoDB!"
  exit 1
fi

echo "Exported $(wc -l < /tmp/mongodb_users.json) users from MongoDB"

# Import users into PostgreSQL using Node.js script
cat > /tmp/import_users.js << 'EOF'
const fs = require('fs');
const { Client } = require('pg');

const client = new Client({
  host: process.env.PG_HOST || 'localhost',
  port: process.env.PG_PORT || 5432,
  database: process.env.PG_DB || 'fiber_app',
  user: process.env.PG_USER || 'postgres',
  password: process.env.PG_PASSWORD || 'postgres',
});

async function migrateUsers() {
  await client.connect();

  const data = fs.readFileSync('/tmp/mongodb_users.json', 'utf8');
  const lines = data.trim().split('\n');

  let imported = 0;
  let skipped = 0;

  for (const line of lines) {
    try {
      const user = JSON.parse(line);

      // Check if user already exists
      const existingUser = await client.query(
        'SELECT id FROM users WHERE id = $1',
        [user._id]
      );

      if (existingUser.rows.length > 0) {
        console.log(`Skipping duplicate user: ${user.email}`);
        skipped++;
        continue;
      }

      // Insert user into PostgreSQL
      await client.query(
        `INSERT INTO users (id, company_id, name, email, password, role, is_active, profile_image_url, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
         ON CONFLICT (company_id, email) DO NOTHING`,
        [
          user._id,
          user.company_id || '00000000-0000-0000-0000-000000000000', // Default company if not set
          user.name,
          user.email,
          user.password || '', // Might be empty if using OAuth
          user.role || 'USER',
          user.is_active !== undefined ? user.is_active : true,
          user.profile_image_url || null,
          user.created_at || new Date(),
          user.updated_at || new Date()
        ]
      );

      console.log(`Imported user: ${user.name} (${user.email})`);
      imported++;

    } catch (err) {
      console.error(`Error processing user:`, err.message);
    }
  }

  console.log(`\nMigration complete!`);
  console.log(`Imported: ${imported} users`);
  console.log(`Skipped: ${skipped} users`);

  await client.end();
}

migrateUsers().catch(console.error);
EOF

# Check if Node.js is available
if command -v node &> /dev/null; then
  echo "Running import script..."
  PG_HOST="$PG_HOST" PG_PORT="$PG_PORT" PG_DB="$PG_DB" PG_USER="$PG_USER" PG_PASSWORD="$PG_PASSWORD" node /tmp/import_users.js
else
  echo "Node.js not found. Using Python fallback..."

  # Python fallback
  cat > /tmp/import_users.py << 'EOF'
import json
import psycopg2
import os
from datetime import datetime

# Connect to PostgreSQL
conn = psycopg2.connect(
    host=os.getenv('PG_HOST', 'localhost'),
    port=int(os.getenv('PG_PORT', 5432)),
    database=os.getenv('PG_DB', 'fiber_app'),
    user=os.getenv('PG_USER', 'postgres'),
    password=os.getenv('PG_PASSWORD', 'postgres')
)
cur = conn.cursor()

# Read MongoDB export
with open('/tmp/mongodb_users.json', 'r') as f:
    lines = f.readlines()

imported = 0
skipped = 0

for line in lines:
    try:
        user = json.loads(line.strip())

        # Check if user already exists
        cur.execute('SELECT id FROM users WHERE id = %s', (user['_id'],))
        if cur.fetchone():
            print(f"Skipping duplicate user: {user['email']}")
            skipped += 1
            continue

        # Insert user
        cur.execute("""
            INSERT INTO users (id, company_id, name, email, password, role, is_active, profile_image_url, created_at, updated_at)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            ON CONFLICT (company_id, email) DO NOTHING
        """, (
            user['_id'],
            user.get('company_id', '00000000-0000-0000-0000-000000000000'),
            user.get('name', ''),
            user.get('email', ''),
            user.get('password', ''),
            user.get('role', 'USER'),
            user.get('is_active', True),
            user.get('profile_image_url'),
            user.get('created_at', datetime.now()),
            user.get('updated_at', datetime.now())
        ))

        print(f"Imported user: {user.get('name')} ({user.get('email')})")
        imported += 1

    except Exception as e:
        print(f"Error: {e}")

conn.commit()
cur.close()
conn.close()

print(f"\nMigration complete!")
print(f"Imported: {imported} users")
print(f"Skipped: {skipped} users")
EOF

  python3 /tmp/import_users.py
fi

echo "Migration completed!"
