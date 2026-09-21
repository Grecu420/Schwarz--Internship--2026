CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS properties (    
id BIGSERIAL PRIMARY KEY,    
user_id BIGINT NOT NULL,    
name TEXT NOT NULL,
description TEXT NOT NULL,
address TEXT NOT NULL,
location GEOGRAPHY(Point, 4326),
price INTEGER NOT NULL,
image_urls TEXT[] NOT NULL,

CONSTRAINT fk_property_user 
    FOREIGN KEY (user_id) 
    REFERENCES users (id) 
    ON DELETE CASCADE

);