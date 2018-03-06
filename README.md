# Hyperledger Fabric Demo
Hyperledger Fabric project for the SoftUni Blockchain Dev Camp - Sofia - February 2018

This is a simple online seller shipment and tracking use case.

Imagine that a big online seller works with several courier companies to send shipments to their customers. Both seller and courier companies need to keep real time ledger of shipments and track their whole lifecycle from creation to delivery.

Seller needs to keep a separate private ledger of shipments for each of their courier partners and only they and the courier to have access to write to it. Other courier companies the seller works with must not have access to this private ledger.

The projects consists of the following components:
 
[Hyperledger Fabric network](network/README.md) between the seller and one courier company.
 
[Chaincode](chaincode/README.md) in the seller-courier one channel.

[Seller web app](seller/README.md) <br>
Only seller can create a shipment and decide to which of their courier partners to pass it for delivery.

[Courier web app](courier/README.md) <br>
Only courier can change the shipment status and current location, once they receive it from the seller.

[Customer web app](customer/README.md) <br>
Seller's customer know their own tracking codes and should have a way to track current shipment location and history in real time.
