import { Component, OnInit } from '@angular/core';
import { DefaultService, AuthenticationRequest, User, UserAuthenticationInfo, PackageMetadata, PackageQuery} from 'generated';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  providers: [DefaultService]
})
export class AppComponent implements OnInit {
  title = 'frontend';
  packages: PackageMetadata[];
  constructor(private service: DefaultService) {

  }

  ngOnInit(): void {
    var user:User = {name: "mabaums", isAdmin: true};
    var authInfo:UserAuthenticationInfo = {password: "mabaums"};
    var request:AuthenticationRequest = {User: user, Secret: authInfo};

    this.service.createAuthToken(request).subscribe(body=> {
      console.log(body);
    });

    
    var query:PackageQuery = {Name: "*"}
    var queries = [query]

    this.service.packagesList(queries, "").subscribe(body => {
      this.packages = body;
      console.log(this.packages);
    })

    this.service.packageRetrieve("", "1").subscribe(body => {
      console.log(body);
    });
  }
}