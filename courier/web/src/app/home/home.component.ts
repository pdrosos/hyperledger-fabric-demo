import { Component, OnInit } from '@angular/core';

import { Observable } from 'rxjs/Observable';

import { Shipment } from '../core/models/shipment';
import { ShipmentService } from '../core/services/shipment.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

  public shipments$: Observable<Shipment[]>;

  constructor(private shipmentService: ShipmentService) {
  }

  ngOnInit() {
    this.shipments$ = this.shipmentService.getShipments();
  }
}
