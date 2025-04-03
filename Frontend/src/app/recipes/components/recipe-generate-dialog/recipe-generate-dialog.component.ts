import { Component } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { RecipeService } from '../../services/recipe.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatOptionModule } from '@angular/material/core';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatDialogModule } from '@angular/material/dialog';

@Component({
  selector: 'app-recipe-generate-dialog',
  templateUrl: './recipe-generate-dialog.component.html',
  styleUrls: ['./recipe-generate-dialog.component.css'],
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatFormFieldModule,
    MatSelectModule,
    MatOptionModule,
    MatButtonModule,
    MatProgressSpinnerModule,
    MatDialogModule
  ]
})
export class RecipeGenerateDialogComponent {
  cuisine = '';
  isLoading = false;
  cuisines = [
    'Any', 'Italian', 'Mexican', 'Indian', 'Chinese', 
    'Japanese', 'Mediterranean', 'American', 'Thai', 'French'
  ];

  constructor(
    private dialogRef: MatDialogRef<RecipeGenerateDialogComponent>,
    private recipeService: RecipeService,
    private snackBar: MatSnackBar
  ) {}

  generateRecipes(): void {
    this.isLoading = true;
    this.recipeService.generateRecipes(this.cuisine === 'Any' ? '' : this.cuisine).subscribe({
      next: () => {
        this.snackBar.open('Recipes generated successfully!', 'Close', {
          duration: 3000
        });
        this.dialogRef.close(true);
      },
      error: (err) => {
        console.error('Failed to generate recipes', err);
        this.snackBar.open('Failed to generate recipes', 'Close', {
          duration: 3000
        });
        this.isLoading = false;
      }
    });
  }
}