import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PackageResultsComponent } from './package-results.component';

describe('PackageResultsComponent', () => {
  let component: PackageResultsComponent;
  let fixture: ComponentFixture<PackageResultsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ PackageResultsComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(PackageResultsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
