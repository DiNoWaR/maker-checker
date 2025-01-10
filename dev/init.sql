CREATE TABLE IF NOT EXISTS messages (   id TEXT PRIMARY KEY,
                                            sender_id TEXT NOT NULL,
                                            recipient_id TEXT NOT NULL,
                                            status TEXT,
                                            content TEXT NOT NULL,
                                            ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP);

CREATE INDEX idx_messages_status ON messages(status);