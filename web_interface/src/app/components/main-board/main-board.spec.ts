import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MainBoard } from './main-board';

describe('MainBoard', () => {
  let component: MainBoard;
  let fixture: ComponentFixture<MainBoard>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MainBoard]
    })
    .compileComponents();

    fixture = TestBed.createComponent(MainBoard);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
