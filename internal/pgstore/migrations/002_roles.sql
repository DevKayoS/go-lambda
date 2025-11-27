-- Write your migrate up statements here
CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE permissions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Tabela intermediária: roles podem ter várias permissões
CREATE TABLE role_permissions (
    role_id BIGINT REFERENCES roles(id) ON DELETE CASCADE,
    permission_id BIGINT REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- Adiciona role_id na tabela users
ALTER TABLE users ADD COLUMN role_id BIGINT REFERENCES roles(id);

-- Insere roles padrão
INSERT INTO roles (name, description) VALUES 
    ('admin', 'Administrator with full access'),
    ('user', 'Regular user with limited access'),
    ('moderator', 'Moderator with some admin capabilities');

-- Insere permissões padrão
INSERT INTO permissions (name, description) VALUES 
    ('read:transactions', 'Can view transactions'),
    ('write:transactions', 'Can create transactions'),
    ('delete:transactions', 'Can delete transactions'),
    ('read:users', 'Can view users'),
    ('write:users', 'Can create users'),
    ('delete:users', 'Can delete users'),
    ('manage:all', 'Full system access');

-- Associa permissões aos roles
-- Admin tem tudo
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'admin';

-- User tem apenas leitura e escrita de transactions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'user' AND p.name IN ('read:transactions', 'write:transactions');

-- Moderator tem mais permissões
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'moderator' AND p.name IN (
    'read:transactions', 'write:transactions', 'delete:transactions', 
    'read:users', 'write:users'
);

---- create above / drop below ----
ALTER TABLE users DROP COLUMN IF EXISTS role_id;

-- Remove as associações de permissões aos roles
DELETE FROM role_permissions;

-- Remove as permissões padrão
DELETE FROM permissions;

-- Remove os roles padrão
DELETE FROM roles;

-- Drop da tabela intermediária role_permissions
DROP TABLE IF EXISTS role_permissions;

-- Drop da tabela permissions
DROP TABLE IF EXISTS permissions;

-- Drop da tabela roles
DROP TABLE IF EXISTS roles;
