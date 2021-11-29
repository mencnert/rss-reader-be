CREATE TABLE IF NOT EXISTS rss(
    rss_id serial PRIMARY KEY,
    url VARCHAR(50) UNIQUE NOT NULL,
    rank INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    viewed BOOLEAN DEFAULT false,
    saved BOOLEAN DEFAULT false,
    last_fetch TIMESTAMP NOT NULL
);
