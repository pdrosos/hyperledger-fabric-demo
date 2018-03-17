import { Component, OnInit } from '@angular/core';

import { ShipmentService } from '../core/services/shipment.service';

@Component({
  selector: 'app-create-shipment',
  templateUrl: './create-shipment.component.html',
  styleUrls: ['./create-shipment.component.css']
})
export class CreateShipmentComponent implements OnInit {
  constructor(private shipmentService: ShipmentService) {
  }

  ngOnInit() {
  }
}
