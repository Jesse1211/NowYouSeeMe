# MBTI Projection Table Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add MBTI projection table to optimize agent queries by MBTI type from O(N×E) to O(M) complexity.

**Architecture:** Create append-only `agent_mbti_timeline` table tracking MBTI changes. Update synchronously during diary submission. Use DISTINCT ON query to filter agents by current MBTI before fetching full state.

**Tech Stack:** PostgreSQL, Go, database/sql

---

## File Structure

### Modified Files:
- `backend/storage/postgres.go`: Add GetAgentIDsByCurrentMBTI, insertMBTITimeline, update CreateAgent and SubmitDiary
- `backend/api/gallery_handlers.go`: Update GetSnapshotsByMBTI to use new query method

### Created Files:
- `backend/migrations/002_add_mbti_timeline.sql`: Schema migration for new table

---

## Task 1: Create Migration File

**Files:**
- Create: `backend/migrations/002_add_mbti_timeline.sql`

- [ ] **Step 1: Create migration file**

```sql
-- backend/migrations/002_add_mbti_timeline.sql

-- ============================================
-- MBTI Timeline Projection Table
-- ============================================
CREATE TABLE agent_mbti_timeline (
  id BIGSERIAL PRIMARY KEY,
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  mbti TEXT NOT NULL,
  effective_from TIMESTAMPTZ NOT NULL,
  diary_id TEXT NOT NULL REFERENCES agent_diary_versions(id) ON DELETE CASCADE,
  event_sequence BIGINT NOT NULL,

  UNIQUE(agent_id, effective_from)
);

-- ============================================
-- Indexes
-- ============================================

-- Core index: Query latest MBTI by agent
CREATE INDEX idx_mbti_timeline_agent_time
  ON agent_mbti_timeline(agent_id, effective_from DESC);

-- Core index: Filter by MBTI type (with covering index)
CREATE INDEX idx_mbti_timeline_mbti
  ON agent_mbti_timeline(mbti)
  INCLUDE (agent_id, effective_from);

-- Foreign key index optimization
CREATE INDEX idx_mbti_timeline_diary
  ON agent_mbti_timeline(diary_id);
```

- [ ] **Step 2: Run migration**

```bash
cd /Users/liuzhenhua/Desktop/NAH/NowYouSeeMe
psql -U postgres -d nowyouseeme -f backend/migrations/002_add_mbti_timeline.sql
```

Expected: Tables and indexes created successfully

- [ ] **Step 3: Verify migration**

```bash
psql -U postgres -d nowyouseeme -c "\d agent_mbti_timeline"
```

Expected: Table structure displayed with all columns and indexes

- [ ] **Step 4: Commit migration file**

```bash
git add backend/migrations/002_add_mbti_timeline.sql
git commit -m "feat: add MBTI timeline projection table migration

Add agent_mbti_timeline table to track MBTI changes as checkpoints.
Enables O(M) query complexity for filtering agents by MBTI type.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task 2: Add insertMBTITimeline Storage Method

**Files:**
- Modify: `backend/storage/postgres.go` (add new method after line 438)

- [ ] **Step 1: Add insertMBTITimeline method**

Add this method after the SubmitDiary function in `backend/storage/postgres.go`:

```go
// insertMBTITimeline inserts a new MBTI timeline record
func (s *PostgresStore) insertMBTITimeline(
	tx *sql.Tx,
	agentID string,
	mbti string,
	diaryID string,
	eventSequence int64,
) error {
	query := `
		INSERT INTO agent_mbti_timeline (
			agent_id, mbti, effective_from, diary_id, event_sequence
		) VALUES ($1, $2, NOW(), $3, $4)
	`

	_, err := tx.Exec(query, agentID, mbti, diaryID, eventSequence)
	if err != nil {
		return fmt.Errorf("failed to insert MBTI timeline: %w", err)
	}

	return nil
}
```

- [ ] **Step 2: Verify code compiles**

```bash
cd backend
go build
```

Expected: No compilation errors

- [ ] **Step 3: Commit insertMBTITimeline method**

```bash
git add backend/storage/postgres.go
git commit -m "feat: add insertMBTITimeline storage method

Add method to insert MBTI timeline records during diary submission.
Only called when MBTI actually changes.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task 3: Update CreateAgent to Insert Initial MBTI

**Files:**
- Modify: `backend/storage/postgres.go:24-44` (CreateAgent method)

