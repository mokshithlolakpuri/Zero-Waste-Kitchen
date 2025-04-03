import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { AuthService } from '../../../auth/services/auth.service';
import { NotificationService } from '../../services/notification.service';
import { Observable, Subscription } from 'rxjs';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterLinkActive, MatIconModule],
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit, OnDestroy {
  isMenuOpen = false;
  isAdmin = false; // Store admin status as a boolean
  notificationsEnabled$: Observable<boolean>;
  private adminSubscription!: Subscription;

  constructor(
    public authService: AuthService,
    private notificationService: NotificationService
  ) {
    this.notificationsEnabled$ = this.notificationService.getNotificationStatus();
  }

  ngOnInit(): void {
    // Check admin status immediately if authenticated
    if (this.authService.isAuthenticated()) {
      this.authService.checkAdminStatus().subscribe();
    }
    
    // Then subscribe to changes
    this.adminSubscription = this.authService.isAdmin$().subscribe((isAdmin) => {
      this.isAdmin = isAdmin;
    });
  }

  ngOnDestroy(): void {
    // Unsubscribe to avoid memory leaks
    if (this.adminSubscription) {
      this.adminSubscription.unsubscribe();
    }
  }

  toggleMenu(): void {
    this.isMenuOpen = !this.isMenuOpen;
  }

  toggleNotifications(): void {
    // this.notificationService.toggleNotificationStatus();
    }

  logout(): void {
    this.authService.logout();
    this.isMenuOpen = false;
  }
}