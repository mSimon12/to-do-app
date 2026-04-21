export interface Task {
  id?: number;
  title: string;
  description?: string;
  status?: string; // e.g., 'To Do', 'In Progress', 'Done'
  priority?: number;
  created_at?: string;
  due_date?: string;
}
