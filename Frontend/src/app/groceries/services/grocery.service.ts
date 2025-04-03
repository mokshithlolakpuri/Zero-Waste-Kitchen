import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { environment } from '../../../environments/environment';
import { GroceryItem } from '../models/grocery-item';

@Injectable({
  providedIn: 'root'
})
export class GroceryService {
  constructor(private http: HttpClient) {}

  getAllGroceries(): Observable<GroceryItem[]> {
    return this.http.get<GroceryItem[]>(`${environment.apiUrl}/groceries`).pipe(
      map((groceries) =>
        groceries.map((grocery) => ({
          ...grocery,
          storageLocation: this.formatStorageLocation(grocery.storageLocation),
          expiryDate: grocery.expiry_date,
          manufactureDate: grocery.manufacture_date
        }))
      )
    );
  }

  getGrocery(id: string): Observable<GroceryItem> {
    return this.http.get<GroceryItem>(`${environment.apiUrl}/groceries/${id}`);
  }

  createGrocery(grocery: Omit<GroceryItem, 'id'>): Observable<GroceryItem> {
    return this.http.post<GroceryItem>(`${environment.apiUrl}/groceries`, grocery);
  }

  updateGrocery(id: string, grocery: Partial<GroceryItem>): Observable<GroceryItem> {
    return this.http.put<GroceryItem>(`${environment.apiUrl}/groceries/${id}`, grocery);
  }

  deleteGrocery(id: string): Observable<void> {
    return this.http.delete<void>(`${environment.apiUrl}/groceries/${id}`);
  }

  private formatStorageLocation(location: string): string {
    switch (location) {
      case 'refrigerator':
        return 'Refrigerator';
      case 'deep_freeze':
        return 'Deep Freeze';
      case 'dry_pantry':
        return 'Dry Pantry';
      default:
        return location;
    }
  }
}