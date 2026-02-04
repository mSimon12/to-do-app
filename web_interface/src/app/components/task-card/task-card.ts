import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Task } from '../../models/task';
import { LucideAngularModule } from 'lucide-angular';
import { isNull } from 'util';

@Component({
  selector: 'app-task-card',
  standalone: true,
  imports: [CommonModule, LucideAngularModule],
  templateUrl: './task-card.html',
  styleUrl: './task-card.css',
})
export class TaskCard {
  @Input() task!: Task;

  getPriorityClass() {
    if (this.task.priority == null) return 'undefined';
    if (this.task.priority <= 1) return 'p-low';
    if (this.task.priority === 2) return 'p-med';
    return 'p-high';
  }

  getStatusIcon() {
    switch (this.task.status) {
      case 'Done': return 'check-circle';
      case 'In Progress': return 'play-circle';
      default: return 'list-todo';
    }
  }
}
