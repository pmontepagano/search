// @generated by protoc-gen-connect-es v1.1.3 with parameter "target=ts"
// @generated from file search/v1/broker.proto (package search.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { BrokerChannelRequest, BrokerChannelResponse, RegisterProviderRequest, RegisterProviderResponse } from "./broker_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service search.v1.BrokerService
 */
export const BrokerService = {
  typeName: "search.v1.BrokerService",
  methods: {
    /**
     * @generated from rpc search.v1.BrokerService.BrokerChannel
     */
    brokerChannel: {
      name: "BrokerChannel",
      I: BrokerChannelRequest,
      O: BrokerChannelResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc search.v1.BrokerService.RegisterProvider
     */
    registerProvider: {
      name: "RegisterProvider",
      I: RegisterProviderRequest,
      O: RegisterProviderResponse,
      kind: MethodKind.Unary,
    },
  }
} as const;

