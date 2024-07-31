ALTER TABLE subscribers
ADD CONSTRAINT email_unique UNIQUE (email);
