CREATE TABLE IF NOT EXISTS public.staging_followers (
    email VARCHAR(255),
    full_name VARCHAR(255) NOT NULL,
    load_id UUID NOT NULL
);
