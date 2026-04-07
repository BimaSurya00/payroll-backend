#!/bin/bash

# ============================================
# Quick Migration: MongoDB Users → PostgreSQL
# ============================================

echo "Migrating users from MongoDB to PostgreSQL..."

# Create temporary file for user data
TEMP_FILE=$(mktemp)

# Export users from MongoDB
echo "Exporting from MongoDB..."
mongosh --quiet localhost:27017/fiber_app --eval "
  db.users.find().forEach(function(user) {
    var company_id = user.company_id || '00000000-0000-0000-0000-000000000000';
    var name = user.name || '';
    var email = user.email || '';
    var password = user.password || '';
    var role = user.role || 'USER';
    var is_active = (user.is_active !== undefined) ? user.is_active : true;
    var profile_image_url = user.profile_image_url || null;
    var created_at = user.created_at || new Date();
    var updated_at = user.updated_at || new Date();

    print(user._id + '|' + company_id + '|' + name.replace(/\|/g, '') + '|' + email + '|' + password + '|' + role + '|' + is_active + '|' + (profile_image_url || '') + '|' + created_at.toISOString() + '|' + updated_at.toISOString());
  });
" > "$TEMP_FILE"

# Check if we got users
USER_COUNT=$(wc -l < "$TEMP_FILE")
echo "Found $USER_COUNT users in MongoDB"

if [ "$USER_COUNT" -eq 0 ]; then
  echo "No users to migrate!"
  rm "$TEMP_FILE"
  exit 1
fi

# Import into PostgreSQL
echo "Importing into PostgreSQL..."
{
  echo "BEGIN;"
  while IFS='|' read -r _id company_id name email password role is_active profile_image_url created_at updated_at; do
    # Skip empty lines
    [ -z "$_id" ] && continue

    # Generate SQL
    echo "INSERT INTO users (id, company_id, name, email, password, role, is_active, profile_image_url, created_at, updated_at)"
    echo "VALUES ('$_id', '$company_id', '\$(echo "$name" | sed \"s/'/''/g\")', '$email', '\$(echo "$password" | sed \"s/'/''/g\")', '$role', $is_active, $(if [ -n "$profile_image_url" ]; then echo "'$profile_image_url'"; else echo "NULL"; fi), '$created_at', '$updated_at')"
    echo "ON CONFLICT (company_id, email) DO NOTHING;"
  done < "$TEMP_FILE"
  echo "COMMIT;"
} | psql -h localhost -U postgres -d fiber_app -q

# Clean up
rm "$TEMP_FILE"

echo "Migration completed!"

# Verify
echo "Verifying..."
PGPASSWORD=postgres psql -h localhost -U postgres -d fiber_app -c "SELECT COUNT(*) as total_users FROM users;"
