import { NgModule, Optional, SkipSelf } from '@angular/core';

import { ShipmentService } from './services/shipment.service';
import { ErrorHandlerService } from './services/error-handler.service';

@NgModule({
  imports: [],
  providers: [
    ShipmentService,
    ErrorHandlerService
  ]
})
export class CoreModule {
  constructor (@Optional() @SkipSelf() parentModule: CoreModule) {
      if (parentModule) {
          throw new Error('CoreModule is already loaded. Import it in the AppModule only');
      }
  }
}
