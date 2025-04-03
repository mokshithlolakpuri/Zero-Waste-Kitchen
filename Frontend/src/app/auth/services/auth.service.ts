import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable, tap, catchError, of } from 'rxjs';
import { map } from 'rxjs/operators';
import { Router } from '@angular/router';
import { environment } from '../../../environments/environment';
import { ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';

interface AuthResponse {
  token: string;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private tokenSubject = new BehaviorSubject<string | null>(null);
  private isAdminSubject = new BehaviorSubject<boolean>(false); // Store admin status

  constructor(private http: HttpClient, private router: Router) {
    const token = localStorage.getItem('token');
    this.tokenSubject.next(token);
  
    // Check admin status on initialization if token exists
    if (token) {
      this.checkAdminStatus().subscribe();
    }
  }

  register(name: string, email: string, password: string): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${environment.apiUrl}/auth/register`, {
      name, email, password
    }).pipe(
      tap(response => this.storeToken(response.token))
    );
  }

  login(email: string, password: string): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${environment.apiUrl}/auth/login`, {
      email, password
    }).pipe(
      tap(response => {
        this.storeToken(response.token);
        this.checkAdminStatus().subscribe(); // Check admin status on login
      })
    );
  }

  logout(): void {
    localStorage.removeItem('token');
    this.tokenSubject.next(null);
    this.isAdminSubject.next(false); // Reset admin status on logout
    this.router.navigate(['/auth/login']);
  }

  getToken(): string | null {
    return this.tokenSubject.value || localStorage.getItem('token');
  }

  isAuthenticated(): boolean {
    return !!this.getToken();
  }

  isAdmin$(): Observable<boolean> {
    return this.isAdminSubject.asObservable(); // Expose admin status as an observable
  }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): boolean {
    if (this.isAuthenticated()) {
      return true;
    } else {
      this.router.navigate(['/auth/login'], { queryParams: { returnUrl: state.url } });
      return false;
    }
  }

  checkAdminStatus$(): Observable<boolean> {
    const token = this.getToken();
    if (!token) {
      this.isAdminSubject.next(false);
      return of(false);
    }

    return this.http.get<{ is_admin: boolean }>(`${environment.apiUrl}/auth/isadmin`).pipe(
      tap(response => this.isAdminSubject.next(response.is_admin)),
      catchError(() => {
        this.isAdminSubject.next(false);
        return of(false);
      }),
      map(response => typeof response === 'object' && response.is_admin) // Transform the response to a boolean
    );
  }

  // auth.service.ts

  checkAdminStatus(): Observable<boolean> {
    const token = this.getToken();
    if (!token) {
      this.isAdminSubject.next(false);
      return of(false);
    }

    return this.http.get<{ is_admin: boolean }>(`${environment.apiUrl}/auth/isadmin`).pipe(
      tap(response => this.isAdminSubject.next(response.is_admin)),
      catchError(() => {
        this.isAdminSubject.next(false);
        return of(false);
      }),
      map(response => typeof response === 'object' && response.is_admin)
    );
  }

  // private checkAdminStatus(): void {
  //   const token = this.getToken();
  //   if (!token) {
  //     this.isAdminSubject.next(false);
  //     return;
  //   }

  //   this.http.get<{ is_admin: boolean }>(`${environment.apiUrl}/auth/isadmin`).pipe(
  //     tap(response => this.isAdminSubject.next(response.is_admin)),
  //     catchError(() => {
  //       this.isAdminSubject.next(false);
  //       return of(false);
  //     })
  //   ).subscribe();
  // }

  private storeToken(token: string): void {
    localStorage.setItem('token', token);
    this.tokenSubject.next(token);
  }
}