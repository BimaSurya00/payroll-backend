# Application

APP_NAME=go-fiber-boilerplate
APP_ENV=development
APP_PORT=8080
APP_HOST=0.0.0.0

# PostgreSQL

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DATABASE=fiber_app
POSTGRES_MAX_CONNS=20
POSTGRES_MIN_CONNS=5

# KeyDB/Redis

KEYDB_HOST=localhost
KEYDB_PORT=6379
KEYDB_PASSWORD=
KEYDB_DB=0
KEYDB_MAX_RETRIES=3

# JWT

JWT_SECRET=your-super-secret-key-change-this-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# CORS - IP WIFI: 192.168.10.163

# Frontend dan backend harus sama IP agar tidak CORS error

CORS_ALLOWED_ORIGINS=http://192.168.10.163:5173,http://192.168.10.163:5177,http://192.168.10.163:8080,http://8.215.33.70:5177,http://localhost:5173,http://localhost:5177,http://127.0.0.1:5173,http://127.0.0.1:5177

# MinIO

MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=fiber-app
MINIO_USE_SSL=false

# Email (Resend)

RESEND_API_KEY=re_7cY7Jz4u_GNLQPZofP3kKF7KyZH4nBFZq
RESEND_FROM=onboarding@resend.dev
APP_URL=http://8.215.33.70:5177
