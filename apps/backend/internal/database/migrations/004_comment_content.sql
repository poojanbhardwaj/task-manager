-- Keep the comments schema consistent with the repository, for both fresh and existing databases.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = current_schema()
          AND table_name = 'todo_comments'
          AND column_name = 'comment'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = current_schema()
          AND table_name = 'todo_comments'
          AND column_name = 'content'
    ) THEN
        ALTER TABLE todo_comments RENAME COLUMN comment TO content;
    END IF;
END $$;
---- create above / drop below ----
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = current_schema()
          AND table_name = 'todo_comments'
          AND column_name = 'content'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = current_schema()
          AND table_name = 'todo_comments'
          AND column_name = 'comment'
    ) THEN
        ALTER TABLE todo_comments RENAME COLUMN content TO comment;
    END IF;
END $$;