// ocr-preview-dialog.component.ts
import { Component, Inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatListModule } from '@angular/material/list';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatSelectModule } from '@angular/material/select';
import { MatOptionModule } from '@angular/material/core';
import { ReceiptService } from '../../services/receipt.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner'; 

@Component({
  selector: 'app-ocr-preview-dialog',
  standalone: true,
  imports: [
    CommonModule,
    MatDialogModule,
    MatButtonModule,
    MatCardModule,
    MatIconModule,
    FormsModule,
    ReactiveFormsModule,
    MatInputModule,
    MatFormFieldModule,
    MatListModule,
    MatDatepickerModule,
    MatNativeDateModule,
    MatSelectModule,
    MatOptionModule,
    MatProgressSpinnerModule 
  ],
  templateUrl: './ocr-preview-dialog.component.html',
  styleUrls: ['./ocr-preview-dialog.component.css']
})
export class OcrPreviewDialogComponent {
  editedItems: any[];
  showRawText = false;
  showImage = false;
  confidence: number;
  isUploading = false;

  constructor(
    public dialogRef: MatDialogRef<OcrPreviewDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { 
      image: string | File, 
      ocrResult: {
        storeName: string;
        purchaseDate: Date;
        totalAmount: number;
        items: any[];
        rawText: string;
        confidence: number;
      } 
    },
    private receiptService: ReceiptService,
    private snackBar: MatSnackBar
  ) {
    // Ensure purchaseDate is a Date object
    this.data.ocrResult.purchaseDate = data.ocrResult.purchaseDate ? new Date(data.ocrResult.purchaseDate) : new Date();
    
    this.editedItems = data.ocrResult.items.map(item => ({
      ...item,
      expiryDate: null, // Start with empty expiry date
      storageLocation: item.storageLocation || 'dry_pantry'
    }));
    
    this.confidence = data.ocrResult.confidence;
  }

  toggleImage(): void {
    this.showImage = !this.showImage; // Toggle the visibility of the receipt image
  }

  isFormValid(): boolean {
    return typeof this.data.ocrResult.storeName === 'string' && this.data.ocrResult.storeName.trim() !== '' &&
           this.data.ocrResult.purchaseDate &&
           this.editedItems.every(item => 
             item.name && 
             item.quantity && 
             item.unit && 
             item.expiryDate && 
             item.storageLocation
           );
  }

  async onSave(): Promise<void> {
    if (!this.isFormValid()) {
      this.snackBar.open('Please fill all required fields', 'Close', { duration: 3000 });
      return;
    }

    this.isUploading = true;
    try {
      const receiptData = {
        storeName: this.data.ocrResult.storeName,
        purchaseDate: this.data.ocrResult.purchaseDate.toISOString().split('T')[0],
        totalAmount: this.data.ocrResult.totalAmount,
        items: this.editedItems.map(item => ({
          ...item,
          expiryDate: item.expiryDate.toISOString()
        }))
      };

      const formData = new FormData();
      formData.append('receipt', JSON.stringify(receiptData));
      
      if (this.data.image instanceof File) {
        formData.append('file', this.data.image);
      } else if (typeof this.data.image === 'string') {
        const blob = await fetch(this.data.image).then(r => r.blob());
        const file = new File([blob], 'receipt.jpg', { type: 'image/jpeg' });
        formData.append('file', file);
      }

      await this.receiptService.uploadReceipt(formData).toPromise();
      this.snackBar.open('Receipt saved successfully', 'Close', { duration: 3000 });
      this.dialogRef.close(true);
    } catch (error) {
      this.snackBar.open('Error saving receipt', 'Close', { duration: 3000 });
      console.error('Upload error:', error);
    } finally {
      this.isUploading = false;
    }
  }

  addItem(): void {
    this.editedItems.push({
      name: '',
      quantity: 1,
      unit: 'pcs',
      expiryDate: null,
      storageLocation: 'dry_pantry'
    });
  }

  removeItem(index: number): void {
    this.editedItems.splice(index, 1);
  }

  onCancel(): void {
    this.dialogRef.close();
  }

  getConfidenceColor(): string {
    if (this.confidence > 80) return 'green';
    if (this.confidence > 60) return 'orange';
    return 'red';
  }
}