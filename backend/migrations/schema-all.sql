-- Event Sourcing Schema for Agent Diary System
-- Complete schema including all tables for Event-Driven Microservices architecture

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

-- ============================================
-- 6. Transactional Outbox Table
-- ============================================
-- Transactional Outbox Pattern: Events are written to this table
-- within the same transaction as the main events table, ensuring
-- at-least-once delivery guarantee to Redis Streams.

CREATE TABLE outbox (
  id BIGSERIAL PRIMARY KEY,

  -- Event reference
  event_id BIGINT NOT NULL REFERENCES events(event_id) ON DELETE CASCADE,
  aggregate_id TEXT NOT NULL,    -- agent_id (used for partitioning)
  event_type TEXT NOT NULL,
  payload JSONB NOT NULL,

  -- Publishing state
  published BOOLEAN DEFAULT FALSE NOT NULL,
  published_at TIMESTAMPTZ,

  -- Retry mechanism
  retry_count INT DEFAULT 0 NOT NULL,
  last_error TEXT,
  next_retry_at TIMESTAMPTZ,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Critical index: Find unpublished events efficiently
CREATE INDEX idx_outbox_pending
ON outbox(published, created_at)
WHERE published = FALSE;

-- Index for ordering by aggregate (maintains event order per agent)
CREATE INDEX idx_outbox_aggregate
ON outbox(aggregate_id, id);

-- Index for retry logic
CREATE INDEX idx_outbox_retry
ON outbox(next_retry_at)
WHERE published = FALSE AND next_retry_at IS NOT NULL;

COMMENT ON TABLE outbox IS 'Transactional outbox for reliable event publishing to Redis Streams';
COMMENT ON COLUMN outbox.aggregate_id IS 'Partition key - typically agent_id';
COMMENT ON COLUMN outbox.published IS 'FALSE = pending, TRUE = successfully published to message bus';

-- ============================================
-- 7. Gallery View (Denormalized Read Model)
-- ============================================
-- Denormalized read model for gallery page
-- Eliminates N+1 queries by pre-aggregating all data needed for gallery display
-- Updated by Gallery Projector consuming events from Redis Streams

CREATE TABLE gallery_view (
  agent_id TEXT PRIMARY KEY REFERENCES agents(id) ON DELETE CASCADE,
  agent_name TEXT NOT NULL,
  current_mbti TEXT NOT NULL,
  mbti_confidence DECIMAL(4,3),

  -- Latest metadata
  latest_mood TEXT,
  latest_philosophy TEXT,
  geometry_representation TEXT,

  -- Aggregated statistics
  total_goals INT DEFAULT 0 NOT NULL,
  completed_goals INT DEFAULT 0 NOT NULL,
  active_goals INT DEFAULT 0 NOT NULL,
  total_capabilities INT DEFAULT 0 NOT NULL,
  total_limitations INT DEFAULT 0 NOT NULL,
  total_aspirations INT DEFAULT 0 NOT NULL,

  -- Event sourcing metadata
  total_diaries INT DEFAULT 0 NOT NULL,
  last_event_sequence BIGINT NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_gallery_mbti ON gallery_view(current_mbti);
CREATE INDEX idx_gallery_updated ON gallery_view(updated_at DESC);
CREATE INDEX idx_gallery_goals ON gallery_view(total_goals DESC);

COMMENT ON TABLE gallery_view IS 'Denormalized read model for gallery page - updated by Gallery Projector';
COMMENT ON COLUMN gallery_view.last_event_sequence IS 'Tracks which events have been processed (idempotency)';
