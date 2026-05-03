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
