-- Test accounts for each actor role (manager seeded in 009).
-- Password for both: Test1234
INSERT INTO users (
    full_name,
    email,
    password_hash,
    role,
    language,
    is_active
) VALUES
(
    'Test Registrar',
    'registrar@qirsmezgeb.gov.et',
    '$2a$12$CRO6MShI7JIey4WrxUUDGOJpsRg3ltG7ruDK4GBfx6on8A4Tdq8DS',
    'registrar',
    'am',
    TRUE
),
(
    'Test Supervisor',
    'supervisor@qirsmezgeb.gov.et',
    '$2a$12$CRO6MShI7JIey4WrxUUDGOJpsRg3ltG7ruDK4GBfx6on8A4Tdq8DS',
    'supervisor',
    'en',
    TRUE
);
