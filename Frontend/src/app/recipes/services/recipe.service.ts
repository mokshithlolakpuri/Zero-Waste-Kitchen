import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { Recipe } from '../models/recipe';
import { Observable, throwError } from 'rxjs';
import { map, catchError } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class RecipeService {
  private apiUrl = `${environment.apiUrl}/recipes`;

  constructor(private http: HttpClient) {}

  generateRecipes(cuisine: string = ''): Observable<Recipe[]> {
    return this.http.post<Recipe[]>(`${this.apiUrl}/generate`, { cuisine });
  }

  getAllRecipes(): Observable<Recipe[]> {
    return this.http.get<{ recipes: Recipe[] }>(this.apiUrl).pipe(
      map(response => response.recipes), // Extract the recipes array
      catchError(err => {
        console.error('Failed to load recipes', err);
        return throwError(() => err);
      })
    );
  }

  getRecipeById(id: number): Observable<Recipe> {
    return this.http.get<{ recipe: Recipe }>(`${this.apiUrl}/${id}`).pipe(
      map(response => response.recipe), // Extract the recipe object
      catchError(err => {
        console.error('Failed to load recipe', err);
        return throwError(() => err);
      })
    );
  }
}