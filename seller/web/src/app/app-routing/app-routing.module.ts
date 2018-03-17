import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { HomeComponent } from '../home/home.component';
import { ShipmentStateComponent } from '../shipment-state/shipment-state.component';
import { ShipmentHistoryComponent } from '../shipment-history/shipment-history.component';
import { CreateShipmentComponent } from '../create-shipment/create-shipment.component';

const appRoutes: Routes = [
  {
    path: 'shipments/new',
    component: CreateShipmentComponent
  },
  {
    path: 'shipments/:trackingCode/history',
    component: ShipmentHistoryComponent
  },
  {
    path: 'shipments/:trackingCode',
    component: ShipmentStateComponent
  },
  {
    path: '',
    component: HomeComponent
  }
];

@NgModule({
  imports: [
    RouterModule.forRoot(appRoutes)
  ],
  exports: [
    RouterModule
  ],
  providers: []
})
export class AppRoutingModule { }
