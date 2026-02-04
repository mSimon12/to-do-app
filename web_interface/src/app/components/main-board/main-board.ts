import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DragDropModule, CdkDragDrop, moveItemInArray, transferArrayItem } from '@angular/cdk/drag-drop';
import { TasksApi } from '../../services/tasks-api';
import { Theme } from '../../services/theme';
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
  public board_name: string = 'TaskFlow';
  selectedTask: Task | null = null;

   columns: Column[] = [
    { name: 'To Do', status: 'To Do', tasks: [] },
    { name: 'In Progress', status: 'In Progress', tasks: [] },
    { name: 'Done', status: 'Done', tasks: [] }
  ];

  constructor(
    private taskService: TasksApi,
    public themeService: Theme,
    private cd: ChangeDetectorRef
  ) {}

  ngOnInit() {
    this.loadTasks();
  }

  loadTasks() {
    this.taskService.getTasks().subscribe({
      next: (tasks) => {
        this.columns.forEach(col => {
          col.tasks = tasks.filter(t => (t.status || '').toLowerCase() === col.status.toLowerCase());
        });
        try { this.cd.detectChanges(); } catch (e) { }
      },
      error: (err) => {
        console.error('Error loading tasks:', err);
        const mockTasks: Task[] = [
          { id: 1, title: 'Refactor Auth Service', description: 'Clean up the login logic and implement better error handling.', status: 'To Do', priority: 3 },
          { id: 2, title: 'Update Dashboard UI', description: 'Apply the new design system to the main dashboard views.', status: 'In Progress', priority: 2 },
          { id: 3, title: 'Fix Header Bug', description: 'The mobile menu is not closing correctly on navigation.', status: 'Done', priority: 1 }
        ];
        this.columns.forEach(col => {
          col.tasks = mockTasks.filter(t => (t.status || '').toLowerCase() === col.status.toLowerCase());
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