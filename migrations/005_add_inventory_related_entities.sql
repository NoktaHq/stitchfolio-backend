-- Migration: 005_add_inventory_related_entities
-- Generated: 2026-02-15T04:58:52+05:30

-- ====================================
-- UP Migration
-- ====================================

-- Create table: stich.Inventories
CREATE TABLE IF NOT EXISTS stich."Inventories" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  product_id INTEGER NOT NULL UNIQUE,
  quantity INTEGER NOT NULL DEFAULT 0,
  low_stock_threshold INTEGER DEFAULT 0,
  PRIMARY KEY (id)
);

-- Create table: stich.InventoryLogs
CREATE TABLE IF NOT EXISTS stich."InventoryLogs" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  product_id INTEGER NOT NULL,
  change_type VARCHAR(20) NOT NULL,
  quantity INTEGER NOT NULL,
  reason TEXT NOT NULL,
  notes TEXT,
  logged_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (id)
);

-- Create table: stich.Products
CREATE TABLE IF NOT EXISTS stich."Products" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  name TEXT NOT NULL,
  sku TEXT UNIQUE,
  category_id INTEGER NOT NULL,
  description TEXT,
  cost_price DECIMAL(10,2) NOT NULL,
  selling_price DECIMAL(10,2) NOT NULL,
  PRIMARY KEY (id)
);

-- Create table: stich.Categories
CREATE TABLE IF NOT EXISTS stich."Categories" (
  id BIGSERIAL NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  is_active BOOL DEFAULT true,
  created_by_id INTEGER,
  updated_by_id INTEGER,
  channel_id INTEGER,
  name TEXT NOT NULL,
  PRIMARY KEY (id)
);

-- Add foreign key to stich.Inventories
ALTER TABLE stich."Inventories" ADD CONSTRAINT fk_Inventory_product_id FOREIGN KEY (product_id) REFERENCES stich."Products" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

-- Add foreign key to stich.InventoryLogs
ALTER TABLE stich."InventoryLogs" ADD CONSTRAINT fk_InventoryLog_product_id FOREIGN KEY (product_id) REFERENCES stich."Products" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

-- Add foreign key to stich.Products
ALTER TABLE stich."Products" ADD CONSTRAINT fk_Product_category_id FOREIGN KEY (category_id) REFERENCES stich."Categories" (id) ON DELETE RESTRICT ON UPDATE RESTRICT;


-- ====================================
-- DOWN Migration (Rollback)
-- ====================================

-- TODO: Add rollback statements manually
