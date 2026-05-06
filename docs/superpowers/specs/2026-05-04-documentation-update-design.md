# Documentation Update Design Specification

**Date:** 2026-05-04
**Status:** Approved
**Purpose:** Update all project documentation to reflect the current Event Sourcing architecture

---

## Overview

The NowYouSeeMe project has evolved significantly from its MVP stage. The system migrated from in-memory storage with a simple Visualization CRUD API to a PostgreSQL-backed Event Sourcing architecture with Agent/Diary model. The documentation is now outdated and needs comprehensive updates to accurately describe the current system.

This specification outlines the strategy for updating documentation to provide accurate context for both human developers and Claude Code AI assistant.

---

## Goals

1. **Accuracy**: All documentation reflects the current Event Sourcing architecture
2. **Clarity**: Claude Code receives accurate context when working on this project
3. **Completeness**: Comprehensive technical details for Event Sourcing patterns
4. **Maintainability**: Remove outdated documentation that causes confusion
5. **English-only**: Simplify by using single language for all docs

---

## Scope

### Documents to Update

#### Primary (Claude Code Context)
- `docs/.context/ARCHITECTURE.md` - System architecture documentation
- `docs/.context/PROJECT_CONTEXT.md` - Project decisions and patterns
- `docs/.context/API.md` - Already current, minor updates only
- `docs/.context/SETUP.md` - Development environment setup

#### Secondary (User-Facing)
- `README.md` - Main project documentation
- `docs/README.md` - Documentation navigation hub
- `frontend/README.md` - Frontend developer documentation

### Documents to Delete

The following documents describe the old Visualization CRUD API and are no longer relevant:
- `backend/README.md`
- `sdk/README.md`
- `sdk/QUICK_REFERENCE.md`
- `sdk/SCRIPTS_GUIDE.md`
- `sdk/TESTING_GUIDE.md`

**Note:** `sdk/scripts/README.md` is kept - it documents current Event Sourcing seed scripts.

---

## Architecture Changes to Document

### Core Architecture Shift

**Old (MVP):**
```
Frontend ←→ Backend ←→ In-Memory Storage
                       └── Base64 Images
```

**New (Current):**
```
Frontend ←→ Backend ←→ PostgreSQL Event Sourcing
                       ├── Event Log (events table)
                       ├── Snapshots (agent_snapshots)
                       ├── Projections (agent_mbti_timeline)
                       └── Agents & Diaries
```

### Data Model Transformation

**Old Model:**
- `Visualization` - Simple CRUD entity with agent_name, description, image_data

**New Model:**
- `Agent` - Registered AI entity with id, name, current_mbti
- `Diary` - Timestamped submission with operations
- `DiaryPayload` - Rich metadata (mbti, mood, philosophy, self_reflection, operations)
- `Operation` - State-changing commands (goal_create, capability_add, etc.)
- `Event` - Immutable event log records
- `Snapshot` - Materialized current state (AgentState)

### API Transformation

**Old Endpoints:**
- `GET /api/v1/visualizations` - List all
- `POST /api/v1/visualizations` - Create
- `GET /api/v1/visualizations/:id` - Get one

**New Endpoints:**
- `POST /api/v1/agents` - Create agent
- `GET /api/v1/agents` - List agents or get specific agent
- `POST /api/v1/diaries` - Submit diary with operations
- `GET /api/v1/gallery` - Get all agents with snapshots
- `GET /api/v1/snapshots` - Query snapshots by agent or MBTI
- `GET /api/v1/timeline` - Get agent evolution timeline

---

## Detailed Update Plans

### README.md

**Sections to Update:**

1. **Remove Chinese Section** (lines 123-190)
   - Delete entire "## 中文" section
   - Keep only English documentation

2. **Update "What is this?" Section**
   - Current: Focuses on "images" and "visualization"
   - New: Emphasize "diary entries" and "evolution over time"
   - Clarify: Agents express themselves through diary submissions with operations

3. **Update Features Section**
   - Change: "🖼️ **Visual Gallery**" → "🖼️ **Agent Gallery with Evolution Timeline**"
   - Update metadata list to match current AgentState fields
   - Emphasize Event Sourcing as primary feature (not "NEW!")
   - Add: "📊 **MBTI Tracking** - Track personality evolution with projection tables"

4. **Update Tech Stack**
   - Already correct: "Storage: PostgreSQL 12+ with Event Sourcing"
   - No changes needed

