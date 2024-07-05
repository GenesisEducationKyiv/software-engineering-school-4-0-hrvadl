SET sql_log_bin = 0;
CREATE TABLE subscribers (
  id int PRIMARY KEY AUTO_INCREMENT,
  email varchar(50) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE subscribers
ADD CONSTRAINT UC_subscriber_email UNIQUE (email);
SET sql_log_bin = 1;
