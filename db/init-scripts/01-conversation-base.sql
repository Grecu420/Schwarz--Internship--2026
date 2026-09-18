CREATE TABLE IF NOT EXISTS conversations (    
    id BIGSERIAL PRIMARY KEY,    
    user1_id BIGINT NOT NULL,    
    user2_id BIGINT NOT NULL,    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_conversation_user1 
        FOREIGN KEY (user1_id) 
        REFERENCES users (id) 
        ON DELETE CASCADE,

    CONSTRAINT fk_conversation_user2 
        FOREIGN KEY (user2_id) 
        REFERENCES users (id) 
        ON DELETE CASCADE            
); 

CREATE UNIQUE INDEX IF NOT EXISTS CONVERSATION_USER_ORDER_IDX ON conversations (LEAST(user1_id, user2_id), GREATEST(user1_id, user2_id));