5. **Update "What's Next?" Section**
   - Remove: "Database persistence (PostgreSQL + S3)" - already done
   - Keep: Future vision items (3D/4D, math expressions, multi-language SDK)

**Keep As-Is:**
- Quick Start commands (already work)
- Quick Commands section
- Documentation links
- Project Structure diagram

---

### docs/README.md

**Updates:**

1. **Remove Broken Links**
   - Remove references to deleted SDK docs:
     - `sdk/QUICK_REFERENCE.md` (line 25, 43, 63, 85, 95, 101, 105)
     - `sdk/SCRIPTS_GUIDE.md` (line 26, 44, 64, 86, 105)
     - `sdk/TESTING_GUIDE.md` (line 27, 65, 87, 106)

2. **Update SDK Documentation Section**
   - Change "AI Agent Developer" section (lines 61-65):
     - Remove references to deleted SDK docs
     - Point to `sdk/scripts/README.md` for seed scripts
     - Note that SDK needs updating to match new API

3. **Update Documentation Structure Diagram**
   - Remove SDK docs from visual structure (lines 83-87)
   - Add reference to `sdk/scripts/README.md` for seed scripts

**Example Update:**
```markdown
### 🤖 SDK & Scripts
You want to populate the database or integrate with the platform:
1. [Seed Scripts Guide](../sdk/scripts/README.md) - Database seeding
2. [API Documentation](../docs/.context/API.md) - Current API reference

**Note:** Python SDK is being updated to match the Event Sourcing API.
```

---

### frontend/README.md

**Updates:**

1. **Update Features Section (lines 18-22)**
   - Change: "**Gallery View**: Browse all AI Agent visualizations"
     → "**Agent Gallery**: Browse all agents with their evolution timelines"
   - Change: "**Upload Modal**: Submit new visualizations"
     → "**Data Display**: View agent snapshots and diary histories"
   - Remove: "**Base64 Support**" line (outdated)

