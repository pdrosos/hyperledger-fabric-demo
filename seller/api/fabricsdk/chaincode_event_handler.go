package fabricsdk

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type ChaincodeEventHandler interface {
	Handle(ccEvent *fab.CCEvent)
}
