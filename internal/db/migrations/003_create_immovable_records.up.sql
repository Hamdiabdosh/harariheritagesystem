CREATE TYPE record_status AS ENUM (
    'draft',
    'pending_review',
    'under_review',
    'returned',
    'approved'
);

CREATE TYPE immovable_owner_type AS ENUM (
    'public',
    'government',
    'religion',
    'private',
    'waqf'
);

CREATE TYPE age_method AS ENUM ('estimated', 'exact', 'relative');

CREATE TYPE overall_condition AS ENUM (
    'very_good',
    'good',
    'damaged',
    'severely_damaged'
);

CREATE TYPE damage_level AS ENUM ('minor', 'moderate', 'medium', 'severe');

CREATE TYPE quality_level AS ENUM (
    'very_good',
    'good',
    'medium',
    'low',
    'very_low'
);

CREATE TYPE accessibility_level AS ENUM (
    'very_good',
    'good',
    'medium',
    'low',
    'very_low',
    'none'
);

CREATE TYPE sex_type AS ENUM ('male', 'female');

CREATE TYPE record_type AS ENUM ('immovable', 'movable');

CREATE TABLE immovable_records (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id             VARCHAR(30) NOT NULL UNIQUE,
    registrar_id          UUID NOT NULL REFERENCES users(id),
    status                record_status NOT NULL DEFAULT 'draft',

    name_amharic          VARCHAR(200) NOT NULL,
    name_local            VARCHAR(200),
    category              VARCHAR(10)[],
    current_use           VARCHAR(50)[],
    current_use_other     VARCHAR(200),
    previous_id           VARCHAR(100),

    woreda                VARCHAR(100) NOT NULL,
    kebele                VARCHAR(100) NOT NULL,
    house_number          VARCHAR(50),
    street_number         VARCHAR(50),
    gate                  VARCHAR(50),

    owner_type            immovable_owner_type,
    owner_name            VARCHAR(200),
    map_reference         VARCHAR(200),
    gps_east              DECIMAL(10, 6),
    gps_north             DECIMAL(10, 6),
    elevation_m           DECIMAL(8, 2),

    built_by              VARCHAR(200),
    construction_period   VARCHAR(100),
    age_method            age_method,
    height_m              DECIMAL(8, 2),
    length_m              DECIMAL(8, 2),
    width_m               DECIMAL(8, 2),
    num_doors             INTEGER,
    num_windows           INTEGER,
    num_rooms             INTEGER,
    material              TEXT,
    description           TEXT,
    harari_house_grades   VARCHAR(50)[],
    neighborhood_type     VARCHAR(50),

    overall_condition     overall_condition,
    damage_roof           damage_level,
    damage_cornice        damage_level,
    damage_wall           damage_level,
    damage_floor          damage_level,
    damage_door           damage_level,
    damage_cupboard       damage_level,
    damage_upper_floor    damage_level,
    damage_dera           damage_level,
    damage_pillar         damage_level,

    value_historical      TEXT,
    value_craftsmanship   TEXT,
    value_artistic        TEXT,
    value_scientific      TEXT,
    value_cultural        TEXT,

    has_threat            BOOLEAN,
    maintenance_done      BOOLEAN,
    maintenance_reason    VARCHAR(300),
    maintenance_by        VARCHAR(200),
    maintenance_date      DATE,
    maintenance_count     INTEGER,
    preventive_level      quality_level,
    accessibility         accessibility_level,
    notes                 TEXT,

    related_docs          VARCHAR(30)[],
    has_oral_history      BOOLEAN,

    caretaker_name        VARCHAR(200),
    caretaker_role        VARCHAR(200),
    informant_name        VARCHAR(200),
    informant_sex         sex_type,
    informant_age         INTEGER,
    registrar_date        DATE,

    approved_at           TIMESTAMPTZ,
    approved_by           UUID REFERENCES users(id),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER immovable_records_set_updated_at
    BEFORE UPDATE ON immovable_records
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();
