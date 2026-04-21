import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Task } from '../../models/task';
import { PRIORITY, PRIORITY_LABELS, STATUS_ICONS } from '../../models/constants';
import { LucideAngularModule } from 'lucide-angular';

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
    if (this.task.priority <= PRIORITY.LOW) return 'p-low';
    if (this.task.priority === PRIORITY.MEDIUM) return 'p-med';
    return 'p-high';
  }

  getPriorityLabel() {
    if (this.task.priority == null) return PRIORITY_LABELS['undefined'];
    if (this.task.priority <= PRIORITY.LOW) return PRIORITY_LABELS[PRIORITY.LOW];
    if (this.task.priority === PRIORITY.MEDIUM) return PRIORITY_LABELS[PRIORITY.MEDIUM];
    return PRIORITY_LABELS[PRIORITY.HIGH];
  }

  getStatusIcon() {
    const status = this.task.status?.toLowerCase() || '';
    return STATUS_ICONS[status as keyof typeof STATUS_ICONS] || 'list-todo';
  }
}
