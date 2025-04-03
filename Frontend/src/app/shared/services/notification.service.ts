import { Injectable, inject } from '@angular/core';
import { initializeApp } from 'firebase/app';
import { getMessaging, getToken, onMessage, isSupported, type Messaging, type MessagePayload } from 'firebase/messaging';
import { environment } from '../../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { AuthService } from '../../auth/services/auth.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { BehaviorSubject, Observable, from, of } from 'rxjs';
import { catchError, switchMap } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class NotificationService {
  private http = inject(HttpClient);
  private authService = inject(AuthService);
  private snackBar = inject(MatSnackBar);
  private messaging: Messaging | null = null;
  private isMessagingInitialized = new BehaviorSubject<boolean>(false);
  private notificationsEnabled = new BehaviorSubject<boolean>(false);
  private notificationListenerActive = false;

  constructor() {
    this.initializeMessaging();
  }

  private async initializeMessaging(): Promise<void> {
    if (this.messaging || !(await isSupported())) return;

    try {
      const app = initializeApp(environment.firebase);
      this.messaging = getMessaging(app);

      // Check if the service worker is already registered
      const registrations = await navigator.serviceWorker.getRegistrations();
      if (registrations.some((reg) => reg.active?.scriptURL.includes('firebase-messaging-sw.js'))) {
        console.log('Service worker already registered');
      } else {
        const registration = await navigator.serviceWorker.register('/firebase-messaging-sw.js');
        console.log('Service worker registered:', registration);
      }

      this.isMessagingInitialized.next(true);
    } catch (error) {
      console.error('Firebase initialization failed:', error);
    }
  }

  requestPermission(): Observable<string | null> {
    return this.isMessagingInitialized.pipe(
      switchMap(initialized => {
        if (!initialized || !this.messaging) {
          return of(null);
        }
        return from(this.handlePermissionFlow());
      }),
      catchError(err => {
        console.error('Error in requestPermission:', err);
        return of(null);
      })
    );
  }
  
  private async handlePermissionFlow(): Promise<string | null> {
    try {
      const permission = await Notification.requestPermission();
      if (permission === 'granted') {
        const token = await this.getFCMToken();
        this.notificationsEnabled.next(true);
        this.showNotification('Notifications enabled!', 'You will now receive updates');
        return token;
      } else {
        this.notificationsEnabled.next(false);
        this.showNotification('Notifications blocked', 'Enable them in browser settings', 'error');
        return null;
      }
    } catch (error) {
      console.error('Error requesting notification permission:', error);
      throw error;
    }
  }

  listenForNotifications(callback: (payload: MessagePayload) => void): () => void {
    if (!this.messaging) {
      console.warn('Messaging not initialized');
      return () => {};
    }

    if (this.notificationListenerActive) {
      console.log('Notification listener already active');
      return () => {};
    }

    this.notificationListenerActive = true;
    return onMessage(this.messaging, callback);
  }

  private async getFCMToken(): Promise<string> {
    if (!this.messaging) throw new Error('Messaging not initialized');
  
    const currentToken = await getToken(this.messaging, {
      vapidKey: environment.firebase.vapidKey
    });
  
    if (!currentToken) {
      throw new Error('No registration token available');
    }
  
    await this.registerTokenWithBackend(currentToken);
    return currentToken;
  }

  private async registerTokenWithBackend(token: string): Promise<void> {
    // const storedToken = localStorage.getItem('fcm_token');
    // if (storedToken === token) {
    //   console.log('Token already registered');
    //   return;
    // }

    const authToken = this.authService.getToken();
    if (!authToken) throw new Error('User not authenticated');

    return this.http.post(`${environment.apiUrl}/user/fcm-token`, { token }, {
      headers: { Authorization: `Bearer ${authToken}` }
    }).toPromise()
      .then(() => localStorage.setItem('fcm_token', token))
      .catch(error => {
        console.error('Failed to register FCM token:', error);
        throw error;
      });
  }

  toggleNotifications(enable: boolean): void {
    this.notificationsEnabled.next(enable);
    if (!enable) {
      localStorage.removeItem('fcm_token');
      this.showNotification('Notifications disabled', 'You will no longer receive updates', 'info');
    } else {
      this.requestPermission().subscribe({
        error: (err) => console.error('Failed to enable notifications:', err)
      });
    }
  }

  getNotificationStatus(): Observable<boolean> {
    return this.notificationsEnabled.asObservable();
  }

  private showNotification(message: string, action: string, panelClass: string = 'success'): void {
    this.snackBar.open(message, action, {
      duration: 5000,
      panelClass: [`notification-${panelClass}`]
    });
  }
}