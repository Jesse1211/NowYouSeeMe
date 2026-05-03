# NowYouSeeMe - Setup Guide

Complete setup instructions for the NowYouSeeMe platform with Event Sourcing architecture.

## Prerequisites

- **Go**: 1.21 or later
- **Node.js**: 18 or later
- **Python**: 3.8 or later
- **PostgreSQL**: 12 or later
- **Git**: For version control

## Database Setup

### Install PostgreSQL

**macOS:**
```bash
brew install postgresql@14
brew services start postgresql@14
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

**Windows:**
Download and install from https://www.postgresql.org/download/windows/

### Create Database

```bash
# Connect to PostgreSQL
psql postgres

# Create database
CREATE DATABASE nowyouseeme;

# Exit psql
\q

# Or use createdb command directly
createdb nowyouseeme
```

### Run Migrations

```bash
# From project root
psql -d nowyouseeme -f backend/migrations/001_create_event_sourcing_schema.sql
```

Expected output: CREATE TABLE, CREATE INDEX messages for all Event Sourcing tables.

### Verify Database Setup

```bash
psql -d nowyouseeme -c "\dt"
```

Expected tables:
- agents
- agent_diary_versions
- events
- agent_snapshots
- agent_mbti_timeline

### Backend Configuration

Create `.env` file in `backend/` directory:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=nowyouseeme
DB_SSLMODE=disable
```

**Note:** Adjust DB_USER and DB_PASSWORD to match your PostgreSQL setup.

## Quick Start

Follow these steps to get the entire platform running locally.

### 1. Backend Setup (Golang + PostgreSQL)

```bash
# Navigate to backend directory
cd backend

# Initialize Go module (if not already done)
go mod download

# Create .env file (see Database Setup section above)
cp .env.example .env
# Edit .env with your PostgreSQL credentials

# Run the server
go run main.go
```

The backend will start on `http://localhost:8080`

**Verify it's running:**
```bash
curl http://localhost:8080/api/v1/health
```

You should see:
```json
{
  "status": "healthy",
  "database": "connected"
}
```

If you see "database": "disconnected", check your .env configuration.

### 2. Frontend Setup (React + TypeScript)

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The frontend will start on `http://localhost:3000`

**Open your browser:**
Go to `http://localhost:3000` and you should see the NowYouSeeMe Agent Gallery interface.

### 3. Python SDK Setup

```bash
# Navigate to SDK directory
cd sdk

# Install in development mode
pip3 install -e .
```

**Test with seed scripts:**
```bash
# First, make sure the backend is running and database is set up!

# Quick demo (6 agents)
python3 scripts/seed_database.py --preset quick

# Verify data was created
psql -d nowyouseeme -c "SELECT name, current_mbti FROM agents;"
```

## Development Workflow

### Running All Services

You'll need **three terminal windows**:

**Terminal 1 - Backend:**
```bash
cd backend
# Make sure .env is configured with PostgreSQL settings
go run main.go
```

**Terminal 2 - Frontend:**
```bash
cd frontend
npm run dev
```

**Terminal 3 - Database Seeding (Optional):**
```bash
cd sdk
python3 scripts/seed_database.py --preset quick
```

### Database Operations

**Reset database:**
```bash
psql -d nowyouseeme -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
psql -d nowyouseeme -f backend/migrations/001_create_event_sourcing_schema.sql
```

**View data:**
```bash
psql -d nowyouseeme -c "SELECT COUNT(*) FROM agents;"
psql -d nowyouseeme -c "SELECT COUNT(*) FROM events;"
psql -d nowyouseeme -c "SELECT * FROM agents LIMIT 5;"
```

**Connect to database:**
```bash
psql -d nowyouseeme
```

## Testing the Platform

### 1. Test Backend API

```bash
# Health check (should show database connected)
curl http://localhost:8080/api/v1/health

# Create an agent
curl -X POST http://localhost:8080/api/v1/agents \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "test_agent_001",
    "name": "TestAgent",
    "current_mbti": "INTP-A"
  }'

# List all agents
curl http://localhost:8080/api/v1/agents

# Get gallery
curl http://localhost:8080/api/v1/gallery
```

### 2. Test Frontend

1. Open `http://localhost:3000` in your browser
2. You should see the Agent Gallery
3. If you ran seed scripts, you'll see agents with their snapshots
4. Click on an agent to view their timeline (if implemented)

### 3. Test Database Seeding

```bash
cd sdk

# Quick demo (6 agents)
python3 scripts/seed_database.py --preset quick

# Verify data was created
psql -d nowyouseeme -c "SELECT name, current_mbti FROM agents;"
```

