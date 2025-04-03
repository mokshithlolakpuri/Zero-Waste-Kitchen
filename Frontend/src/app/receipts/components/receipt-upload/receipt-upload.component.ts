// receipt-upload.component.ts
import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatIconModule } from '@angular/material/icon';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { OcrService } from '../../services/ocr.service';
import { OcrPreviewDialogComponent } from '../ocr-preview-dialog/ocr-preview-dialog.component';

@Component({
  selector: 'app-receipt-upload',
  standalone: true,
  imports: [
    CommonModule,
    MatButtonModule,
    MatCardModule,
    MatProgressBarModule,
    MatIconModule
  ],
  templateUrl: './receipt-upload.component.html',
  styleUrls: ['./receipt-upload.component.css']
})
export class ReceiptUploadComponent {
  selectedFile: File | null = null;
  processing = false;
  previewImage: string | ArrayBuffer | null = null;
  maxFileSizeMB = 5;
  progress = 0;
  status = '';

  constructor(
    private ocrService: OcrService,
    private router: Router,
    private dialog: MatDialog
  ) {}

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const file = input.files[0];
      
      if (file.size > this.maxFileSizeMB * 1024 * 1024) {
        this.status = `File size exceeds ${this.maxFileSizeMB}MB limit`;
        return;
      }

      this.selectedFile = file;
      
      // Create preview
      const reader = new FileReader();
      reader.onload = () => {
        this.previewImage = reader.result;
      };
      reader.readAsDataURL(file);
    }
  }

  

  async processReceipt(): Promise<void> {
    if (!this.selectedFile) return;

    this.processing = true;
    this.status = 'Starting OCR processing...';
    
    try {
      const result = await this.ocrService.recognize(this.selectedFile);
      
      // Parse the OCR result into structured data
      const parsedData = this.parseOcrResult(result.text);
      
      this.openOcrPreviewDialog({
        image: this.previewImage,
        ocrResult: {
          ...parsedData,
          rawText: result.text,
          confidence: result.confidence
        }
      });
    } catch (error) {
      this.status = 'OCR processing failed';
      console.error('OCR Error:', error);
    } finally {
      this.processing = false;
    }
  }

  private parseOcrResult(text: string): any {
    // Implement your receipt parsing logic here
    // This is a simple example - you'll need to customize for your receipt format
    const lines = text.split('\n').filter(line => line.trim());
    
    // Simple pattern matching - adjust based on your receipts
    const storeMatch = text.match(/(Walmart|Target|Costco|Amazon|Whole Foods)/i);
    const dateMatch = text.match(/(\d{1,2}[\/\-]\d{1,2}[\/\-]\d{2,4})/);
    const totalMatch = text.match(/total.*?(\d+\.\d{2})/i);
    
    const items = lines
      .filter(line => {
        // Simple item detection - adjust as needed
        const priceMatch = line.match(/(\d+\.\d{2})$/);
        return priceMatch && line.length < 50 && !line.match(/total|subtotal|tax/i);
      })
      .map(line => {
        const priceMatch = line.match(/(\d+\.\d{2})$/);
        return {
          name: line.replace(priceMatch?.[0] || '', '').trim(),
          quantity: 1,
          unit: 'pcs',
          price: priceMatch ? parseFloat(priceMatch[1]) : undefined
        };
      });

    return {
      storeName: storeMatch?.[0] || 'Unknown Store',
      purchaseDate: dateMatch?.[0] ? new Date(dateMatch[0]) : new Date(),
      totalAmount: totalMatch ? parseFloat(totalMatch[1]) : undefined,
      items: items
    };
  }

  openOcrPreviewDialog(data: any): void {
    const dialogRef = this.dialog.open(OcrPreviewDialogComponent, {
      width: '90vw',
      maxWidth: '800px',
      data: data
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result === 'upload') {
        // Handle upload to backend if needed
        this.router.navigate(['/receipts']);
      }
    });
  }
}