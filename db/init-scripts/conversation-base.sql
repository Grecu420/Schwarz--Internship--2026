CREATE TABLE IF NOT EXISTS conversations (    
    id BIGSERIAL PRIMARY KEY,    
    user1_id BIGINT NOT NULL,    
    user2_id BIGINT NOT NULL,    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
); 

CREATE UNIQUE INDEX IF NOT EXISTS CONVERSATION_USER_ORDER_IDX ON conversations (LEAST(user1_id, user2_id), GREATEST(user1_id, user2_id));


-- -- 1. Create the trigger function to update the parent conversation
-- CREATE OR REPLACE FUNCTION touch_parent_conversation()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     UPDATE conversations
--     SET updated_at = NOW()
--     WHERE id = NEW.conversation_id;
    
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- -- 2. Attach the trigger to the "Message" table on INSERT
-- CREATE TRIGGER trigger_touch_conversation_on_message
-- AFTER INSERT ON messages
-- FOR EACH ROW
-- EXECUTE FUNCTION touch_parent_conversation();