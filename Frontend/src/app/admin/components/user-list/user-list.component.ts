// user-list.component.ts
import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { UserService } from '../../services/user.service';
import { NotificationService } from '../../services/notification.service';
import { User } from '../../models/user';
import { NotificationDialogComponent } from '../notification-dialog/notification-dialog.component';
import { MatSnackBar } from '@angular/material/snack-bar';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';

@Component({
  selector: 'app-user-list',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatTableModule,
    MatButtonModule,
    MatIconModule,
    MatProgressBarModule,
  ],
  templateUrl: './user-list.component.html',
  styleUrls: ['./user-list.component.css']
})
export class UserListComponent implements OnInit {
  users: User[] = [];
  displayedColumns: string[] = ['is_admin', 'name', 'email', 'actions'];
  isLoading = true;

  constructor(
    private userService: UserService,
    private notificationService: NotificationService,
    private dialog: MatDialog,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.userService.getUsers().subscribe({
      next: (data) => {
        this.users = data;
        this.isLoading = false;
      },
      error: () => {
        this.isLoading = false;
        this.snackBar.open('Failed to load users', 'Close', { duration: 3000 });
      }
    });
  }

  openNotificationDialog(userId: number): void {
    const dialogRef = this.dialog.open(NotificationDialogComponent, {
      width: '500px',
      data: { userId }
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        this.notificationService
          .sendNotification({ userId, message: result.message })
          .subscribe({
            next: () => {
              this.snackBar.open('Notification sent successfully', 'Close', { duration: 3000 });
            },
            error: () => {
              this.snackBar.open('Failed to send notification', 'Close', { duration: 3000 });
            }
          });
      }
    });
  }
}