import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Task } from '../../models/task';
import { TasksApi } from '../../services/tasks-api';
import { LucideAngularModule } from 'lucide-angular';

@Component({
  selector: 'app-task-details',
  standalone: true,
  imports: [CommonModule, FormsModule, LucideAngularModule],
  templateUrl: './task-details.html',
  styleUrl: './task-details.css',
})
export class TaskDetails {
  @Input() task!: Task;
  @Output() close = new EventEmitter<void>();
  @Output() saved = new EventEmitter<void>();

  constructor(private taskService: TasksApi) {}

  saveTask() {
    if (!this.task.title) return;
    if (this.task.id) {
      this.taskService.updateTask(this.task.id, this.task).subscribe(() => this.saved.emit());
    } else {
      this.taskService.createTask(this.task).subscribe(() => this.saved.emit());
    }
  }

  deleteTask() {
    if (this.task.id && confirm('Are you sure you want to delete this task?')) {
      this.taskService.deleteTask(this.task.id).subscribe(() => this.saved.emit());
    }
  }
}
