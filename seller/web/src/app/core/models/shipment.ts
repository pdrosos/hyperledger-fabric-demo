export interface Shipment {
  trackingCode: string;
  courier: string;
  shippingType: string;
  sender: {
    firstName: string;
    lastName: string;
    country: string;
    city: string;
    postCode: string;
    address: string;
    phone: string;
  },
  recipient: {
    firstName: string;
    lastName: string;
    country: string;
    city: string;
    postCode: string;
    address: string;
    phone: string;
  },
  weightInGrams: number;
  heightInMM: number;
  widthInMM: number;
  depthInMM: number;
  content: string;
  isFragile: boolean;
  lastState: string;
  lastLocation?: {
    country: string;
    city: string;
    postCode: string;
    address: string;
  };
  isDelivered: boolean;
  createdAt: Date;
  updatedAt: Date;
}
