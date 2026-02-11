import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DragDropModule, CdkDragDrop, moveItemInArray, transferArrayItem } from '@angular/cdk/drag-drop';
import { TasksApi } from '../../services/tasks-api';
import { Theme } from '../../services/theme';
import { Task } from '../../models/task';
import { STATUS, STATUS_LABELS } from '../../models/constants';
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

   columns: Column[] = Object.values(STATUS).map(statusValue => ({
    name: STATUS_LABELS[statusValue],
    status: statusValue,
    tasks: []
  }));

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
        this.taskService.updateTask(task.id, { ...task, status: newStatus }).subscribe(
          () => {
            this.loadTasks();
          }
        );
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
      status: STATUS.BACKLOG,
      priority: 1
    };
  }

  onTaskSaved() {
    this.selectedTask = null;
    this.loadTasks();
  }
}