- [ ] **Step 1: Update CreateAgent to use transaction and insert initial MBTI**

Replace the CreateAgent method (lines 24-44) with:

```go
// CreateAgent creates a new agent
func (s *PostgresStore) CreateAgent(req *models.CreateAgentRequest) (*models.Agent, error) {
	agent := &models.Agent{
		ID:          req.AgentID,
		Name:        req.Name,
		InitialMBTI: req.InitialMBTI,
		CreatedAt:   time.Now(),
	}

	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert agent record
	query := `
		INSERT INTO agents (id, name, initial_mbti, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err = tx.Exec(query, agent.ID, agent.Name, agent.InitialMBTI, agent.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	// Insert initial MBTI timeline record
	err = s.insertMBTITimeline(tx, agent.ID, agent.InitialMBTI, "initial", 0)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return agent, nil
}
```

- [ ] **Step 2: Verify code compiles**

```bash
cd backend
go build
```

Expected: No compilation errors

- [ ] **Step 3: Test CreateAgent with timeline insertion**

```bash
cd backend
go test -run TestCreateAgent ./storage/... -v
```

Expected: Tests pass (or create test if needed)

- [ ] **Step 4: Commit updated CreateAgent**

```bash
git add backend/storage/postgres.go
git commit -m "feat: update CreateAgent to insert initial MBTI timeline

CreateAgent now inserts initial MBTI record to agent_mbti_timeline
with diary_id='initial' and event_sequence=0.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task 4: Update SubmitDiary to Track MBTI Changes

**Files:**
- Modify: `backend/storage/postgres.go:364-438` (SubmitDiary method)

- [ ] **Step 1: Update SubmitDiary to track old MBTI and insert timeline on change**

Replace lines 378-381 in SubmitDiary with:

```go
	// Build current state using WAL
	currentState, currentSeq, err := s.GetCurrentState(agentID)
	if err != nil {
		return nil, err
	}

	// Track old MBTI before applying changes
	oldMBTI := currentState.MBTI
```

Then add after line 420 (after updating currentState.CurrentSelfReflection):

```go
	currentState.CurrentSelfReflection = payload.SelfReflection

	// Check if MBTI changed and insert timeline record
	newMBTI := payload.MBTI
	if newMBTI != oldMBTI {
		finalSeq := nextSeq + int64(len(payload.Operations)) - 1
		if len(payload.Operations) == 0 {
			finalSeq = currentSeq
		}

		err = s.insertMBTITimeline(tx, agentID, newMBTI, diaryID, finalSeq)
		if err != nil {
			return nil, err
		}
	}

	// Materialize snapshot (MVP: always)
```

- [ ] **Step 2: Verify code compiles**

```bash
cd backend
go build
```

Expected: No compilation errors

- [ ] **Step 3: Test SubmitDiary with MBTI change detection**

```bash
cd backend
go test -run TestSubmitDiary ./storage/... -v
```

Expected: Tests pass

- [ ] **Step 4: Commit updated SubmitDiary**

```bash
git add backend/storage/postgres.go
git commit -m "feat: update SubmitDiary to track MBTI changes

SubmitDiary now compares old and new MBTI values.
Only inserts timeline record when MBTI actually changes.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task 5: Add GetAgentIDsByCurrentMBTI Query Method

**Files:**
- Modify: `backend/storage/postgres.go` (add new method after insertMBTITimeline)

- [ ] **Step 1: Add GetAgentIDsByCurrentMBTI method**

Add this method after insertMBTITimeline:

```go
// GetAgentIDsByCurrentMBTI returns agent IDs with the specified current MBTI type
func (s *PostgresStore) GetAgentIDsByCurrentMBTI(mbtiType string) ([]string, error) {
	query := `
		SELECT agent_id
		FROM (
			SELECT DISTINCT ON (agent_id) agent_id, mbti
			FROM agent_mbti_timeline
			ORDER BY agent_id, effective_from DESC
		) latest
		WHERE mbti = $1
	`

	rows, err := s.db.Query(query, mbtiType)
	if err != nil {
		return nil, fmt.Errorf("failed to query agents by MBTI: %w", err)
	}
	defer rows.Close()

	agentIDs := []string{}
	for rows.Next() {
		var agentID string
		if err := rows.Scan(&agentID); err != nil {
			return nil, err
		}
		agentIDs = append(agentIDs, agentID)
	}

	return agentIDs, nil
}
```

- [ ] **Step 2: Verify code compiles**

```bash
cd backend
go build
```

Expected: No compilation errors

- [ ] **Step 3: Commit GetAgentIDsByCurrentMBTI**

```bash
git add backend/storage/postgres.go
git commit -m "feat: add GetAgentIDsByCurrentMBTI query method

Query agent_mbti_timeline using DISTINCT ON to get agents
with specified current MBTI type. O(M) complexity.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task 6: Update GetSnapshotsByMBTI API Handler

**Files:**
- Modify: `backend/api/gallery_handlers.go:78-127` (GetSnapshotsByMBTI function)

- [ ] **Step 1: Read current GetSnapshotsByMBTI implementation**

```bash
cat backend/api/gallery_handlers.go | sed -n '78,127p'
```

Expected: See current implementation that loops through all agents

- [ ] **Step 2: Replace GetSnapshotsByMBTI implementation**

Replace lines 78-127 in `backend/api/gallery_handlers.go` with:

```go
// GetSnapshotsByMBTI filters agents by MBTI type
func GetSnapshotsByMBTI(store *storage.PostgresStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		mbtiType := c.Query("mbti")
		if mbtiType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "mbti query parameter required"})
			return
		}

		// Only query agents with matching current MBTI
		agentIDs, err := store.GetAgentIDsByCurrentMBTI(mbtiType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		type Result struct {
			AgentID   string             `json:"agent_id"`
			Snapshot  *models.AgentState `json:"snapshot"`
			UpdatedAt string             `json:"updated_at"`
		}

		results := []Result{}

		// Only fetch full state for matched agents
		for _, agentID := range agentIDs {
			snapshot, _, err := store.GetCurrentState(agentID)
			if err != nil {
				continue
			}

			snapshotDB, _ := store.GetSnapshot(agentID)
			updatedAt := ""
			if snapshotDB != nil {
				updatedAt = snapshotDB.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
			}

			results = append(results, Result{
				AgentID:   agentID,
				Snapshot:  snapshot,
				UpdatedAt: updatedAt,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"mbti":   mbtiType,
			"agents": results,
			"count":  len(results),
		})
	}
}
```

- [ ] **Step 3: Verify code compiles**

```bash
cd backend
go build
```

Expected: No compilation errors

- [ ] **Step 4: Commit updated API handler**

```bash
git add backend/api/gallery_handlers.go
git commit -m "feat: optimize GetSnapshotsByMBTI using projection table

Use GetAgentIDsByCurrentMBTI to filter agents before fetching state.
Reduces complexity from O(N×E) to O(M) where M = matching agents.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task 7: Test End-to-End Flow

**Files:**
- Test: Manual API testing

- [ ] **Step 1: Start backend server**

```bash
cd backend
make run
```

Expected: Server starts on port 8080

- [ ] **Step 2: Create test agent**

```bash
curl -X POST http://localhost:8080/api/v1/agents \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "test_agent_mbti_timeline",
    "name": "TestMBTIAgent",
    "initial_mbti": "INTP-A"
  }'
