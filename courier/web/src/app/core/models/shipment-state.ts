export interface ShipmentState {
  state: string;
  location: {
    country: string;
    city: string;
    postCode: string;
    address: string;
  };
  isDelivered: boolean;
}