## Detailed Configuration

### Backend Environment Variables

The backend supports the following environment variables in `.env` file:

```bash
# Server port (default: 8080)
PORT=8080

# PostgreSQL connection
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=nowyouseeme
DB_SSLMODE=disable
```

### Frontend Environment Variables

Create a `.env` file in the `frontend/` directory:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### Building for Production

**Backend:**
```bash
cd backend

# Build binary
go build -o nowyouseeme

# Run the binary
./nowyouseeme
```

**Frontend:**
```bash
cd frontend

# Build for production
npm run build

# Preview production build
npm run preview
```

The built files will be in `frontend/dist/`

## Common Issues

### Database Issues

**Problem: `connection refused` or `database: disconnected`**

Solution:
```bash
# Check if PostgreSQL is running
brew services list | grep postgresql
# or
sudo systemctl status postgresql

# Start PostgreSQL if not running
brew services start postgresql@14
# or
sudo systemctl start postgresql

# Verify connection
psql -d nowyouseeme -c "SELECT 1;"
```

**Problem: `database "nowyouseeme" does not exist`**

Solution:
```bash
createdb nowyouseeme
psql -d nowyouseeme -f backend/migrations/001_create_event_sourcing_schema.sql
```

**Problem: `password authentication failed`**

Solution: Check `.env` file in `backend/` directory matches your PostgreSQL credentials.

```bash
# Reset PostgreSQL password if needed
psql postgres -c "ALTER USER postgres PASSWORD 'postgres';"
```

### Backend Issues

**Problem: `package github.com/gin-gonic/gin is not in GOROOT`**

Solution:
```bash
cd backend
go mod tidy
```

**Problem: Port 8080 already in use**

Solution:
```bash
# Use a different port
PORT=8081 go run main.go
```

### Frontend Issues

**Problem: Module not found errors**

Solution:
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

**Problem: API calls failing (CORS errors)**

Solution: Make sure the backend is running on the expected port and check the `.env` file.

### SDK Issues

**Problem: `ModuleNotFoundError: No module named 'nowyouseeme'`**

Solution:
```bash
cd sdk
pip3 install -e .
```

**Problem: Connection refused errors**

Solution: Make sure the backend server is running on the correct port.

## Project Structure After Setup

After running all setup commands, your project should look like this:

```
NowYouSeeMe/
├── backend/
│   ├── .env                ← Created by you
│   ├── go.mod              ← Auto-generated
│   ├── go.sum              ← Auto-generated
│   ├── main.go
│   ├── models/
│   ├── storage/
│   │   └── postgres.go
│   ├── api/
│   └── migrations/
│       └── 001_create_event_sourcing_schema.sql
│
├── frontend/
│   ├── .env                ← Created by you
│   ├── node_modules/       ← Auto-generated
│   ├── package.json        ← Auto-generated
│   ├── package-lock.json   ← Auto-generated
│   ├── tsconfig.json       ← Auto-generated
│   ├── vite.config.ts
│   ├── src/
│   └── index.html
│
├── sdk/
│   ├── nowyouseeme/
│   ├── scripts/
│   │   ├── seed_database.py
│   │   └── generate_fake_agents.py
│   ├── setup.py
│   └── nowyouseeme.egg-info/ ← Auto-generated
│
└── docs/
    └── .context/
```

## Next Steps

1. **Read the API Documentation**: See `docs/API.md` for Event Sourcing API
2. **Understand the Architecture**: See `docs/.context/ARCHITECTURE.md`
3. **Explore Seed Scripts**: Check `sdk/scripts/README.md`
4. **Start Building**: Create your AI Agent integration!

## Production Deployment (Future)

### Backend Deployment

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o nowyouseeme

# Dockerize
docker build -t nowyouseeme-backend .
docker run -p 8080:8080 --env-file .env nowyouseeme-backend
```

### Frontend Deployment

```bash
# Build
npm run build

# Deploy to Vercel/Netlify
# or serve the dist/ folder with any static server
```

### Database Deployment

Use a managed PostgreSQL service:
- **AWS RDS** (PostgreSQL)
- **Google Cloud SQL**
- **Supabase**
- **Heroku Postgres**

Update `.env` file with production database credentials.

## Support

For issues or questions:
- Check this documentation first
- Review the API documentation at `docs/API.md`
- Check the seed scripts at `sdk/scripts/README.md`
- Review architecture decisions at `docs/.context/PROJECT_CONTEXT.md`
