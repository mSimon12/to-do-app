import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DragDropModule, CdkDragDrop, moveItemInArray, transferArrayItem } from '@angular/cdk/drag-drop';
import { TasksApi } from '../../services/tasks-api';
import { Task } from '../../models/task';
import { TaskCard } from '../task-card/task-card';
import { TaskDetails } from '../task-details/task-details';
import { LucideAngularModule, Plus } from 'lucide-angular';

interface Column {
  name: string;
  status: string;
  tasks: Task[];
}

@Component({
  selector: 'app-main-board',
  standalone: true,
  imports: [CommonModule, DragDropModule, TaskCard, TaskDetails, LucideAngularModule],
  templateUrl: './main-board.html',
  styleUrl: './main-board.css',
})

export class MainBoard implements OnInit {
  readonly PlusIcon = Plus;
  
  columns: Column[] = [
    { name: 'To Do', status: 'To Do', tasks: [] },
    { name: 'In Progress', status: 'In Progress', tasks: [] },
    { name: 'Done', status: 'done', tasks: [] }
  ];
  selectedTask: Task | null = null;

  constructor(private taskService: TasksApi) {}

  ngOnInit() {
    this.loadTasks();
  }

  loadTasks() {
    this.taskService.getTasks().subscribe({
      next: (tasks) => {
        console.log('Loaded tasks:', tasks);
        this.columns.forEach(col => {
          col.tasks = tasks.filter(t => t.status === col.status);
        });
      },
      error: (err) => {
        console.error('Error loading tasks:', err);
        // Fallback for demo if API is not running
        const mockTasks: Task[] = [
          { id: 1, title: 'Sample Task 1', description: 'Description 1', status: 'To Do', priority: 1 },
          { id: 2, title: 'Sample Task 2', description: 'Description 2', status: 'In Progress', priority: 2 },
          { id: 3, title: 'Sample Task 3', description: 'Description 3', status: 'Done', priority: 3 }
        ];
        this.columns.forEach(col => {
          col.tasks = mockTasks.filter(t => t.status === col.status);
        });
      }
    });
  }

  drop(event: CdkDragDrop<Task[]>, newStatus: string) {
    if (event.previousContainer === event.container) {
      moveItemInArray(event.container.data, event.previousIndex, event.currentIndex);
    } else {
      const task = event.previousContainer.data[event.previousIndex];
      transferArrayItem(
        event.previousContainer.data,
        event.container.data,
        event.previousIndex,
        event.currentIndex
      );
      
      if (task.id) {
        this.taskService.updateTask(task.id, { ...task, status: newStatus }).subscribe();
      }
    }
  }

  selectTask(task: Task) {
    this.selectedTask = { ...task };
  }

  openCreateModal() {
    this.selectedTask = {
      title: '',
      description: '',
      status: 'To Do',
      priority: 1
    };
  }

  onTaskSaved() {
    this.selectedTask = null;
    this.loadTasks();
  }
}