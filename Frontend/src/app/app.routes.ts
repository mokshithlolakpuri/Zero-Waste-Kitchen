import { Routes } from '@angular/router';
import { AuthGuard } from './core/guards/auth.guard';
import { guestGuard } from './core/guards/guest.guard';

export const routes: Routes = [
  {
    path: '',
    redirectTo: 'groceries',
    pathMatch: 'full'
  },
  {
    path: 'auth',
    loadChildren: () => import('./auth/auth.module').then(m => m.AuthModule),
    canActivate: [guestGuard]
  },
  {
    path: 'groceries',
    loadChildren: () => import('./groceries/groceries.module').then(m => m.GroceriesModule),
    canActivate: [AuthGuard]
  },
  {
    path: 'receipts',
    loadChildren: () => import('./receipts/receipts.module').then(m => m.ReceiptsModule),
    canActivate: [AuthGuard]
  },
  {
    path: 'recipes',
    loadChildren: () => import('./recipes/recipes.module').then(m => m.RecipesModule),
    canActivate: [AuthGuard]
  },
  {
    path: 'admin',
    loadChildren: () => import('./admin/admin.module').then(m => m.AdminModule),
    canActivate: [AuthGuard]
  },
  {
    path: '**',
    redirectTo: 'groceries'
  }
];