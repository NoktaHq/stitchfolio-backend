-- Migration: 006_add_expense_detail_entity
-- Child table ExpenseDetails: Source, Price, ExpenseId (FK to Expenses).
-- Migrates existing Expenses into one ExpenseDetail per row: Source='', Price=Expenses.price.

-- Create table: stich.ExpenseDetails
CREATE TABLE IF NOT EXISTS stich."ExpenseDetails" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  source TEXT DEFAULT '',
  price DOUBLE PRECISION NOT NULL,
  expense_id BIGINT NOT NULL,
  PRIMARY KEY (id)
);

-- Foreign key to Expenses
ALTER TABLE stich."ExpenseDetails"
  ADD CONSTRAINT fk_ExpenseDetail_expense_id
  FOREIGN KEY (expense_id) REFERENCES stich."Expenses" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

-- Migrate existing data: one ExpenseDetail per Expense, Source='', Price=Expenses.price
INSERT INTO stich."ExpenseDetails" (
  created_at,
  updated_at,
  is_active,
  created_by_id,
  updated_by_id,
  channel_id,
  source,
  price,
  expense_id
)
SELECT
  E.created_at,
  E.updated_at,
  E.is_active,
  E.created_by_id,
  E.updated_by_id,
  E.channel_id,
  '' AS source,
  COALESCE(E.price, 0) AS price,
  E.id AS expense_id
FROM stich."Expenses" E
WHERE NOT EXISTS (
  SELECT 1 FROM stich."ExpenseDetails" ed WHERE ed.expense_id = E.id
);

-- ====================================
-- DOWN Migration (Rollback)
-- ====================================

-- DROP TABLE IF EXISTS stich."ExpenseDetails";
