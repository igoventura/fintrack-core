-- Write your migrate up statements here

INSERT INTO tenants (id, name) VALUES
('00000000-0000-0000-0000-000000000001', 'Default Tenant')

---- create above / drop below ----

DELETE FROM tenants WHERE id = '00000000-0000-0000-0000-000000000001';
