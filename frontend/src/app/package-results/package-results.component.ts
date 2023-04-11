import { Component, Input, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { DefaultService, PackageMetadata } from 'generated';
import { PackageQuery } from 'generated';
@Component({
  selector: 'app-package-results',
  templateUrl: './package-results.component.html',
  styleUrls: ['./package-results.component.scss']
})
export class PackageResultsComponent implements OnInit {
  name?: string;
  version?: string;
  regex?: string;
  packages!: PackageMetadata[];

  constructor (private route: ActivatedRoute, private service: DefaultService) {}

  ngOnInit(): void {
    this.route.queryParams.subscribe( params=> {
      this.name = params['name'];
      this.version = params['version'];
      console.log("Searching for ", this.name)
      this.searchByNameVersion()
    })
  }
  searchByNameVersion() {
    var query:PackageQuery
    if (this.name && this.version) {
      query = {Name: this.name, Version: this.version}
    } else if (this.name) {
      query = {Name: this.name}
    } else {
      query = {Name: "*"}
    }
 
    var queries = [query]

    this.service.packagesList(queries, "").subscribe(body => {
      this.packages = body;
      console.log(this.packages);
    })
    // TODO: Deal with pagination
  }

  //TODO: LATER
  searchByRegex(regex: string) {
  }
}

