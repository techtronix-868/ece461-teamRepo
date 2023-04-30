import { Component } from '@angular/core';
import { AuthenticationRequest, DefaultService, User } from 'generated';

@Component({
  selector: 'app-login-page',
  templateUrl: './login-page.component.html',
  styleUrls: ['./login-page.component.scss']
})
export class LoginPageComponent {
  username: string = ""
  password: string = ""
  authRequest: AuthenticationRequest
  constructor(private service: DefaultService) {}

  login() {
    this.authRequest.User = {
      name: this.username,
      isAdmin: true
    }
    this.authRequest.Secret = {
      password: this.password
    }
    this.service.createAuthToken(this.authRequest).subscribe(body => {
      console.log(body)
    })
  }
}
