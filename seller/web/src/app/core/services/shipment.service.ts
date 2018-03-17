import { Injectable } from '@angular/core';

import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';
import { catchError } from 'rxjs/operators';
import 'rxjs/add/operator/map';

import { API_BASE_URL } from '../../constants/constants';
import { Shipment } from '../models/shipment';
import { ErrorHandlerService } from './error-handler.service';

const httpOptions = {
  headers: new HttpHeaders({
    'Content-Type': 'application/json'
  })
};

@Injectable()
export class ShipmentService {

  private baseResourceUrl: string;

  constructor(private http: HttpClient, private errorHandlerService: ErrorHandlerService) {
    this.baseResourceUrl = API_BASE_URL + '/shipments'
  }

  public getShipments(): Observable<Shipment[]> {
    return this.http.get(this.baseResourceUrl)
      .map((shipments: Shipment[]) => {
        shipments.forEach((shipment: Shipment) => {
          shipment.createdAt = new Date(shipment.createdAt);
          shipment.updatedAt = new Date(shipment.updatedAt);
        });

        return shipments;
      })
      .pipe(
        catchError(this.errorHandlerService.handleError)
      )
    ;
  }

  public getShipmentHistory(trackingCode: string): Observable<Shipment[]> {
    let url = '/' + trackingCode + '/history';

    return this.http.get(this.baseResourceUrl + url)
      .map((shipments: Shipment[]) => {
        shipments.forEach((shipment: Shipment) => {
          shipment.createdAt = new Date(shipment.createdAt);
          shipment.updatedAt = new Date(shipment.updatedAt);
        });

        return shipments;
      })
      .pipe(
        catchError(this.errorHandlerService.handleError)
      )
    ;
  }

  public getShipmentByTrackingCode(trackingCode: string): Observable<Shipment> {
    let url = '/' + trackingCode;

    return this.http.get(this.baseResourceUrl + url)
      .map((shipment: Shipment) => {
        shipment.createdAt = new Date(shipment.createdAt);
        shipment.updatedAt = new Date(shipment.updatedAt);

        return shipment;
      })
      .pipe(
        catchError(this.errorHandlerService.handleError)
      )
    ;
  }

  public createShipment(shipment: Shipment): Observable<boolean> {
    return this.http
      .post(this.baseResourceUrl, shipment, httpOptions)
      .map(() => {
        return true;
      })
      .pipe(
        catchError(this.errorHandlerService.handleError)
      )
    ;
  }
}
