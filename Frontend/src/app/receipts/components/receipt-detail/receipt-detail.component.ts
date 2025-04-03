// receipt-detail.component.ts
import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatListModule } from '@angular/material/list';
import { MatDividerModule } from '@angular/material/divider';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ReceiptService } from '../../services/receipt.service';
import { Receipt } from '../../models/receipt';

@Component({
  selector: 'app-receipt-detail',
  standalone: true,
  imports: [
    CommonModule,
    MatButtonModule,
    MatCardModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatListModule,
    MatDividerModule
  ],
  templateUrl: './receipt-detail.component.html',
  styleUrls: ['./receipt-detail.component.css']
})
export class ReceiptDetailComponent implements OnInit {
  loading = true;
  error = '';
  receipt: Receipt | null = null;

  constructor(
    private route: ActivatedRoute,
    public  router: Router,
    private receiptService: ReceiptService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (!id) {
      this.router.navigate(['/receipts']);
      return;
    }

    this.loadReceipt(id);
  }

  loadReceipt(id: string): void {
    this.loading = true;
    this.error = '';

    this.receiptService.getReceipt(id).subscribe({
      next: (receipt) => {
        this.receipt = receipt;
        this.loading = false;
      },
      error: (err) => {
        this.loading = false;
        this.error = 'Failed to load receipt details';
        this.snackBar.open(this.error, 'Close', {
          duration: 5000,
          panelClass: ['error-snackbar']
        });
        this.router.navigate(['/receipts']);
      }
    });
  }

  getTotalPrice(item: any): number {
    return item.price ? item.price * item.quantity : 0;
  }
}