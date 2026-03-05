CREATE TABLE IF NOT EXISTS public.urls (
    url_id bigserial PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_url varchar(10) NOT NULL,
    created_at timestamptz DEFAULT NOW(),

    CONSTRAINT urls_original_url_key UNIQUE (original_url),
    CONSTRAINT urls_short_url_key UNIQUE (short_url)
);