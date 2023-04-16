import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { DefaultService, PackageData, PackageHistoryEntry } from 'generated';

@Component({
  selector: 'app-package-page',
  templateUrl: './package-page.component.html',
  styleUrls: ['./package-page.component.scss']
})
export class PackagePageComponent implements OnInit {

  constructor (private route: ActivatedRoute, private service: DefaultService) {}

  name: string = ""
  id: string = ""
  pkg_history: PackageHistoryEntry[]
  ngOnInit(): void {
    this.route.queryParams.subscribe( params => {
      this.name = params['name'];
      this.id = params['id'];
      this.getPackageByName()
    })
  }

  getPackageByName() {
    this.service.packageByNameGet(this.name, "").subscribe(body => {
      this.pkg_history = body;
      console.log("Retrieved pkg ", this.name, " ", this.pkg_history)
    })
  }
}
