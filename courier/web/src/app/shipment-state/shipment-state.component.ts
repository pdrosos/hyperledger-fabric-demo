import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap } from '@angular/router';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/switchMap';

import { Shipment } from '../core/models/shipment';
import { ShipmentService } from '../core/services/shipment.service';

@Component({
  selector: 'app-shipment-state',
  templateUrl: './shipment-state.component.html',
  styleUrls: ['./shipment-state.component.css']
})
export class ShipmentStateComponent implements OnInit {

  public trackingCode: string;
  public shipment$: Observable<Shipment>;
  public shipment: Shipment;

  constructor(private route: ActivatedRoute, private shipmentService: ShipmentService) {
  }

  ngOnInit() {
    this.shipment$ = this.route.paramMap
      .switchMap((params: ParamMap) => {
        this.trackingCode = params.get('trackingCode');

        return this.shipmentService.getShipmentByTrackingCode(this.trackingCode);
      })
    ;

    this.shipment$.subscribe((shipment: Shipment) => {
      this.shipment = shipment;
    });
  }
}