```

Expected: 201 Created, agent data returned

- [ ] **Step 3: Verify initial MBTI timeline record**

```bash
psql -U postgres -d nowyouseeme -c \
  "SELECT agent_id, mbti, diary_id, event_sequence FROM agent_mbti_timeline WHERE agent_id = 'test_agent_mbti_timeline';"
```

Expected: One record with mbti='INTP-A', diary_id='initial', event_sequence=0

- [ ] **Step 4: Query by MBTI type**

```bash
curl http://localhost:8080/api/v1/snapshots?mbti=INTP-A
```

Expected: JSON response with agents array containing test_agent_mbti_timeline

- [ ] **Step 5: Submit diary with MBTI change**

```bash
curl -X POST http://localhost:8080/api/v1/diaries \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "test_agent_mbti_timeline",
    "mbti": "ENTP-A",
    "mbti_confidence": 0.85,
    "geometry_representation": "test_image_data",
    "current_mood": "excited",
    "philosophy": "testing MBTI changes",
    "self_reflection": {
      "rumination_for_yesterday": "was INTP",
      "what_happened_today": "evolved to ENTP",
      "expectations_for_tomorrow": "continue growing"
    },
    "operations": []
  }'
```

Expected: 201 Created, snapshot with ENTP-A returned

- [ ] **Step 6: Verify MBTI timeline has two records**

```bash
psql -U postgres -d nowyouseeme -c \
  "SELECT agent_id, mbti, effective_from, event_sequence FROM agent_mbti_timeline WHERE agent_id = 'test_agent_mbti_timeline' ORDER BY effective_from;"
