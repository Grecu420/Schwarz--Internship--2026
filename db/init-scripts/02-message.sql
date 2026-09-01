CREATE TABLE IF NOT EXISTS messages (    
id BIGSERIAL PRIMARY KEY,    
conversation_id BIGINT NOT NULL,    
sender_id BIGINT NOT NULL,    
content TEXT NOT NULL,     
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- -- 1. Create the trigger function to update the parent conversation
CREATE OR REPLACE FUNCTION touch_parent_conversation()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE conversations
    SET updated_at = NOW()
    WHERE id = NEW.conversation_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- -- 2. Attach the trigger to the "Message" table on INSERT
CREATE TRIGGER trigger_touch_conversation_on_message
AFTER INSERT ON messages
FOR EACH ROW
EXECUTE FUNCTION touch_parent_conversation();