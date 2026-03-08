-- Migration: 009_inventory_phase2_vendor_purchase
-- Vendor, Purchase, PurchaseItem entities; InventoryLog source_type and source_id.

-- ====================================
-- UP Migration
-- ====================================

-- Create table: stich.Vendors
CREATE TABLE IF NOT EXISTS stich."Vendors" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  name TEXT NOT NULL,
  contact_person TEXT,
  phone TEXT,
  email TEXT,
  address TEXT,
  payment_terms TEXT,
  notes TEXT,
  PRIMARY KEY (id)
);

-- Create table: stich.Purchases
CREATE TABLE IF NOT EXISTS stich."Purchases" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  vendor_id INTEGER NOT NULL,
  purchase_number TEXT,
  purchase_date TIMESTAMPTZ NOT NULL,
  status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
  expected_delivery_date TIMESTAMPTZ,
  notes TEXT,
  total_amount DECIMAL(12,2) DEFAULT 0,
  paid_amount DECIMAL(12,2) DEFAULT 0,
  paid_at TIMESTAMPTZ,
  payment_method TEXT,
  PRIMARY KEY (id)
);

ALTER TABLE stich."Purchases"
  ADD CONSTRAINT fk_Purchase_vendor_id
  FOREIGN KEY (vendor_id) REFERENCES stich."Vendors" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

-- Create table: stich.PurchaseItems
CREATE TABLE IF NOT EXISTS stich."PurchaseItems" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  purchase_id INTEGER NOT NULL,
  product_id INTEGER NOT NULL,
  quantity_ordered INTEGER NOT NULL,
  quantity_received INTEGER NOT NULL DEFAULT 0,
  unit_cost DECIMAL(10,2) NOT NULL,
  line_total DECIMAL(12,2) DEFAULT 0,
  PRIMARY KEY (id)
);

ALTER TABLE stich."PurchaseItems"
  ADD CONSTRAINT fk_PurchaseItem_purchase_id
  FOREIGN KEY (purchase_id) REFERENCES stich."Purchases" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE stich."PurchaseItems"
  ADD CONSTRAINT fk_PurchaseItem_product_id
  FOREIGN KEY (product_id) REFERENCES stich."Products" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

-- Alter InventoryLogs: add source_type and source_id
ALTER TABLE stich."InventoryLogs"
  ADD COLUMN IF NOT EXISTS source_type VARCHAR(30) DEFAULT 'MANUAL';

ALTER TABLE stich."InventoryLogs"
  ADD COLUMN IF NOT EXISTS source_id INTEGER;

-- ====================================
-- DOWN Migration (Rollback)
-- ====================================

-- ALTER TABLE stich."InventoryLogs" DROP COLUMN IF EXISTS source_id;
-- ALTER TABLE stich."InventoryLogs" DROP COLUMN IF EXISTS source_type;
-- DROP TABLE IF EXISTS stich."PurchaseItems";
-- DROP TABLE IF EXISTS stich."Purchases";
-- DROP TABLE IF EXISTS stich."Vendors";
