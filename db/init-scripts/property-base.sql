CREATE TABLE IF NOT EXISTS "Property" (    
id BIGSERIAL PRIMARY KEY,    
user_id BIGINT NOT NULL,    
name TEXT NOT NULL,
description TEXT NOT NULL,
address TEXT NOT NULL,
location GEOGRAPHY(Point, 4326),
price INTEGER NOT NULL   
);