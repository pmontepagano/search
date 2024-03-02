import argparse
import asyncio
import decimal
import logging
import signal
import sys

import simplejson as json

# from grpclib.server import Server
from grpclib.client import Channel
from lib.search import v1 as search


logger = logging.getLogger(__name__)
log_formatter = logging.Formatter("%(asctime)s %(levelname)s %(message)s")
log_handler = logging.StreamHandler(sys.stdout)
log_handler.setFormatter(log_formatter)
logger.addHandler(log_handler)
logger.setLevel(logging.DEBUG)

PROVIDER_CONTRACT = """
.outputs Srv
.state graph
q0 ClientApp ? PurchaseRequest q1
q1 ClientApp ! TotalAmount q2
q2 ClientApp ? PurchaseWithPaymentNonce q3
q3 PPS ! RequestChargeWithNonce q4
q4 PPS ? ChargeOK q5
q4 PPS ? ChargeFail q6
q5 ClientApp ! PurchaseOK q7
q6 ClientApp ! PurchaseFail q8
.marking q0
.end""".encode(
    "utf-8"
)
PRICE_PER_BOOK = decimal.Decimal("100.00")
SHIPPING_COST = decimal.Decimal("10.00")


async def main(grpc_channel):
    stub = search.PrivateMiddlewareServiceStub(grpc_channel)
    registered = False
    logger.info("Connected to middleware. Waiting for registration...")
    async for r in stub.register_app(
        search.RegisterAppRequest(
            provider_contract=search.LocalContract(
                format=search.LocalContractFormat.LOCAL_CONTRACT_FORMAT_FSA,
                contract=PROVIDER_CONTRACT,
            )
        )
    ):
        if registered and r.notification:
            logger.info(f"Notification received: {r.notification}")
            # Start a new session for this channel.
            asyncio.create_task(session(grpc_channel, r.notification))
        elif not registered and r.app_id:
            # This should only happen once, in the first iteration.
            registered = True
            logger.info(f"App registered with id {r.app_id}")
            # Create temp file for Docker Compose healthcheck.
            with open("/tmp/registered", "w") as f:
                f.write("OK")
        else:
            logger.error(f"Unexpected response: {r}. Exiting.")
            break

    grpc_channel.close()


async def session(grpc_channel, notification):
    logger.info("Starting new session...")
    stub = search.PrivateMiddlewareServiceStub(grpc_channel)
    channel_id = notification.channel_id

    # Receive purchase request from ClientApp.
    r = await stub.app_recv(
        search.AppRecvRequest(
            channel_id=channel_id,
            participant="ClientApp",
        )
    )
    if r.message.type != "PurchaseRequest":
        logger.error(f"Unexpected message type: {r.message.type}. Exiting.")
        return
    purchase_request = json.loads(r.message.body)
    logger.debug(f"Received purchase request: {purchase_request}")

    # Calculate total amount and send it to ClientApp.
    total_amount = decimal.Decimal(0)
    for book, quantity in purchase_request["items"].items():
        total_amount += PRICE_PER_BOOK * quantity
    total_amount += SHIPPING_COST
    logger.debug(f"Total amount: {total_amount}")
    r = await stub.app_send(
        search.AppSendRequest(
            channel_id=channel_id,
            recipient="ClientApp",
            message=search.AppMessage(
                type="TotalAmount",
                body=str(total_amount).encode("utf-8"),
            ),
        )
    )
    if r.result != search.AppSendResponseResult.RESULT_OK:
        logger.error(
            f"Failure sending message. Received AppSendResponse with status: {r.status}. Exiting."
        )
        return

    # Receive PurchaseWithPaymentNonce from ClientApp.
    r = await stub.app_recv(
        search.AppRecvRequest(
            channel_id=channel_id,
            participant="ClientApp",
        )
    )
    if r.message.type != "PurchaseWithPaymentNonce":
        logger.error(f"Unexpected message type: {r.message.type}. Exiting.")
        return
    purchase_with_payment_nonce = json.loads(r.message.body)
    logger.debug(f"Received payment nonce: {purchase_with_payment_nonce.get('nonce')}")

    # Send RequestChargeWithNonce to PPS.
    r = await stub.app_send(
        search.AppSendRequest(
            channel_id=channel_id,
            recipient="PPS",
            message=search.AppMessage(
                type="RequestChargeWithNonce",
                body=json.dumps(
                    {
                        "nonce": purchase_with_payment_nonce["nonce"],
                        "amount": total_amount,
                    }
                ).encode("utf-8"),
            ),
        )
    )
    if r.result != search.AppSendResponseResult.RESULT_OK:
        logger.error(
            f"Failure sending message. Received AppSendResponse with status: {r.status}. Exiting."
        )
        return

    # Receive ChargeOK or ChargeFail from PPS and send PurchaseOK or PurchaseFail to ClientApp.
    r = await stub.app_recv(
        search.AppRecvRequest(
            channel_id=channel_id,
            participant="PPS",
        )
    )
    if r.message.type == "ChargeOK":
        logger.debug(f"Charge OK.")
        r = await stub.app_send(
            search.AppSendRequest(
                channel_id=channel_id,
                recipient="ClientApp",
                message=search.AppMessage(
                    type="PurchaseOK",
                    body="".encode("utf-8"),
                ),
            )
        )
        if r.result != search.AppSendResponseResult.RESULT_OK:
            logger.error(
                f"Failure sending message. Received AppSendResponse with status: {r.status}. Exiting."
            )
            return
    elif r.message.type == "ChargeFail":
        logger.debug(f"Charge failed.")
        r = await stub.app_send(
            search.AppSendRequest(
                channel_id=channel_id,
                recipient="ClientApp",
                message=search.AppMessage(
                    type="PurchaseFail",
                    body="".encode("utf-8"),
                ),
            )
        )
        if r.result != search.AppSendResponseResult.RESULT_OK:
            logger.error(
                f"Failure sending message. Received AppSendResponse with status: {r.status}. Exiting."
            )
            return
    else:
        logger.error(f"Unexpected message type: {r.message.type}. Exiting.")
        return
    logger.info("Session completed.")


async def shutdown(signal, loop):
    """Cleanup tasks tied to the service's shutdown."""
    logger.info(f"Received exit signal {signal.name}...")
    tasks = [t for t in asyncio.all_tasks() if t is not asyncio.current_task()]

    [task.cancel() for task in tasks]

    logger.info(f"Cancelling {len(tasks)} outstanding tasks")
    await asyncio.gather(*tasks)
    loop.stop()


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Backend service for credit card example."
    )
    parser.add_argument("--middleware-host", type=str, default="middleware-backend")
    parser.add_argument("--middleware-port", type=int, default=11000)
    args = parser.parse_args()

    logger.debug("Starting backend service...")
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    signals = (signal.SIGHUP, signal.SIGTERM, signal.SIGINT)
    for s in signals:
        loop.add_signal_handler(s, lambda s=s: asyncio.create_task(shutdown(s, loop)))

    logger.info("Backend service starting. Connecting to middleware...")
    grpc_channel = Channel(host=args.middleware_host, port=args.middleware_port)

    try:
        loop.run_until_complete(main(grpc_channel=grpc_channel))
    finally:
        loop.close()
        grpc_channel.close()
        logger.info("Backend service stopped.")
