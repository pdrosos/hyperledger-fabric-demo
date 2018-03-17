import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap } from '@angular/router';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/switchMap';

import { Shipment } from '../core/models/shipment';
import { ShipmentService } from '../core/services/shipment.service';

@Component({
  selector: 'app-shipment-history',
  templateUrl: './shipment-history.component.html',
  styleUrls: ['./shipment-history.component.css']
})
export class ShipmentHistoryComponent implements OnInit {

  public trackingCode: string;
  public shipmentHistory$: Observable<Shipment[]>;
  public shipment: Shipment;

  constructor(private route: ActivatedRoute, private shipmentService: ShipmentService) {
  }

  ngOnInit() {
    this.shipmentHistory$ = this.route.paramMap
      .switchMap((params: ParamMap) => {
        this.trackingCode = params.get('trackingCode');

        return this.shipmentService.getShipmentHistory(this.trackingCode);
      })
    ;

    this.shipmentHistory$.subscribe((shipments: Shipment[]) => {
      // last history item
      this.shipment = shipments[shipments.length - 1];
    });
  }
}
