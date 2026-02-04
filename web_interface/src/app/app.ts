import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { MainBoard } from './components/main-board/main-board';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, MainBoard],
  templateUrl: './app.html',
  styleUrl: './app.css'
})

export class App {
  protected readonly title = 'To-Do App';
}
