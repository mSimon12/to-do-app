export const STATUS = {
  BACKLOG: 'backlog',
  OPEN: 'open',
  DONE: 'done'
} as const;

export const STATUS_LABELS= {
  [STATUS.BACKLOG]: 'To Do',
  [STATUS.OPEN]: 'In Progress',
  [STATUS.DONE]: 'Done'
} as const;

export const STATUS_ICONS = {
  [STATUS.DONE]: 'circle-check',
  [STATUS.OPEN]: 'circle-play',
  [STATUS.BACKLOG]: 'list-todo'
} as const;

export const PRIORITY = {
  LOW: 1,
  MEDIUM: 2,
  HIGH: 3
} as const;

export const PRIORITY_LABELS = {
  [PRIORITY.LOW]: 'Low',
  [PRIORITY.MEDIUM]: 'Medium',
  [PRIORITY.HIGH]: 'High',
  'undefined': 'Undefined'
} as const;
