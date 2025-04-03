// receipt.service.ts
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, catchError, map, throwError } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Receipt } from '../models/receipt';
import { MatSnackBar } from '@angular/material/snack-bar';
import { OcrResult } from '../models/receipt';

@Injectable({
  providedIn: 'root'
})
export class ReceiptService {
  private apiUrl = `${environment.apiUrl}/receipts`;

  constructor(
    private http: HttpClient,
    private snackBar: MatSnackBar
  ) {}

  getAllReceipts(): Observable<Receipt[]> {
    return this.http.get<Receipt[]>(this.apiUrl).pipe(
      catchError(err => {
        this.showError('Failed to load receipts');
        return throwError(() => err);
      })
    );
  }

  getReceipt(id: string): Observable<Receipt> {
    return this.http.get<Receipt>(`${this.apiUrl}/${id}`).pipe(
      catchError(err => {
        this.showError('Failed to load receipt details');
        return throwError(() => err);
      })
    );
  }

  uploadReceipt(formData: FormData): Observable<any> {
      return this.http.post(`${this.apiUrl}/upload`, formData).pipe(
          catchError(error => {
              this.snackBar.open('Upload failed', 'Close', { duration: 3000 });
              return throwError(() => error);
          })
      );
  }

  processReceipt(file: File): Observable<OcrResult> {
    const formData = new FormData();
    formData.append('file', file);

    return this.http.post<OcrResult>(`${this.apiUrl}/process`, formData).pipe(
      catchError(err => {
        this.showError('OCR processing failed');
        return throwError(() => err);
      })
    );
  }

  deleteReceipt(id: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`).pipe(
      catchError(err => {
        this.showError('Failed to delete receipt');
        return throwError(() => err);
      })
    );
  }

  showError(message: string): void {
    this.snackBar.open(message, 'Close', {
      duration: 5000,
      panelClass: ['error-snackbar']
    });
  }
}