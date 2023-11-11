import { createPromiseClient } from "@connectrpc/connect";
import { PrivateMiddlewareService } from "./lib/search/v1/middleware_connect";
import { createConnectTransport } from "@connectrpc/connect-node";

const transport = createConnectTransport({
//   baseUrl: "http://middleware-payments:11000",
  baseUrl: "http://localhost:11000",
  httpVersion: "2"
});

async function main() {
  const client = createPromiseClient(PrivateMiddlewareService, transport);
  const res = await client.registerApp({
    provider_contract: 'caca'
  });
  console.log(res);
}
void main();