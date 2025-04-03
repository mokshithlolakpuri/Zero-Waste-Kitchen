// receipt.module.ts
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBarModule } from '@angular/material/snack-bar';

import { ReceiptListComponent } from './components/receipt-list/receipt-list.component';
import { ReceiptUploadComponent } from './components/receipt-upload/receipt-upload.component';
import { ReceiptDetailComponent } from './components/receipt-detail/receipt-detail.component';
import { AuthGuard } from '../core/guards/auth.guard';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatListModule } from '@angular/material/list';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatDialogModule } from '@angular/material/dialog';
import { MatNativeDateModule } from '@angular/material/core';

const routes: Routes = [
  { path: '', component: ReceiptListComponent, canActivate: [AuthGuard] },
  { path: 'upload', component: ReceiptUploadComponent },
  { path: ':id', component: ReceiptDetailComponent }
];

@NgModule({
  imports: [
    CommonModule,
    RouterModule.forChild(routes),
    MatButtonModule,
    MatCardModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatIconModule,
    MatToolbarModule,
    MatListModule,
    MatInputModule,
    MatFormFieldModule,
    FormsModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatNativeDateModule
  ],
  exports: [RouterModule]
})
export class ReceiptsModule {}