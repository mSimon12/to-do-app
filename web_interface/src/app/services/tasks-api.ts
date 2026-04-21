import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, expand, reduce, EMPTY } from 'rxjs';
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
  pagination: { offset: number; limit: number; total_tasks: number };
}

@Injectable({
  providedIn: 'root',
})

export class TasksApi {
  private apiUrl = 'http://localhost:8080/api/tasks'; // Default URL for the Go API

  constructor(private http: HttpClient) {}

  private toIsoNoMs(timestamp: number): string {
    return new Date(timestamp * 1000).toISOString().split('T')[0];
  }

  private toTask(t: ApiTask): Task {
    return {
      id: t.id,
      title: t.title,
      status: t.status,
      priority: t.priority,
      description: t.description,
      created_at: t.created_at ? this.toIsoNoMs(t.created_at) : undefined,
      due_date: t.due_date ? this.toIsoNoMs(t.due_date) : undefined,
    };
  }

  private fetchPage(offset: number): Observable<ApiResponse> {
    return this.http.get<ApiResponse>(this.apiUrl, {
      params: { limit: '100', offset: String(offset) }
    });
  }

  getTasks(): Observable<Task[]> {
    return this.fetchPage(0).pipe(
      expand(res => {
        const { offset, limit, total_tasks } = res.pagination;
        return offset + limit < total_tasks ? this.fetchPage(offset + limit) : EMPTY;
      }),
      reduce((all: Task[], res) => {
        const page = Array.isArray(res.data) ? res.data.map(t => this.toTask(t)) : [];
        return all.concat(page);
      }, [])
    );
  }

  getTask(id: number): Observable<Task> {
    interface SingleResponse { task: ApiTask }
    return this.http.get<SingleResponse>(`${this.apiUrl}/${id}`).pipe(
      map(res => {
        const t = res && res.task ? res.task : null;
        if (!t) {
          return {} as Task;
        }

        return {
          id: t.id,
          title: t.title,
          status: t.status,
          priority: t.priority,
          description: t.description,
          created_at: t.created_at ? this.toIsoNoMs(t.created_at) : undefined,
          due_date: t.due_date ? this.toIsoNoMs(t.due_date) : undefined,
        } as Task;
      })
    );
  }

  createTask(task: Task): Observable<Task> {
    const payload = {
      ...task,
      due_date: task.due_date ? Math.floor(new Date(task.due_date).getTime() / 1000) : undefined,
    };
    return this.http.post<Task>(this.apiUrl, payload);
  }

  updateTask(id: number, task: Partial<Task>): Observable<Task> {
    const payload = {
      ...task,
      due_date: task.due_date ? Math.floor(new Date(task.due_date).getTime() / 1000) : undefined,
    };
    return this.http.put<Task>(`${this.apiUrl}/${id}`, payload);
  }

  deleteTask(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }
}