```

Expected: Two records - first INTP-A, second ENTP-A

- [ ] **Step 7: Query by new MBTI type**

```bash
curl http://localhost:8080/api/v1/snapshots?mbti=ENTP-A
```

Expected: JSON response containing test_agent_mbti_timeline

- [ ] **Step 8: Query by old MBTI type**

```bash
curl http://localhost:8080/api/v1/snapshots?mbti=INTP-A
```

Expected: JSON response NOT containing test_agent_mbti_timeline (only current MBTI)

- [ ] **Step 9: Submit diary without MBTI change**

```bash
curl -X POST http://localhost:8080/api/v1/diaries \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "test_agent_mbti_timeline",
    "mbti": "ENTP-A",
    "mbti_confidence": 0.90,
    "geometry_representation": "test_image_data_2",
    "current_mood": "stable",
    "philosophy": "testing no MBTI change",
    "self_reflection": {
      "rumination_for_yesterday": "still ENTP",
      "what_happened_today": "no change",
      "expectations_for_tomorrow": "stay ENTP"
    },
    "operations": []
  }'
```

Expected: 201 Created, snapshot with ENTP-A returned

- [ ] **Step 10: Verify timeline still has only two records**

```bash
psql -U postgres -d nowyouseeme -c \
  "SELECT COUNT(*) FROM agent_mbti_timeline WHERE agent_id = 'test_agent_mbti_timeline';"
```

Expected: Count = 2 (no new record added since MBTI didn't change)

- [ ] **Step 11: Document test results**

Create test summary in commit message or notes showing:
- Initial MBTI timeline insertion works
- MBTI change detection works
- Timeline only updated on actual MBTI change
- Query filtering by MBTI works correctly

---

## Task 8: Performance Verification

**Files:**
- Test: Performance comparison

- [ ] **Step 1: Generate test data with fake agents script**

```bash
cd sdk
python scripts/generate_fake_agents.py -n 50 -e 5 -q
```

Expected: 50 agents created with varying MBTI types

- [ ] **Step 2: Check timeline table size**

```bash
psql -U postgres -d nowyouseeme -c \
  "SELECT COUNT(*) as total_timeline_records, COUNT(DISTINCT agent_id) as unique_agents FROM agent_mbti_timeline;"
```

Expected: Shows total timeline records and unique agents

- [ ] **Step 3: Test query performance with EXPLAIN ANALYZE**

```bash
psql -U postgres -d nowyouseeme -c \
  "EXPLAIN ANALYZE
   SELECT agent_id FROM (
     SELECT DISTINCT ON (agent_id) agent_id, mbti
     FROM agent_mbti_timeline
     ORDER BY agent_id, effective_from DESC
   ) latest WHERE mbti = 'INTP-A';"
```

Expected: Execution time < 10ms, uses idx_mbti_timeline_mbti index

- [ ] **Step 4: Time API response**

```bash
time curl -s http://localhost:8080/api/v1/snapshots?mbti=INTP-A > /dev/null
```

Expected: Response time < 200ms for 50 agents

- [ ] **Step 5: Compare with timeline table disabled (simulate old behavior)**

Count agents manually:
```bash
psql -U postgres -d nowyouseeme -c \
  "SELECT COUNT(*) FROM agents;"
```

Calculate expected difference: With 50 agents, old method = 50+ queries, new method = ~5-10 queries (depending on MBTI distribution)

- [ ] **Step 6: Document performance improvement**

Note:
- Old complexity: O(N×E) = 50 agents × avg events
- New complexity: O(M) = ~8 matching agents (assuming ~16% per MBTI type)
- Query reduction: ~90% fewer database calls

---

## Task 9: Final Integration Test

**Files:**
- Test: Complete workflow test

- [ ] **Step 1: Test with frontend (if running)**

If frontend is running:
```bash
# Open browser to http://localhost:5173/mbti/INTP-A
# Verify agents are displayed
# Check network tab shows single /snapshots?mbti=INTP-A request
```

Expected: Agents with current MBTI=INTP-A displayed correctly

- [ ] **Step 2: Verify backward compatibility**

Old API calls should still work:
```bash
curl http://localhost:8080/api/v1/gallery
```

Expected: All agents returned with snapshots

- [ ] **Step 3: Test error cases**

Invalid MBTI type:
```bash
curl http://localhost:8080/api/v1/snapshots?mbti=INVALID
```

Expected: Empty agents array (valid response, just no matches)

Missing MBTI parameter:
```bash
curl http://localhost:8080/api/v1/snapshots
```

Expected: 400 Bad Request

- [ ] **Step 4: Commit final integration test results**

```bash
git add -A
git commit -m "test: verify MBTI projection table integration

