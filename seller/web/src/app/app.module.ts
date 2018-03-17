import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { CommonModule, registerLocaleData } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import localeUk from '@angular/common/locales/uk';

import { CoreModule } from './core/core.module';
import { AppRoutingModule } from './app-routing/app-routing.module';

import { AppComponent } from './app.component';
import { NavbarComponent } from './navbar/navbar.component';
import { HomeComponent } from './home/home.component';
import { ShipmentStateComponent } from './shipment-state/shipment-state.component';
import { ShipmentHistoryComponent } from './shipment-history/shipment-history.component';
import { CreateShipmentComponent } from './create-shipment/create-shipment.component';

@NgModule({
  declarations: [
    AppComponent,
    NavbarComponent,
    HomeComponent,
    ShipmentStateComponent,
    ShipmentHistoryComponent,
    CreateShipmentComponent
  ],
  imports: [
    BrowserModule,
    CommonModule,
    FormsModule,
    HttpClientModule,
    AppRoutingModule,
    CoreModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }

// the second parameter 'uk' is optional
registerLocaleData(localeUk, 'uk');
