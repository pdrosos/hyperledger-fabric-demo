import { Injectable } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';

import { ErrorObservable } from 'rxjs/observable/ErrorObservable';

@Injectable()
export class ErrorHandlerService {

  constructor() {
  }

  public handleError(error: HttpErrorResponse) {
    let messageFromError: string;

    if (error.error instanceof ErrorEvent) {
      // A client-side or network error occurred. Handle it accordingly.
      messageFromError = `An error occurred: ${error.error.message}`;
    } else {
      // The backend returned an unsuccessful response code.
      // The response body may contain clues as to what went wrong
      messageFromError = `Server returned error ${error.status}`;

      console.error(
        `Server returned error ${error.status}, ` +
        `body was: ${error.error}`);
    }

    // return an ErrorObservable with a user-facing error message
    return new ErrorObservable(messageFromError);
  };
}
