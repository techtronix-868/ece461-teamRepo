import { Component, OnInit } from '@angular/core';
import { DefaultService, AuthenticationRequest, User, UserAuthenticationInfo} from 'generated';
@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
  providers: [DefaultService]
})
export class AppComponent implements OnInit {
  title = 'frontend';

  constructor(private service: DefaultService) {

  }

  ngOnInit(): void {
    var user:User = {name: "mabaums", isAdmin: true};
    var authInfo:UserAuthenticationInfo = {password: "mabaums"};
    var request:AuthenticationRequest = {user: user, secret: authInfo};

    this.service.createAuthToken(request).subscribe(body=> {
      console.log(body);
    });
  }
}
