# Hyperledger Fabric Demo
Hyperledger Fabric project for the SoftUni Blockchain Dev Camp - Sofia - February 2018

This is a simple online seller shipment and tracking use case.

Imagine that a big online seller works with several courier companies to send shipments to their customers. Both seller and courier companies need to keep real time ledger of shipments and track their whole lifecycle from creation to delivery.

Seller needs to keep a separate private ledger of shipments for each of their courier partners and only they and the courier to have access to write to it. Other courier companies the seller works with must not have access to this private ledger.

The projects consists of the following components:
 
[Hyperledger Fabric network](fabric-starter/README.md) between the seller and one courier company.
 
[Chaincode](fabric-starter/chaincode/go/shipment) installed in the seller-courier1 channel.

[Seller api and web app](seller/api/README.md) <br>
Only seller can create a shipment and decide to which of their courier partners to pass it for delivery.

[Courier api and web app](courier/api/README.md) <br>
Only courier can change the shipment state and current location, once they receive it from the seller.

[Customer api and web app](customer/api/README.md) <br>
Seller's customer know their own tracking codes and should have a way to track current shipment location and history in real time.

## Acknowledgements

The Hyperledger Fabric network uses a very helpful [fabric-starter](https://github.com/olegabu/fabric-starter) script, which integrates [fabric-rest](https://github.com/Altoros/fabric-rest) REST API server and admin web app.

The scripts are inspired by [first-network](https://github.com/hyperledger/fabric-samples/tree/release/first-network) and 
 [balance-transfer](https://github.com/hyperledger/fabric-samples/tree/release/balance-transfer) of Hyperledger Fabric samples.