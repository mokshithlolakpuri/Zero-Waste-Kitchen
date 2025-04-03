import { Component, OnInit, OnDestroy } from '@angular/core';
import { NotificationService } from '../../services/notification.service';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { MessagePayload } from 'firebase/messaging';

interface Notification {
  title: string;
  body: string;
  timestamp: Date;
}

@Component({
  selector: 'app-notification',
  templateUrl: './notification.component.html',
  styleUrls: ['./notification.component.css']
})
export class NotificationComponent implements OnInit, OnDestroy {
  notifications: Notification[] = [];
  showNotifications = false;
  notificationsEnabled = false;
  private destroy$ = new Subject<void>();

  constructor(private notificationService: NotificationService) {}

  ngOnInit(): void {
    this.notificationService.getNotificationStatus()
      .pipe(takeUntil(this.destroy$))
      .subscribe(status => this.notificationsEnabled = status);

    this.notificationService.listenForNotifications((payload: MessagePayload) => {
      this.notifications.unshift({
        title: payload.notification?.title || 'Notification',
        body: payload.notification?.body || '',
        timestamp: new Date()
      });
    });

    this.notificationService.requestPermission()
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        error: (err) => console.error('Failed to get notification permission:', err)
      });
  }

  toggleNotifications(): void {
    this.notificationService.toggleNotifications(!this.notificationsEnabled);
  }

  clearNotifications(): void {
    this.notifications = [];
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}