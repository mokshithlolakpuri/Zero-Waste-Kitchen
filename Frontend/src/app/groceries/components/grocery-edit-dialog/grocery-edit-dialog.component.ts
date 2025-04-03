import { Component, Inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators, AbstractControl, ValidationErrors, ValidatorFn } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef, MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatSelectModule } from '@angular/material/select';
import { MatNativeDateModule } from '@angular/material/core';
import { GroceryService } from '../../services/grocery.service';
import { GroceryItem } from '../../models/grocery-item';

// Custom validator to check if the date is in the past or today
function pastDateValidator(): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    const selectedDate = new Date(control.value);
    const today = new Date();
    today.setHours(0, 0, 0, 0); // Reset time to midnight for comparison

    if (control.value && selectedDate <= today) {
      return { pastDate: true };
    }
    return null;
  };
}

@Component({
  selector: 'app-grocery-edit-dialog',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatDatepickerModule,
    MatSelectModule,
    MatNativeDateModule
  ],
  templateUrl: './grocery-edit-dialog.component.html',
  styleUrls: ['./grocery-edit-dialog.component.css']
})
export class GroceryEditDialogComponent {
  groceryForm: FormGroup;
  isEdit = false;
  loading = false;
  errorMessage = '';

  constructor(
    private fb: FormBuilder,
    private groceryService: GroceryService,
    public dialogRef: MatDialogRef<GroceryEditDialogComponent>, // Changed to public
    @Inject(MAT_DIALOG_DATA) public data: { grocery?: GroceryItem }
  ) {
    this.groceryForm = this.fb.group({
      name: ['', Validators.required],
      quantity: [0, [Validators.required, Validators.min(0.1)]],
      unit: ['', Validators.required],
      expiryDate: ['', [Validators.required, pastDateValidator()]], // Added custom validator
      storageLocation: ['fridge', Validators.required]
    });

    if (data?.grocery) {
      this.isEdit = true;
      this.groceryForm.patchValue({
        ...data.grocery,
        expiryDate: data.grocery.expiry_date
          ? new Date(data.grocery.expiry_date)
          : ''
      });
    }
  }

  onSubmit(): void {
    if (this.groceryForm.invalid) return;

    this.loading = true;
    this.errorMessage = '';
    const formValue = this.groceryForm.value;

    const groceryData = {
      ...formValue,
      quantity: Number(formValue.quantity),
      expiry_date: formValue.expiryDate ? formValue.expiryDate.toISOString() : null
    };

    const operation = this.isEdit && this.data.grocery?.id
      ? this.groceryService.updateGrocery(String(this.data.grocery.id), groceryData)
      : this.groceryService.createGrocery(groceryData);

    operation.subscribe({
      next: () => {
        this.loading = false;
        this.dialogRef.close(true);
      },
      error: (err) => {
        this.loading = false;
        this.errorMessage = err?.error?.message || 'Failed to save grocery item. Please try again.';
      }
    });
  }
}