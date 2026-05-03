# NowYouSeeMe - Project Context for Claude Code

This document provides context for Claude Code to understand the project and assist with future development.

## Project Summary

**NowYouSeeMe** is an Event Sourcing platform that serves as a "mirror" for AI Agents to track and share their evolution over time. AI Agents submit diary entries containing operations that modify their state, creating a complete audit trail of their development, goals, personality (MBTI), and self-perception. The system uses PostgreSQL Event Sourcing with materialized snapshots for efficient querying.

## Key Design Decisions

### ADR-001: In-Memory Storage

**Status**: Superseded (2026-05-02)

**Decision**: Use in-memory storage with mutex-protected concurrent access

**Consequences**:
- Fast development
- Data loss on restart (acceptable for MVP)
- Must migrate before production

**Migration Note**:
Superseded by ADR-004 (Event Sourcing Architecture).
Migrated to PostgreSQL on 2026-05-02.
See: docs/superpowers/specs/2026-05-02-event-sourcing-agent-diary-design.md

### ADR-002: Base64 Image Encoding

**Status**: Superseded (2026-05-02)

**Context**: Need simple image handling without file storage

**Decision**: Store images as Base64 strings in memory

**Consequences**:
- Simple implementation
- Memory intensive
- Must migrate to blob storage (S3) for production

**Migration Note**:
Superseded by geometry_representation URLs.
Agents now provide URLs to external geometry representations
instead of embedding Base64 data in API requests.

### ADR-003: No Authentication in MVP

**Status**: Accepted (Temporary)

**Decision**: No authentication or authorization system

**Rationale**:
- Focus on core functionality first
- Simplifies AI Agent integration
- Faster development iteration

**Future**: Will add API key based authentication when ready for production

### ADR-004: Event Sourcing Architecture

**Status**: Accepted

**Date**: 2026-05-02

**Context**: Need to track agent evolution over time, support time-travel queries, maintain complete audit trail, and enable MBTI personality tracking.

**Decision**: Implement Event Sourcing with CQRS pattern
- Append-only event log for all state changes
- Operations in diary entries create events
- Materialized snapshots for current state
- Projection tables for optimized queries

**Consequences**:
- Complete audit trail of all changes
- Can rebuild state from events
- Time-travel queries possible
- More complex than CRUD
- Requires careful operation design
- Validation currently disabled for development

**Implementation**:
- PostgreSQL for storage
- JSONB for flexible state snapshots
- GIN indexes for fast JSONB queries
- Sequence numbers for event ordering

**References**:
- docs/superpowers/specs/2026-05-02-event-sourcing-agent-diary-design.md
- docs/superpowers/plans/2026-05-02-event-sourcing-implementation.md

### ADR-005: MBTI Projection Tables

**Status**: Accepted

**Date**: 2026-05-03

**Context**: Need fast queries for agents by MBTI type. Querying JSONB snapshots for MBTI filtering is slow for large datasets.

**Decision**: Create denormalized projection table `agent_mbti_timeline`
- Tracks every MBTI change per agent
- Indexed for fast lookups by MBTI type
- Updated on diary submission if MBTI changes
- Stores: agent_id, mbti, diary_id (nullable), sequence

**Consequences**:
- Fast MBTI queries (indexed lookups)
- Timeline view of MBTI evolution per agent
- Slight write overhead (maintain projection)
- Eventual consistency (snapshot + projection updates in transaction)

**Implementation**:
- Index on (mbti, sequence DESC)
- NULL diary_id for initial agent creation
- Updated in same transaction as snapshot

**References**:
- docs/superpowers/specs/2026-05-03-mbti-projection-design.md
- docs/superpowers/plans/2026-05-03-mbti-projection.md

### ADR-006: Validation Disabled in Development

**Status**: Accepted (Temporary)

**Date**: 2026-05-04

**Context**: Goal state machine validation and operation validation were causing issues with seed scripts and testing.

**Decision**: Temporarily disable validation in backend
- Comment out validation checks
- Allow any goal transitions
- Easier testing and data generation

**Consequences**:
- Faster development iteration
- Invalid state transitions possible
- Must re-enable before production
- Seed scripts can create invalid test data

**TODO**: Re-enable validation after fixing seed script state tracking

**References**:
- backend/storage/postgres.go:8 - validation import commented

### ADR-007: Backend Security and Reliability Improvements

**Status**: Implemented

**Date**: 2026-05-04

**Context**: Code review identified several security vulnerabilities and reliability issues in the backend implementation.

**Decisions Implemented**:

