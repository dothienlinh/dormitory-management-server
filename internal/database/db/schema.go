package db

const CreateEnumTypes = `
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('student', 'admin');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_gender') THEN
        CREATE TYPE user_gender AS ENUM ('male', 'female', 'other');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status') THEN
        CREATE TYPE user_status AS ENUM ('active', 'inactive', 'absent');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'room_status') THEN
        CREATE TYPE room_status AS ENUM ('available', 'occupied', 'maintenance');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'room_rent_status') THEN
        CREATE TYPE room_rent_status AS ENUM ('active', 'inactive', 'graduated');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'contract_status') THEN
        CREATE TYPE contract_status AS ENUM ('active', 'inactive');
    END IF;
END
$$;
`
