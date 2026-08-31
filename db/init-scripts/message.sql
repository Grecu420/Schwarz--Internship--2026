CREATE TABLE IF NOT EXISTS messages (    
id BIGSERIAL PRIMARY KEY,    
conversation_id BIGINT NOT NULL,    
sender_id BIGINT NOT NULL,    
content TEXT NOT NULL,     
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);