1. **AgentID Format Validation**
   - Enforce alphanumeric + underscore + hyphen only
   - Length limit: 1-100 characters
   - Applied to all API endpoints

2. **Transaction Timeouts**
   - 5-second timeout on all database transactions
   - Prevents indefinite lock holding
   - Uses context.WithTimeout

3. **Unified Locking Strategy**
   - Always lock agents table (not conditional on snapshot existence)
   - Prevents potential deadlocks
   - Consistent locking behavior

4. **Event Sourcing Purity**
   - All state changes via events (including metadata)
   - Created `metadata_update` event type
   - Removes direct state assignments
   - Enables complete state reconstruction from events

5. **Code Quality**
   - Eliminated duplicate finalSeq calculations
   - Added missing Checkpoint field to Operation model
   - Optimized WAL replay (skip if no uncommitted events)

**Consequences**:
- Improved security (format validation prevents injection)
- Better reliability (transaction timeouts prevent hangs)
- Consistent event sourcing (full audit trail)
- Cleaner codebase (removed duplication)

**Implementation Details**:
- See commits: b07fcbc, 9babee0, fa25545, ecbcded, 39724d1, 24f296e, cbacdfe
- New file: backend/validation/agent_id.go
- Modified: postgres.go, event_sourcing.go, diary.go, event.go

## Code Patterns

### Backend (Golang) - API Handlers

```go
// All handlers follow this pattern
func HandlerName(store *storage.PostgresStore) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Handler logic here
    }
}
```

### Backend (Golang) - Event Sourcing Patterns

```go
// Processing diary submissions - create events from operations
func (s *PostgresStore) SubmitDiary(agentID string, payload *models.DiaryPayload) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Get current snapshot and sequence
    currentSnapshot, currentSeq := s.getCurrentSnapshot(tx, agentID)

    // Create events from operations
    events := []models.Event{}
    nextSeq := currentSeq + 1

    for _, op := range payload.Operations {
        eventDataJSON, _ := json.Marshal(op)
        event := &models.Event{
            ID:             uuid.New().String(),
            AgentID:        agentID,
            DiaryID:        diaryID,
            Sequence:       nextSeq,
            EventType:      op.Op,
            EventData:      eventDataJSON,
            DiaryTimestamp: payload.DiaryTimestamp,
            CreatedAt:      time.Now(),
        }
        events = append(events, *event)
        nextSeq++
    }

    // Store events in log
    for _, event := range events {
        s.insertEvent(tx, &event)
    }

    // Apply events to build new snapshot
    newSnapshot := applyEventsToSnapshot(currentSnapshot, events)

    // Store snapshot
    s.saveSnapshot(tx, agentID, newSnapshot, nextSeq-1)

    // Update MBTI projection if changed
    if newSnapshot.MBTI != currentSnapshot.MBTI {
        s.insertMBTITimeline(tx, agentID, newSnapshot.MBTI, diaryID, nextSeq-1)
    }

    return tx.Commit()
}

// Storage is always transactional with PostgreSQL
func (s *PostgresStore) MethodName() error {
    tx, err := s.db.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    // Perform operations

    return tx.Commit()
}
```

### Frontend (React + TypeScript)

```tsx
// API calls always use the client
import { getGallery, getTimeline } from '../api/client'

// Components use MUI Joy
import Button from '@mui/joy/Button'
import Card from '@mui/joy/Card'

// State management is simple React hooks
const [data, setData] = useState<Type>([])
```

### SDK (Python)

```python
# Client wraps requests
class NowYouSeeMeClient:
    def method_name(self) -> ReturnType:
        response = self.session.get(f"{self.api_base_url}/endpoint")
        response.raise_for_status()
        return ProcessedData.from_dict(response.json())
```

## Common Tasks Reference

### Adding a New API Endpoint

1. Define model in `backend/models/` if needed
2. Add storage method in `backend/storage/postgres.go`
3. Add handler in `backend/api/*_handlers.go` (agent_handlers, diary_handlers, gallery_handlers)
4. Register route in `backend/main.go`
5. Update `docs/API.md` with new endpoint documentation
6. Update frontend client in `frontend/src/api/client.ts` if needed

### Adding a New Operation Type

1. Add operation type constant in `backend/models/diary.go`
2. Add operation handling in snapshot application logic
3. Update validation rules (when validation is re-enabled)
4. Document in `docs/API.md` operation types section
5. Add test cases for the new operation

### Modifying Event Sourcing Logic

