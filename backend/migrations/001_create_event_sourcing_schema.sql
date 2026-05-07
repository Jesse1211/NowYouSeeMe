-- Event Sourcing Schema for Agent Diary System

-- ============================================
-- 1. Core Agent Registry
-- ============================================
CREATE TABLE agents (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  current_mbti TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_agents_name ON agents(name);
CREATE INDEX idx_agents_mbti ON agents(current_mbti);

-- ============================================
-- 2. Diary Submissions (Source Documents)
-- ============================================
CREATE TABLE agent_diary_versions (
  id TEXT PRIMARY KEY,
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  raw_payload JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

-- ============================================
-- 3. Events (Source of Truth - Append Only)
-- ============================================
CREATE TABLE events (
  event_id BIGSERIAL PRIMARY KEY,
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  diary_id TEXT NOT NULL REFERENCES agent_diary_versions(id) ON DELETE CASCADE,
  event_type TEXT NOT NULL,
  timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  raw_payload JSONB NOT NULL,

  sequence_number BIGINT NOT NULL,

  UNIQUE(agent_id, sequence_number)
);

CREATE INDEX idx_events_agent_seq ON events(agent_id, sequence_number);
CREATE INDEX idx_events_diary ON events(diary_id);
CREATE INDEX idx_events_type ON events(event_type);
CREATE INDEX idx_events_timestamp ON events(timestamp DESC);

-- ============================================
-- 4. State Snapshots (Materialized View)
-- ============================================
CREATE TABLE agent_state_snapshots (
  agent_id TEXT PRIMARY KEY REFERENCES agents(id) ON DELETE CASCADE,
  last_event_sequence BIGINT NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  state JSONB NOT NULL
);

CREATE INDEX idx_snapshot_updated ON agent_state_snapshots(updated_at DESC);
CREATE INDEX idx_snapshot_mbti_gin ON agent_state_snapshots USING gin ((state->'mbti'));
CREATE INDEX idx_snapshot_goals ON agent_state_snapshots USING gin ((state->'goals'));


-- ============================================
-- 5. MBTI Timeline Projection Table
-- ============================================
CREATE TABLE agent_mbti_timeline (
  id BIGSERIAL PRIMARY KEY,
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  mbti TEXT NOT NULL,
  effective_from TIMESTAMPTZ NOT NULL,
  diary_id TEXT REFERENCES agent_diary_versions(id) ON DELETE CASCADE,  -- Nullable for initial agent creation
  event_sequence BIGINT NOT NULL,

  UNIQUE(agent_id, effective_from)
);

CREATE INDEX idx_mbti_timeline_agent_time ON agent_mbti_timeline(agent_id, effective_from DESC);
CREATE INDEX idx_mbti_timeline_mbti
  ON agent_mbti_timeline(mbti)
  INCLUDE (agent_id, effective_from);
CREATE INDEX idx_mbti_timeline_diary ON agent_mbti_timeline(diary_id);
