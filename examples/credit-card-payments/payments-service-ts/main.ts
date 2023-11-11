import { createPromiseClient } from "@connectrpc/connect";
import { PrivateMiddlewareService } from "search/v1/middleware_connect";
import { LocalContract, LocalContractFormat } from "search/v1/contracts_pb";
import { createConnectTransport } from "@connectrpc/connect-node";

const transport = createConnectTransport({
//   baseUrl: "http://middleware-payments:11000",
  baseUrl: "http://localhost:11000",
  httpVersion: "2"
});


const ppsContract = new LocalContract();
ppsContract.format = LocalContractFormat.FSA;
ppsContract.contract = Buffer.from(`
.outputs PPS
.state graph
q0 ClientApp ? CardDetailsWithTotalAmount q1
q1 ClientApp ! PaymentNonce q2
q2 Srv ? RequestChargeWithNonce q3
q3 Srv ! ChargeOK q4
q3 Srv ! ChargeFail q5
.marking q0
.end`);


async function main() {
  const client = createPromiseClient(PrivateMiddlewareService, transport);
  const res = await client.registerApp({
    providerContract: ppsContract,
  });
  console.log(res);
}
void main();
