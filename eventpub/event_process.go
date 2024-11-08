package eventpub

import (
	"context"
	"github.com/dapr-platform/common"
)

func PublishInternalMessage(ctx context.Context, msg *common.InternalMessage) error {

	err := common.GetDaprClient().PublishEvent(ctx, common.PUBSUB_NAME, common.INTERNAL_MESSAGE_TOPIC_NAME, msg)
	if err != nil {
		common.Logger.Error("publishInternalmessage error ", err)
	}
	common.Logger.Debug("PublishInternalMessage type=", msg.GetType())
	return err
}
