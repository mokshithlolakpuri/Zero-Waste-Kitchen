import { Injectable } from '@angular/core';
import { CanActivate, CanActivateChild, Router } from '@angular/router';
import { Observable } from 'rxjs';
import { map, tap } from 'rxjs/operators';
import { AuthService } from '../../auth/services/auth.service';
import { MatSnackBar } from '@angular/material/snack-bar';

@Injectable({
  providedIn: 'root'
})
export class AuthGuard implements CanActivate, CanActivateChild {
  constructor(
    private authService: AuthService,
    private router: Router,
    private snackBar: MatSnackBar
  ) {}

  canActivate(): boolean {
    if (this.authService.isAuthenticated()) {
      return true;
    } else {
      this.snackBar.open('You need to log in to access this page.', 'Dismiss', {
        duration: 3000,
        panelClass: ['error-snackbar']
      });
      this.router.navigate(['/auth/login']);
      return false;
    }
  }

  canActivateChild(): Observable<boolean> {
    return this.authService.checkAdminStatus$().pipe(
      tap((isAdmin) => {
        if (!isAdmin) {
          this.snackBar.open('You are not authorized to access this page.', 'Dismiss', {
            duration: 5000,
            panelClass: ['error-snackbar']
          });
          this.router.navigate(['/groceries']);
        }
      }),
      map((isAdmin) => isAdmin)
    );
  }
}
