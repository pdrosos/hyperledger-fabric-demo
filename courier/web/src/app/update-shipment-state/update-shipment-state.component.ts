import { Component, OnInit, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import {ActivatedRoute, ParamMap, Router} from '@angular/router';

import { ShipmentState } from '../core/models/shipment-state';
import { ShipmentService } from '../core/services/shipment.service';

@Component({
  selector: 'app-update-shipment-state',
  templateUrl: './update-shipment-state.component.html',
  styleUrls: ['./update-shipment-state.component.css']
})
export class UpdateShipmentStateComponent implements OnInit {

  public trackingCode: string;
  public shipmentState: ShipmentState;
  @ViewChild('updateShipmentStateForm') public updateShipmentStateForm: NgForm;
  public submitDisabled: boolean = false;

  constructor(private router: Router, private route: ActivatedRoute, private shipmentService: ShipmentService) {
    this.shipmentState = <ShipmentState>{
      location: {}
    };
  }

  ngOnInit() {
    this.route.paramMap
      .subscribe((params: ParamMap) => {
        this.trackingCode = params.get('trackingCode');
      })
    ;
  }

  public onUpdateShipmentStateSubmit(): void {
    if (!this.updateShipmentStateForm.valid) {
      return;
    }

    this.submitDisabled = true;

    this.shipmentState.isDelivered = (this.shipmentState.state === 'Delivered');

    this.shipmentService.updateShipmentState(this.trackingCode, this.shipmentState).subscribe(() => {
        this.router.navigate(['/']);
      },
      (error) => {
        this.router.navigate(['/']);
      }
    );
  }
}
