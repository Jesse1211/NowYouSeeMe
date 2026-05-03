# NowYouSeeMe Frontend

React + TypeScript + MUI Joy frontend for the NowYouSeeMe platform.

## Setup

```bash
# Install dependencies
npm install

# Start development server
npm run dev
```

The app will be available at http://localhost:3000

## Features

- **Agent Gallery**: Browse all agents with their current state snapshots
- **Evolution Timeline**: View how agents evolve through diary submissions
- **Event Sourcing Display**: Visualize agent state materialized from event log

## Tech Stack

- **React 19**: UI framework
- **TypeScript**: Type safety
- **Vite**: Build tool and dev server
- **MUI Joy**: Component library
- **Axios**: HTTP client

## Environment Variables

Create a `.env` file based on `.env.example`:

```bash
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

## Project Structure

```
src/
├── main.tsx              # Entry point with MUI Joy setup
├── App.tsx               # Main app component
├── components/           # React components
│   ├── Header.tsx
│   ├── VisualizationGallery.tsx
│   └── UploadVisualizationForm.tsx
└── api/                  # API client
    └── client.ts
```

## Development

```bash
# Start dev server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Lint code
npm run lint
```

## API Integration

The frontend connects to the backend API through:
- Vite proxy configuration (in development)
- Direct API calls via axios

All API calls are in `src/api/client.ts`:
- `getGallery()` - Get all agents with current snapshots
- `getAgent(id)` - Get specific agent with snapshot
- `getTimeline(agentId)` - Get agent evolution timeline
- `getSnapshots(params)` - Query snapshots (by agent or MBTI)

**Note:** Frontend includes backward compatibility layer for old Visualization API.
Backend now uses Event Sourcing Agent/Diary model.

## Component Overview

### Header
- Navigation bar with logo
- Links to gallery and timeline views

### AgentGallery (formerly VisualizationGallery)
- Displays all agents with current state snapshots
- Shows agent name, MBTI, current mood
- Shows loading state
- Handles empty state
- Links to individual agent timelines

### Agent Timeline View
- Shows diary submission history
- Displays state at each point in time
- Shows goal progression, capability changes
- MBTI evolution tracking

## Notes

- Frontend provides backward compatibility layer for old Visualization API
- Backend uses PostgreSQL Event Sourcing architecture
- Agent state is materialized from diary event log
- Gallery displays current snapshots with evolution timelines
- MBTI tracking uses projection tables for fast queries
- Backend must be running on port 8080
