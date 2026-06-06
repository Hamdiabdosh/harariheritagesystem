CREATE TYPE movable_owner_type AS ENUM (
    'public',
    'government',
    'religion',
    'private'
);

CREATE TYPE storage_location AS ENUM (
    'museum',
    'store',
    'church',
    'private_home',
    'other'
);

CREATE TYPE movable_condition AS ENUM ('good', 'fair', 'damaged', 'incomplete');

CREATE TABLE movable_records (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id              VARCHAR(30) NOT NULL UNIQUE,
    registrar_id           UUID NOT NULL REFERENCES users(id),
    status                 record_status NOT NULL DEFAULT 'draft',

    name_amharic           VARCHAR(200) NOT NULL,
    name_local             VARCHAR(200),
    category               VARCHAR(10),

    location_name          VARCHAR(200),
    woreda                 VARCHAR(100),
    kebele                 VARCHAR(100),
    house_number           VARCHAR(50),
    current_use            VARCHAR(200),
    previous_id            VARCHAR(100),

    owner_type             movable_owner_type,
    owner_name             VARCHAR(200),
    storage_location       storage_location,
    storage_location_other VARCHAR(200),

    made_by                VARCHAR(200),
    period_made            VARCHAR(100),
    age_method             age_method,
    acquisition_methods    VARCHAR(30)[],

    height_cm              DECIMAL(8, 2),
    width_cm               DECIMAL(8, 2),
    length_cm              DECIMAL(8, 2),
    diameter_cm            DECIMAL(8, 2),
    thickness_cm           DECIMAL(8, 2),
    weight_kg              DECIMAL(8, 2),
    num_pages              INTEGER,
    num_chapters           INTEGER,
    num_illustrations      INTEGER,

    color_type             VARCHAR(200),
    has_decoration         BOOLEAN,
    materials              VARCHAR(30)[],
    material_other         VARCHAR(200),
    description            TEXT,
    notable_because        VARCHAR(30)[],
    notable_other          TEXT,
    significance           TEXT,

    condition              movable_condition,
    has_threat             BOOLEAN,
    threat_description     TEXT,
    maintenance_done       BOOLEAN,
    maintenance_by         VARCHAR(200),
    maintenance_date       DATE,
    maintenance_count      INTEGER,
    preventive_level       quality_level,
    accessibility          accessibility_level,
    notes                  TEXT,

    related_docs           VARCHAR(30)[],

    informant_name         VARCHAR(200),
    informant_sex          sex_type,
    informant_age          INTEGER,
    informant_occupation   VARCHAR(200),
    caretaker_name         VARCHAR(200),
    caretaker_role         VARCHAR(200),
    registrar_date         DATE,

    approved_at            TIMESTAMPTZ,
    approved_by            UUID REFERENCES users(id),
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER movable_records_set_updated_at
    BEFORE UPDATE ON movable_records
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();
