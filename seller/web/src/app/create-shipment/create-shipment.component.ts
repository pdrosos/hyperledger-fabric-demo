import { Component, OnInit, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import { Router } from '@angular/router';

import { Shipment } from '../core/models/shipment';
import { ShipmentService } from '../core/services/shipment.service';

@Component({
  selector: 'app-create-shipment',
  templateUrl: './create-shipment.component.html',
  styleUrls: ['./create-shipment.component.css']
})
export class CreateShipmentComponent implements OnInit {

  public shipment: Shipment;
  @ViewChild('createShipmentForm') public createShipmentForm: NgForm;
  public submitDisabled: boolean = false;

  constructor(private router: Router, private shipmentService: ShipmentService) {
    this.shipment = <Shipment>{
      sender: {},
      recipient: {}
    };
  }

  ngOnInit() {
  }

  public onCreateShipmentSubmit(): void {
    if (!this.createShipmentForm.valid) {
      return;
    }

    this.submitDisabled = true;

    this.shipmentService.createShipment(this.shipment).subscribe(() => {
        this.router.navigate(['/']);
      },
      (error) => {
        this.router.navigate(['/']);
      }
    );
  }
}
