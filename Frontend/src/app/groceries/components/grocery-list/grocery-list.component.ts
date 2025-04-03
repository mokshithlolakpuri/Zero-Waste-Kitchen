import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginatorModule, PageEvent } from '@angular/material/paginator';
import { MatSelectModule } from '@angular/material/select';
import { MatInputModule } from '@angular/material/input';
import { MatCardModule } from '@angular/material/card';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatChipsModule } from '@angular/material/chips';
import { MatMenuModule } from '@angular/material/menu';
import { MatFormFieldModule } from '@angular/material/form-field';

import { GroceryService } from '../../services/grocery.service';
import { GroceryItem } from '../../models/grocery-item';
import { ConfirmationDialogComponent } from '../../../shared/components/confirmation-dialog/confirmation-dialog.component';
import { GroceryEditDialogComponent } from '../grocery-edit-dialog/grocery-edit-dialog.component';

@Component({
  selector: 'app-grocery-list',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    FormsModule,
    MatButtonModule,
    MatIconModule,
    MatDialogModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatPaginatorModule,
    MatSelectModule,
    MatInputModule,
    MatCardModule,
    MatTooltipModule,
    MatChipsModule,
    MatMenuModule,
    MatFormFieldModule 
  ],
  templateUrl: './grocery-list.component.html',
  styleUrls: ['./grocery-list.component.css']
})
export class GroceryListComponent implements OnInit {
  groceries: GroceryItem[] = [];
  filteredGroceries: GroceryItem[] = [];
  loading = true;
  error: string | null = null;
  
  // Filter properties
  searchTerm = '';
  expiryFilter: string | null = null;
  storageLocationFilter: string | null = null;
  storageLocations: string[] = [];
  
  // Pagination properties
  pageSize = 8;
  pageIndex = 0;
  pageSizeOptions = [4, 8, 12, 24];

  constructor(
    private groceryService: GroceryService,
    private dialog: MatDialog,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.loadGroceries();
  }

  loadGroceries(): void {
    this.loading = true;
    this.error = null;

    this.groceryService.getAllGroceries().subscribe({
      next: (groceries) => {
        this.groceries = groceries;
        this.storageLocations = [...new Set(groceries.map(g => g.storageLocation))];
        this.applyFilters();
        this.loading = false;
      },
      error: (err) => {
        this.loading = false;
        this.error = 'Failed to load groceries. Please try again later.';
        console.error('Failed to load groceries:', err);
      }
    });
  }

  applyFilters(): void {
    let filtered = [...this.groceries];
    
    // Apply search filter
    if (this.searchTerm.trim()) {
      const term = this.searchTerm.toLowerCase().trim();
      filtered = filtered.filter(grocery => 
        grocery.name.toLowerCase().includes(term)
      );
    }
    
    // Apply expiry filter
    if (this.expiryFilter) {
      const days = parseInt(this.expiryFilter);
      filtered = filtered.filter(grocery => {
        if (!grocery.expiry_date) return false;
        const expiryDate = new Date(grocery.expiry_date);
        const today = new Date();
        const diffTime = expiryDate.getTime() - today.getTime();
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
        return diffDays <= days && diffDays >= 0;
      });
    }
    
    // Apply storage location filter
    if (this.storageLocationFilter) {
      filtered = filtered.filter(grocery => 
        grocery.storageLocation === this.storageLocationFilter
      );
    }

    this.filteredGroceries = filtered;
    this.pageIndex = 0; // Reset to first page when filters change
  }

  clearFilters(): void {
    this.searchTerm = '';
    this.expiryFilter = null;
    this.storageLocationFilter = null;
    this.applyFilters();
  }

  getPaginatedItems(): GroceryItem[] {
    const startIndex = this.pageIndex * this.pageSize;
    return this.filteredGroceries.slice(startIndex, startIndex + this.pageSize);
  }

  onPageChange(event: PageEvent): void {
    this.pageIndex = event.pageIndex;
    this.pageSize = event.pageSize;
  }

  openAddDialog(): void {
    const dialogRef = this.dialog.open(GroceryEditDialogComponent, {
      width: '400px',
      data: { grocery: null }
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        this.loadGroceries();
        this.snackBar.open('Grocery item added successfully!', 'Close', {
          duration: 3000,
          panelClass: ['success-snackbar']
        });
      }
    });
  }

  openEditDialog(grocery: GroceryItem): void {
    const dialogRef = this.dialog.open(GroceryEditDialogComponent, {
      width: '400px',
      data: { grocery }
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        this.loadGroceries();
        this.snackBar.open('Grocery item updated successfully!', 'Close', {
          duration: 3000,
          panelClass: ['success-snackbar']
        });
      }
    });
  }

  deleteGrocery(id: number): void {
    const dialogRef = this.dialog.open(ConfirmationDialogComponent, {
      width: '300px',
      data: {
        title: 'Delete Grocery',
        message: 'Are you sure you want to delete this grocery item? This action cannot be undone.'
      }
    });

    dialogRef.afterClosed().subscribe((confirmed) => {
      if (confirmed) {
        this.groceryService.deleteGrocery(id.toString()).subscribe({
          next: () => {
            this.groceries = this.groceries.filter((grocery) => grocery.id !== id);
            this.applyFilters();
            this.snackBar.open('Grocery item deleted successfully!', 'Close', {
              duration: 3000,
              panelClass: ['success-snackbar']
            });
          },
          error: (err) => {
            console.error('Failed to delete grocery:', err);
            this.snackBar.open('Failed to delete grocery item. Please try again.', 'Close', {
              duration: 3000,
              panelClass: ['error-snackbar']
            });
          }
        });
      }
    });
  }

  isExpired(expiryDate: string | undefined): boolean {
    if (!expiryDate) return false;
    const today = new Date();
    const expiry = new Date(expiryDate);
    return expiry < today;
  }
}