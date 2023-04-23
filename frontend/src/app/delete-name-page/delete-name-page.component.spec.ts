import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DeleteNamePageComponent } from './delete-name-page.component';

describe('DeleteNamePageComponent', () => {
  let component: DeleteNamePageComponent;
  let fixture: ComponentFixture<DeleteNamePageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ DeleteNamePageComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DeleteNamePageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
