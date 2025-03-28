
CREATE TABLE IF NOT EXISTS tasks (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  status VARCHAR(20),
  priority INT,
  created_at DATE,
  due_date DATE
)