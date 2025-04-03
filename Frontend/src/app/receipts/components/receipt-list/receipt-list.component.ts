// receipt-list.component.ts
import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { RouterModule } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatToolbarModule } from '@angular/material/toolbar';
import { ReceiptService } from '../../services/receipt.service';
import { Receipt } from '../../models/receipt';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
  selector: 'app-receipt-list',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatButtonModule,
    RouterModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatToolbarModule
  ],
  templateUrl: './receipt-list.component.html',
  styleUrls: ['./receipt-list.component.css']
})
export class ReceiptListComponent implements OnInit {
  receipts: Receipt[] = [];
  loading = true;
  error = '';

  constructor(
    private receiptService: ReceiptService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.loadReceipts();
  }

  loadReceipts(): void {
    this.loading = true;
    this.error = '';

    this.receiptService.getAllReceipts().subscribe({
      next: (receipts) => {
        this.receipts = receipts;
        this.loading = false;
      },
      error: (err) => {
        this.loading = false;
        this.error = 'Failed to load receipts';
        this.snackBar.open(this.error, 'Close', {
          duration: 5000,
          panelClass: ['error-snackbar']
        });
      }
    });
  }

  deleteReceipt(id: string, event: Event): void {
    event.stopPropagation();
    if (confirm('Are you sure you want to delete this receipt?')) {
      this.receiptService.deleteReceipt(id).subscribe({
        next: () => {
          this.snackBar.open('Receipt deleted', 'Close', { duration: 3000 });
          this.loadReceipts();
        },
        error: () => {
          this.snackBar.open('Failed to delete receipt', 'Close', { duration: 3000 });
        }
      });
    }
  }

  getTotalItems(receipt: Receipt): number {
    return receipt.items?.length || 0;
  }
}