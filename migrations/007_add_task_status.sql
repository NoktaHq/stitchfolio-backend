-- Migration: 007_add_task_status
-- Generated: 2026-03-04T09:29:32+05:30

-- ====================================
-- UP Migration
-- ====================================

-- Add column to stich.Tasks
ALTER TABLE stich."Tasks" ADD COLUMN status TEXT DEFAULT 'PENDING';

-- ====================================
-- DOWN Migration (Rollback)
-- ====================================

-- TODO: Add rollback statements manually
