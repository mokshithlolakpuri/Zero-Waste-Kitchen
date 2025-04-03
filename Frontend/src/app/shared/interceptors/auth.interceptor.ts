// shared/interceptors/auth.interceptor.ts
import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { AuthService } from '../../auth/services/auth.service';

export const AuthInterceptor: HttpInterceptorFn = (req, next) => {
  console.log('[AuthInterceptor] Intercepting request:', req.url);
  
  const authService = inject(AuthService);
  const authToken = authService.getToken(); // Or localStorage.getItem('auth_token')

  if (authToken) {
    console.log('[AuthInterceptor] Adding auth token to request');
    const authReq = req.clone({
      setHeaders: {
        Authorization: `Bearer ${authToken}`
      }
    });
    return next(authReq);
  }

  return next(req);
};