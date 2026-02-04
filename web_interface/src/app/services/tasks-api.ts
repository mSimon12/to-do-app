import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable} from 'rxjs';
import { map } from 'rxjs/operators';
import { Task } from '../models/task';

interface ApiTask {
  id: number;
  title: string;
  status: string;
  priority: number;
  description: string;
  created_at: number;
  due_date: number;
}

interface ApiResponse {
  data: ApiTask[];
}

@Injectable({
  providedIn: 'root',
})

export class TasksApi {
  private apiUrl = 'http://localhost:8080/api/tasks'; // Default URL for the Go API

  constructor(private http: HttpClient) {}

  getTasks(): Observable<Task[]> {
    return this.http.get<ApiResponse>(this.apiUrl).pipe(
      map(res => {
        // Ensure we have a valid array before processing
        const list = res && Array.isArray(res.data) ? res.data : [];
        if (list.length === 0) {
          return [];
        }

        return list.map(t => ({
          id: t.id,
          title: t.title,
          status: t.status,
          priority: t.priority,
          description: t.description,
          created_at: t.created_at ? new Date(t.created_at * 1000).toISOString() : undefined,
          due_date: t.due_date ? new Date(t.due_date * 1000).toISOString() : undefined,
        }));
      })
    );
  }

  getTask(id: number): Observable<Task> {
    return this.http.get<Task>(`${this.apiUrl}/${id}`);
  }

  createTask(task: Task): Observable<Task> {
    return this.http.post<Task>(this.apiUrl, task);
  }

  updateTask(id: number, task: Partial<Task>): Observable<Task> {
    return this.http.put<Task>(`${this.apiUrl}/${id}`, task);
  }

  deleteTask(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }
}