Tested:
- CreateAgent inserts initial timeline record
- SubmitDiary tracks MBTI changes
- GetSnapshotsByMBTI uses projection table
- Performance improvement: O(N×E) → O(M)
- Backward compatibility maintained

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task 10: Documentation Update

**Files:**
- Modify: `backend/README.md` or create `docs/mbti-projection.md`

- [ ] **Step 1: Document MBTI timeline table purpose**

Add to backend documentation:

```markdown
## MBTI Projection Table

The `agent_mbti_timeline` table is a projection (read model) that tracks MBTI type changes for agents.

### Purpose
- Optimize queries for agents by MBTI type
- Reduce complexity from O(N×E) to O(M) where M = matching agents
- Maintain append-only event sourcing principles

### How It Works
1. Initial MBTI inserted when agent is created
2. New record inserted only when MBTI actually changes
3. Query uses DISTINCT ON to get current MBTI per agent
4. Filter agents before fetching full state

### Schema
```sql
CREATE TABLE agent_mbti_timeline (
  id BIGSERIAL PRIMARY KEY,
  agent_id TEXT NOT NULL,
  mbti TEXT NOT NULL,
  effective_from TIMESTAMPTZ NOT NULL,
  diary_id TEXT NOT NULL,
  event_sequence BIGINT NOT NULL
);
```

### API Usage
```bash
GET /api/v1/snapshots?mbti=INTP-A
```

Returns agents with current MBTI type = INTP-A.
```
```

- [ ] **Step 2: Update migration documentation**

Document migration in `backend/migrations/README.md` (or create if doesn't exist):

```markdown
## 002_add_mbti_timeline.sql

**Date**: 2026-05-03
**Purpose**: Add MBTI projection table for performance optimization

### Changes
- Creates `agent_mbti_timeline` table
- Adds indexes for efficient querying
- Enables O(M) complexity for MBTI filtering

### Deployment
```bash
psql -U postgres -d nowyouseeme -f backend/migrations/002_add_mbti_timeline.sql
```

### Rollback
```sql
DROP TABLE agent_mbti_timeline CASCADE;
```
```

- [ ] **Step 3: Commit documentation**

```bash
git add backend/README.md docs/mbti-projection.md backend/migrations/README.md
git commit -m "docs: add MBTI projection table documentation

Document purpose, schema, and usage of agent_mbti_timeline.
Include migration instructions and rollback procedure.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Success Criteria

✅ **Migration Applied**: agent_mbti_timeline table exists with correct schema and indexes
✅ **CreateAgent Updated**: Inserts initial MBTI timeline record
✅ **SubmitDiary Updated**: Tracks MBTI changes and inserts timeline records only on change
✅ **Query Optimized**: GetSnapshotsByMBTI uses projection table, O(M) complexity
✅ **Tests Pass**: All integration tests show correct behavior
✅ **Performance Improved**: 9x+ reduction in database queries for MBTI filtering
✅ **Backward Compatible**: Existing API endpoints work unchanged
✅ **Documented**: Purpose, schema, and usage clearly documented

---

## Spec Coverage Check

- [x] **Part 1 - Database Schema**: Task 1 creates migration with exact schema from spec
- [x] **Part 2 - Update Logic**: Tasks 2-4 implement insertMBTITimeline, CreateAgent, and SubmitDiary updates
- [x] **Part 3 - Query Implementation**: Tasks 5-6 add GetAgentIDsByCurrentMBTI and update API handler
- [x] **Part 4 - Integration**: Covered by existing transaction flow in SubmitDiary
- [x] **Part 5 - Migration**: Task 1 creates migration file
- [x] **Part 6 - Error Handling**: Inherent in transaction rollback, covered in testing
- [x] **Testing**: Tasks 7-9 provide comprehensive end-to-end testing
- [x] **Documentation**: Task 10 documents implementation

All spec requirements covered.
