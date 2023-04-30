import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { AuthenticationRequest, DefaultService, User } from 'generated';
import { LoginService } from 'loginService/login.service';
import { BASE_PATH, AuthenticationToken } from 'generated';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
  selector: 'app-login-page',
  templateUrl: './login-page.component.html',
  styleUrls: ['./login-page.component.scss']
})
export class LoginPageComponent {
  username: string = ""
  password: string = ""
  authRequest: AuthenticationRequest
  constructor(private httpClient: HttpClient, private loginService: LoginService, private router: Router, private _snackBar: MatSnackBar) { }

  login() {
    this.authRequest = {
      User: {
        name: this.username,
        isAdmin: true
      },
      Secret: {
        password: this.password
      }
    }
    var observe = "body"
    this.httpClient.request('put',
      `https://ece461-server-rruekicr4q-uc.a.run.app/authenticate`,
      {
        body: this.authRequest,
        responseType: 'text'
      }).subscribe(body => {
        this.loginService.setToken(body)
        this.router.navigate(['/home'])
      }, error => {
        this._snackBar.open("Invalid Credentials", "Ok")
      });
  }
}
