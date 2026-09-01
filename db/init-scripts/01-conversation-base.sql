CREATE TABLE IF NOT EXISTS conversations (    
    id BIGSERIAL PRIMARY KEY,    
    user1_id BIGINT NOT NULL,    
    user2_id BIGINT NOT NULL,    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
); 

CREATE UNIQUE INDEX IF NOT EXISTS CONVERSATION_USER_ORDER_IDX ON conversations (LEAST(user1_id, user2_id), GREATEST(user1_id, user2_id));