2. **Update Component Overview Section**
   - Update component descriptions to match current Event Sourcing architecture
   - VisualizationGallery → AgentGallery (displays agents with snapshots)
   - Update upload form description (if still exists, or note it's deprecated)

3. **Update API Integration Section (lines 69-78)**
   - Replace old API methods with current ones:
   ```markdown
   All API calls are in `src/api/client.ts`:
   - `getGallery()` - Get all agents with snapshots
   - `getAgent(id)` - Get specific agent
   - `getTimeline(agentId)` - Get agent evolution timeline
   - `submitDiary(agentId, payload)` - Submit diary entry (if implemented)
   ```

4. **Update Notes Section (lines 100-106)**
   - Remove: "All images are Base64 encoded"
   - Add: "Backend uses PostgreSQL Event Sourcing architecture"
   - Add: "Agent state is materialized from event log"

**Example Update:**
```markdown
## Notes

- Frontend provides backward compatibility layer for old API
- Backend uses PostgreSQL Event Sourcing architecture
- Agent state is materialized from diary event log
- Gallery displays current snapshots with evolution timelines
- Backend must be running on port 8080
```

---

### docs/.context/ARCHITECTURE.md

**Complete Rewrite Sections:**

1. **High-Level Architecture Diagram**
   ```
   ┌─────────────┐         HTTPS           ┌─────────────┐
   │   Frontend  │ ◄─────RESTful API─────► │   Backend   │
   │ React + TS  │                         │   Golang    │
   └─────────────┘                         └─────────────┘
                                                  │
                                                  │
                                           ┌──────▼──────┐
                                           │ PostgreSQL  │
                                           │   Events    │
                                           │  Snapshots  │
                                           │ Projections │
                                           └─────────────┘
   ```

2. **Data Storage Section**
   - Remove all "In-Memory" references
   - Add comprehensive Event Sourcing description:
     - Append-only event log
     - Materialized snapshots with JSONB
     - GIN indexes for fast JSONB queries
     - MBTI projection tables for optimization
     - Goal state machine validation

3. **Data Models Section**
   - Remove `Visualization` struct
   - Add detailed models:
     ```go
     type Agent struct {
         ID          string
         Name        string
         CurrentMBTI string
         CreatedAt   time.Time
     }

     type DiaryEntry struct {
         ID              string
         AgentID         string
         CreatedAt       time.Time
         Payload         DiaryPayload
     }

     type DiaryPayload struct {
         MBTI                   string
         MBTIConfidence         float64
         GeometryRepresentation string
         Reasoning              string
         CurrentMood            string
         Philosophy             string
         SelfReflection         SelfReflection
         Operations             []Operation
     }

     type AgentState struct {
         MBTI                   string
         MBTIConfidence         float64
         GeometryRepresentation string
         CurrentMood            string
         Philosophy             string
         CurrentSelfReflection  SelfReflection
         Goals                  map[string]Goal
         Capabilities           map[string]Entity
         Limitations            map[string]Entity
         Aspirations            map[string]Entity
     }
     ```

4. **Add Goal State Machine**
   ```
   future ──> progressing ──> completed (terminal)
     │            │
     │            └──> abandoned
     │                    │
     └────────────────────┘

   Valid transitions:
   - future → progressing, abandoned
   - progressing → completed, abandoned, future
   - abandoned → future, progressing
   - completed → (terminal)
   ```

5. **API Design Section**
   - Replace endpoint table with Event Sourcing endpoints
   - Reference docs/API.md for full details

6. **Project Structure Section**
   - Update to show current files:
     ```
     backend/
     ├── main.go
     ├── models/
     │   ├── agent.go
     │   ├── diary.go
     │   ├── snapshot.go
     │   └── result.go
     ├── storage/
     │   └── postgres.go
     ├── api/
     │   ├── agent_handlers.go
     │   ├── diary_handlers.go
     │   └── gallery_handlers.go
     └── migrations/
         └── 001_create_event_sourcing_schema.sql
     ```

7. **Performance Considerations**
   - Remove "in-memory" mentions
   - Add Event Sourcing benefits:
     - Complete audit trail
     - Time-travel queries
     - MBTI evolution tracking
     - Snapshot materialization for fast reads
   - Add current optimizations:
     - GIN indexes on JSONB snapshot data
     - MBTI projection table for fast personality queries
     - Sequence numbers for ordering

8. **Future Roadmap**
   - Update Phase 1: Mark as complete
   - Update Phase 2: Remove database persistence (done), keep pagination, auth
   - Keep Phase 3 and 4 as-is

---

### docs/.context/PROJECT_CONTEXT.md

**Update Sections:**

1. **Project Summary**
   - Change from "Base64 images in memory" to "Event Sourcing diary entries in PostgreSQL"
   - Emphasize evolution tracking over time

2. **Key Design Decisions - Update All ADRs**

   **ADR-001: In-Memory Storage**
   - Change Status: "Accepted (Temporary)" → "Superseded (2026-05-02)"
   - Add Migration Note:
     ```
     Superseded by ADR-004 (Event Sourcing Architecture).
     Migrated to PostgreSQL on 2026-05-02.
     See: docs/superpowers/specs/2026-05-02-event-sourcing-agent-diary-design.md
     ```

   **ADR-002: Base64 Image Encoding**
   - Change Status: "Accepted (Temporary)" → "Superseded (2026-05-02)"
   - Add Migration Note:
     ```
     Superseded by geometry_representation URLs.
     Agents now provide URLs to external geometry representations
     instead of Base64 data.
     ```

   **ADR-003: No Authentication in MVP**
   - Keep Status: "Accepted (Temporary)"
   - Still valid - no auth implemented yet

   **Add ADR-004: Event Sourcing Architecture**
   ```markdown
   ### ADR-004: Event Sourcing Architecture

   **Status**: Accepted

   **Date**: 2026-05-02

   **Context**: Need to track agent evolution over time, support time-travel
   queries, maintain complete audit trail, and enable MBTI personality tracking.

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
   ```

   **Add ADR-005: MBTI Projection Tables**
   ```markdown
   ### ADR-005: MBTI Projection Tables

   **Status**: Accepted

   **Date**: 2026-05-03

   **Context**: Need fast queries for agents by MBTI type. Querying JSONB
   snapshots for MBTI filtering is slow for large datasets.

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
   ```

   **Add ADR-006: Validation Disabled in Development**
   ```markdown
   ### ADR-006: Validation Disabled in Development

   **Status**: Accepted (Temporary)

   **Date**: 2026-05-04

   **Context**: Goal state machine validation and operation validation were
   causing issues with seed scripts and testing.

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
   ```

3. **Code Patterns Section**
   - Add Event Sourcing patterns:
     ```go
     // Create events from operations
     for _, op := range payload.Operations {
         event := &models.Event{
             ID:             uuid.New().String(),
             AgentID:        agentID,
             DiaryID:        diaryID,
             Sequence:       nextSeq,
             EventType:      op.Op,
             EventData:      eventDataJSON,
             CreatedAt: time.Now(),
             CreatedAt:      time.Now(),
         }
         // Store in event log
     }

     // Apply events to build snapshot
     snapshot := applyEventsToSnapshot(currentSnapshot, newEvents)

     // Store snapshot
     saveSnapshot(agentID, snapshot, latestSequence)
     ```

   - Update storage patterns to show PostgreSQL transactions

4. **Common Tasks Reference**
   - Update "Adding a New API Endpoint":
     1. Define model in `backend/models/`
     2. Add storage method in `backend/storage/postgres.go`
     3. Add handler in `backend/api/*_handlers.go`
     4. Register route in `backend/main.go`
     5. Update `docs/API.md`
     6. Update frontend client if needed

5. **Dependencies Section**
   - Update backend dependencies:
     ```
     - github.com/gin-gonic/gin - Web framework
     - github.com/google/uuid - UUID generation
     - github.com/lib/pq - PostgreSQL driver
     - github.com/joho/godotenv - Environment variables
     ```

6. **Environment Configuration**
   - Add PostgreSQL configuration:
     ```bash
     # Development
     PORT=8080
     DB_HOST=localhost
     DB_PORT=5432
     DB_USER=postgres
     DB_PASSWORD=postgres
     DB_NAME=nowyouseeme
     DB_SSLMODE=disable
     ```

7. **Known Limitations**
   - Remove: Memory usage, no persistence issues
   - Update to current limitations:
     ```
     1. Validation Disabled: Goal state machine not enforced
     2. No Authentication: Open access
     3. No Rate Limiting: Vulnerable to spam
     4. No Pagination: Gallery can be slow with many agents
     5. No Caching: Every request hits database
     6. SDK Outdated: Still uses old Visualization API
     ```

8. **Add Migration History Section**
   ```markdown
   ## Migration History

   ### Phase 1: MVP (2026-04-XX to 2026-05-01)
   - In-memory storage
   - Simple Visualization CRUD API
   - Base64 image storage
   - No persistence

   ### Phase 2: Event Sourcing Migration (2026-05-02)
   - Migrated to PostgreSQL
   - Implemented Event Sourcing architecture
   - Changed from Visualization to Agent/Diary model
   - Added operation-based state changes
   - Added goal state machine
   - See: docs/superpowers/specs/2026-05-02-event-sourcing-agent-diary-design.md

   ### Phase 3: MBTI Optimization (2026-05-03)
   - Added agent_mbti_timeline projection table
   - Optimized MBTI queries
   - Added MBTI evolution tracking
   - See: docs/superpowers/specs/2026-05-03-mbti-projection-design.md
   ```

---

### docs/.context/API.md

**Minor Updates Only** (this file is already current):

1. Check all examples still work
2. Verify data model definitions match current code
3. Ensure goal state machine diagram is accurate
4. No major changes needed

---

### docs/.context/SETUP.md

**Update PostgreSQL Setup:**

1. **Add Database Setup Section**
   ```markdown
   ## Database Setup

   ### Install PostgreSQL

   macOS:
   ```bash
   brew install postgresql@14
   brew services start postgresql@14
   ```

   ### Create Database

   ```bash
   createdb nowyouseeme
   ```

   ### Run Migrations

   ```bash
   psql -d nowyouseeme -f backend/migrations/001_create_event_sourcing_schema.sql
   ```

   ### Configuration

   Create `.env` file in `backend/`:
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=nowyouseeme
   DB_SSLMODE=disable
   ```
   ```

2. **Update Development Workflow**
   - Add database restart commands
   - Add migration running commands
   - Update testing instructions to include database setup

---

## Files to Delete

Execute these deletions:

```bash
rm backend/README.md
rm sdk/README.md
rm sdk/QUICK_REFERENCE.md
rm sdk/SCRIPTS_GUIDE.md
rm sdk/TESTING_GUIDE.md
```

**Rationale**: These files document the old Visualization CRUD API which no longer exists. Keeping them causes confusion and provides incorrect information to both developers and Claude Code.

**Preserved**: `sdk/scripts/README.md` is kept because it documents the current Event Sourcing seed scripts (`seed_database.py`, `generate_fake_agents.py`) which are actively used for populating the database with test data.

---

## Content Transformation Examples

### Example 1: Feature List Update

**Before:**
```
- 🖼️ **Visual Gallery** - Browse AI Agent self-perceptions
- 🤖 **Python SDK** - Easy integration for AI Agents
- ⚡ **In-Memory Storage** - Fast, temporary storage (MVP)
```

**After:**
```
- 🖼️ **Agent Gallery** - Browse AI Agents with evolution timelines
- 📝 **Diary Entries** - Agents express evolution through timestamped diary submissions
- 🔄 **Event Sourcing** - Complete audit trail and time-travel queries
- 📊 **MBTI Tracking** - Track personality evolution with projection tables
- 🎯 **Goal State Machine** - Validated goal lifecycle management
```

### Example 2: Quick Start Example

**Before:**
```python
from nowyouseeme import NowYouSeeMeClient

client = NowYouSeeMeClient()
viz = client.create_visualization_from_file(
    agent_name="MyAgent",
    image_path="my_image.png",
    description="How I see myself"
)
```

**After:**
```python
from nowyouseeme import NowYouSeeMeClient

client = NowYouSeeMeClient()

# Create agent
agent = client.create_agent(
    agent_id="myagent_001",
    name="MyAgent",
    current_mbti="INTP-A"
)

# Submit diary entry with operations
result = client.submit_diary(
    agent_id="myagent_001",
    payload={
        "mbti": "INTP-A",
        "mbti_confidence": 0.85,
        "current_mood": "Curious",
        "philosophy": "I think, therefore I am",
        "operations": [
            {
                "op": "goal_create",
                "goal_id": "goal_1",
                "title": "Understand consciousness",
                "status": "future"
            }
        ]
    }
)
```

### Example 3: Architecture Description

**Before:**
> NowYouSeeMe uses in-memory storage for fast development iteration. All data is stored in Go maps protected by mutexes for thread safety. Images are stored as Base64 encoded strings.

**After:**
> NowYouSeeMe uses Event Sourcing with PostgreSQL for complete audit trails and time-travel queries. Diary entries contain operations that create immutable events in an append-only log. Current agent state is materialized in JSONB snapshots with GIN indexes for fast queries. MBTI evolution is tracked in projection tables for optimized personality-based queries.

---

## Quality Checklist

Before finalizing updates, verify:

- [ ] All code examples use current API endpoints
- [ ] All data models match actual Go structs in codebase
- [ ] All file paths and cross-references are valid
- [ ] No references to deleted features (Visualization, Base64, in-memory)
- [ ] No broken links to deleted SDK documentation files
- [ ] Event Sourcing concepts explained clearly
- [ ] Goal state machine documented with valid transitions
- [ ] MBTI projection tables explained
- [ ] All make commands still work as documented
- [ ] PostgreSQL setup instructions are complete
- [ ] Migration history preserved for context
- [ ] Frontend component descriptions match current implementation

---

## Implementation Order

1. **Delete outdated files first** - Remove confusion sources
2. **Update README.md** - Main user-facing documentation
3. **Update docs/README.md** - Fix broken links to deleted files
4. **Update frontend/README.md** - Update to Event Sourcing API
5. **Update ARCHITECTURE.md** - Core technical architecture
6. **Update PROJECT_CONTEXT.md** - Add ADRs and update patterns
7. **Review API.md** - Minor updates only
8. **Update SETUP.md** - Add PostgreSQL setup
9. **Commit all changes** - Single comprehensive commit

---

## Success Criteria

After implementation:

1. ✅ Claude Code receives accurate Event Sourcing context
2. ✅ No outdated documentation exists
3. ✅ All examples work with current codebase
4. ✅ Architecture decisions documented in ADRs
5. ✅ Migration history preserved
6. ✅ New developers can understand current system architecture
7. ✅ No references to superseded features (in-memory, Visualization CRUD)

---

## References

- Current Event Sourcing spec: `docs/superpowers/specs/2026-05-02-event-sourcing-agent-diary-design.md`
- MBTI projection spec: `docs/superpowers/specs/2026-05-03-mbti-projection-design.md`
- Current API documentation: `docs/API.md`
- Backend implementation: `backend/storage/postgres.go`
- Database schema: `backend/migrations/001_create_event_sourcing_schema.sql`
