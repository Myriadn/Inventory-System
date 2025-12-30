-- schema.sql

-- 1. Setup Enum untuk Role (Agar konsisten sesuai requirement)
CREATE TYPE user_role AS ENUM ('super_admin', 'admin', 'staff');

-- 2. Tabel Users
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL, -- Akan menyimpan hash bcrypt
    role user_role NOT NULL DEFAULT 'staff',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Tabel Sessions (Untuk Authentication & Session Management)
-- Sesuai requirement: token UUID, expired_at, revoked_at
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- ID Sesi sekaligus Token
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token UUID NOT NULL DEFAULT gen_random_uuid(), -- Bisa pakai kolom ini atau ID sebagai token
    ip_address VARCHAR(45), -- Opsional: untuk security logging
    user_agent TEXT, -- Opsional: untuk mengetahui device user
    is_revoked BOOLEAN DEFAULT FALSE, -- Flag simple untuk cek status
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE -- Diisi ketika user logout
);

-- 4. Tabel Master: Categories
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 5. Tabel Master: Warehouses (Gudang)
CREATE TABLE warehouses (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 6. Tabel Master: Racks (Rak)
CREATE TABLE racks (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL, -- Misal: "Rak A-01"
    category VARCHAR(50), -- Opsional: Jenis rak
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 7. Tabel Products (Barang)
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    sku VARCHAR(50) NOT NULL UNIQUE, -- Stock Keeping Unit (Kode Unik Barang)
    name VARCHAR(150) NOT NULL,
    description TEXT,
    stock INT NOT NULL DEFAULT 0 CHECK (stock >= 0), -- Prevent stock negatif
    price DECIMAL(15, 2) NOT NULL DEFAULT 0.00, -- 15 digit total, 2 desimal

    -- Relasi ke Master Data
    category_id BIGINT REFERENCES categories(id) ON DELETE RESTRICT,
    warehouse_id BIGINT REFERENCES warehouses(id) ON DELETE RESTRICT,
    rack_id BIGINT REFERENCES racks(id) ON DELETE RESTRICT,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 8. Tabel Sales (Transaksi Header)
CREATE TABLE sales (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL, -- Siapa kasirnya
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 9. Tabel Sale Details (Item Transaksi)
CREATE TABLE sale_details (
    id BIGSERIAL PRIMARY KEY,
    sale_id BIGINT NOT NULL REFERENCES sales(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(15, 2) NOT NULL, -- Harga saat transaksi terjadi (Snapshot)
    subtotal DECIMAL(15, 2) NOT NULL -- quantity * unit_price
);

-- INDEXING (Untuk mempercepat query)
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_sessions_token ON sessions(token);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_sales_date ON sales(transaction_date);

-- SEED DATA (Data awal untuk Super Admin supaya bisa login pertama kali)
-- Password 'password123' di-hash menggunakan bcrypt (cost 10)
-- Hash ini mungkin berbeda tergantung generator, tapi ini valid untuk bcrypt standard.
INSERT INTO users (username, email, password_hash, role)
VALUES (
    'SuperAdmin',
    'superadmin@inventory.com',
    '$2a$10$2.d.w/M7.u7.t.k.h.0.u.e.r.0.h.a.s.h.e.d.p.a.s.s.w.o.r.d', -- GANTI INI DENGAN HASH ASLI DARI KODE GO KAMU NANTI
    'super_admin'
);