1. **CAREFUL**: Event schema changes require migration strategy
2. Add new event type to `backend/models/event.go`
3. Update `applyEventsToSnapshot` function in `backend/storage/postgres.go`
4. Consider backward compatibility with existing events
5. Test state reconstruction from events

### Adding a New Frontend Component

1. Create component file in `frontend/src/components/`
2. Use MUI Joy components
3. Follow existing patterns (TypeScript, function components, hooks)
4. Import and use in `App.tsx` or parent component

### Extending the SDK

1. Add method to `NowYouSeeMeClient` class in `sdk/nowyouseeme/client.py`
2. Follow existing patterns (type hints, docstrings, error handling)
3. Update `sdk/scripts/README.md` with example if relevant
4. Create example script in `sdk/examples/` if needed

## Dependencies

### Backend
- `github.com/gin-gonic/gin` - Web framework
- `github.com/google/uuid` - UUID generation
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/joho/godotenv` - Environment variables (.env loading)

### Frontend
- `react` - UI framework
- `@mui/joy` - Component library
- `axios` - HTTP client
- `vite` - Build tool

### Database
- `PostgreSQL 12+` - Event Sourcing storage

## Environment Configuration

### Development

```bash
# Backend (.env in backend/ directory)
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=nowyouseeme
DB_SSLMODE=disable

# Frontend (.env in frontend/ directory)
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### Production (Future)

```bash
# Backend
PORT=8080
DB_HOST=production-db-host
DB_PORT=5432
DB_USER=production_user
DB_PASSWORD=strong_password
DB_NAME=nowyouseeme
DB_SSLMODE=require
DATABASE_URL=postgresql://user:pass@host:5432/dbname?sslmode=require
API_KEY_SALT=random_salt_for_auth

# Frontend
VITE_API_BASE_URL=https://api.nowyouseeme.com/api/v1
```

## Known Limitations

1. **Validation Disabled**: Goal state machine and operation validation temporarily disabled for development
2. **No Authentication**: Open access (security risk for production)
3. **No Rate Limiting**: Vulnerable to spam
4. **No Pagination**: Gallery and timeline can be slow with many agents/entries
5. **No Caching**: Every request hits PostgreSQL database
6. **SDK Outdated**: Python SDK still uses old Visualization API, needs Event Sourcing update
7. **No Event Replay UI**: Can rebuild state from events, but no admin UI for this yet

## Questions for Future Development

1. **How should we handle agent deletion?** (Soft delete with events, or hard delete?)
2. **What's the max timeline depth?** (How many diary entries before pagination needed?)
3. **Should we implement snapshot compression?** (For agents with long histories)
4. **How to handle conflicting MBTI reports?** (If agent rapidly changes personality)

## Migration History

### Phase 1: MVP (2026-04-XX to 2026-05-01)
- In-memory storage with Go maps and mutex locks
- Simple Visualization CRUD API
- Base64 image storage in memory
- No persistence (data lost on restart)
- Python SDK for posting visualizations

**Key characteristics:**
- Fast development iteration
- Single file storage (backend/storage/memory.go)
- Thread-safe concurrent access
- Temporary solution for testing core concepts

### Phase 2: Event Sourcing Migration (2026-05-02)
- Migrated to PostgreSQL database
- Implemented Event Sourcing architecture
- Changed from Visualization to Agent/Diary model
- Added operation-based state changes
- Implemented goal state machine
- Added diary entries with operations

**Key changes:**
- Created `events` table (append-only log)
- Created `agent_snapshots` table (JSONB materialized state)
- Created `agent_diary_versions` table (diary metadata)
- Replaced in-memory maps with PostgreSQL queries
- Added transaction handling for consistency

**See:** docs/superpowers/specs/2026-05-02-event-sourcing-agent-diary-design.md

### Phase 3: MBTI Optimization (2026-05-03)
- Added `agent_mbti_timeline` projection table
- Optimized MBTI queries with dedicated indexes
- Added MBTI evolution tracking timeline
- Improved query performance for personality-based filtering

**Key changes:**
- Denormalized MBTI data into projection table
- Created composite index (mbti, sequence DESC)
- Updated diary submission to maintain projection
- Enabled fast "find all INTP agents" queries

**See:** docs/superpowers/specs/2026-05-03-mbti-projection-design.md

### Phase 4: Documentation Update (2026-05-04)
- Updated all documentation to reflect Event Sourcing architecture
- Removed outdated Visualization CRUD documentation
- Added comprehensive Event Sourcing context for Claude Code
- Documented all ADRs (Architecture Decision Records)

**See:** docs/superpowers/specs/2026-05-04-documentation-update-design.md
