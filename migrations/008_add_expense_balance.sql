-- Migration: 008_add_expense_balance
-- Generated: 2026-03-04T09:30:26+05:30

-- ====================================
-- UP Migration
-- ====================================

-- Add column to stich.Expenses
ALTER TABLE stich."Expenses" ADD COLUMN balance DOUBLE PRECISION;

-- ====================================
-- DOWN Migration (Rollback)
-- ====================================

-- TODO: Add rollback statements manually
