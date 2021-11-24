CREATE TABLE IF NOT EXISTS rss(
    rss_id serial PRIMARY KEY,
    url VARCHAR(50) UNIQUE NOT NULL,
    rank INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    last_update TIMESTAMP NOT NULL
);
