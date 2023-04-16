import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AppComponent } from './app.component';
import { PackageResultsComponent } from './package-results/package-results.component';
import { AdvancedSearchComponent } from './advanced-search/advanced-search.component';
import { HomePageComponent } from './home-page/home-page.component';
import { PackagePageComponent } from './package-page/package-page.component';

const routes: Routes = [
  { path: '', component: HomePageComponent},
  { path: 'search', component: PackageResultsComponent},
  { path: 'search/advanced', component: AdvancedSearchComponent},
  { path: 'package', component: PackagePageComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
