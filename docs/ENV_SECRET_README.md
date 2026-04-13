# GitHub Secret ENV_FILE - Production

# ============================================================
# CARA UPDATE:
# 1. GitHub repo → Settings → Secrets and variables → Actions
# 2. Klik ENV_FILE → Delete
# 3. New repository secret → Name: ENV_FILE
# 4. Paste isi ENV_FILE_STAGING.txt (untuk staging) atau
#    ENV_FILE_PRODUCTION.txt (untuk production)
# 5. Add secret
# 6. Re-run workflow di Actions tab
# ============================================================

# PENTING:
# - CORS_ALLOWED_ORIGINS JANGAN gunakan * (wildcard)
# - POSTGRES_HOST harus "postgres" (nama container Docker)
# - KEYDB_HOST harus "keydb" (nama container Docker)
# - MINIO_ENDPOINT harus "minio:9000" (nama container Docker)
# - Nilai di bawah adalah DEFAULT, sesuaikan dengan environment masing-masing
