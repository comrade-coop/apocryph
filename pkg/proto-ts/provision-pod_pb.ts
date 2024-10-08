// SPDX-License-Identifier: GPL-3.0

// @generated by protoc-gen-es v1.6.0 with parameter "target=ts"
// @generated from file provision-pod.proto (package apocryph.proto.v0.provisionPod, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3, protoInt64 } from "@bufbuild/protobuf";
import { Pod } from "./pod_pb.js";

/**
 * @generated from message apocryph.proto.v0.provisionPod.ProvisionPodRequest
 */
export class ProvisionPodRequest extends Message<ProvisionPodRequest> {
  /**
   * @generated from field: apocryph.proto.v0.pod.Pod pod = 1;
   */
  pod?: Pod;

  /**
   * @generated from field: apocryph.proto.v0.provisionPod.PaymentChannel payment = 3;
   */
  payment?: PaymentChannel;

  constructor(data?: PartialMessage<ProvisionPodRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.provisionPod.ProvisionPodRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "pod", kind: "message", T: Pod },
    { no: 3, name: "payment", kind: "message", T: PaymentChannel },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ProvisionPodRequest {
    return new ProvisionPodRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ProvisionPodRequest {
    return new ProvisionPodRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ProvisionPodRequest {
    return new ProvisionPodRequest().fromJsonString(jsonString, options);
  }

  static equals(a: ProvisionPodRequest | PlainMessage<ProvisionPodRequest> | undefined, b: ProvisionPodRequest | PlainMessage<ProvisionPodRequest> | undefined): boolean {
    return proto3.util.equals(ProvisionPodRequest, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.provisionPod.DeletePodRequest
 */
export class DeletePodRequest extends Message<DeletePodRequest> {
  constructor(data?: PartialMessage<DeletePodRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.provisionPod.DeletePodRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DeletePodRequest {
    return new DeletePodRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DeletePodRequest {
    return new DeletePodRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DeletePodRequest {
    return new DeletePodRequest().fromJsonString(jsonString, options);
  }

  static equals(a: DeletePodRequest | PlainMessage<DeletePodRequest> | undefined, b: DeletePodRequest | PlainMessage<DeletePodRequest> | undefined): boolean {
    return proto3.util.equals(DeletePodRequest, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.provisionPod.DeletePodResponse
 */
export class DeletePodResponse extends Message<DeletePodResponse> {
  /**
   * @generated from field: bool success = 1;
   */
  success = false;

  /**
   * @generated from field: string error = 2;
   */
  error = "";

  constructor(data?: PartialMessage<DeletePodResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.provisionPod.DeletePodResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "success", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 2, name: "error", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): DeletePodResponse {
    return new DeletePodResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): DeletePodResponse {
    return new DeletePodResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): DeletePodResponse {
    return new DeletePodResponse().fromJsonString(jsonString, options);
  }

  static equals(a: DeletePodResponse | PlainMessage<DeletePodResponse> | undefined, b: DeletePodResponse | PlainMessage<DeletePodResponse> | undefined): boolean {
    return proto3.util.equals(DeletePodResponse, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.provisionPod.UpdatePodRequest
 */
export class UpdatePodRequest extends Message<UpdatePodRequest> {
  /**
   * @generated from field: apocryph.proto.v0.pod.Pod pod = 1;
   */
  pod?: Pod;

  /**
   * @generated from field: apocryph.proto.v0.provisionPod.PaymentChannel payment = 2;
   */
  payment?: PaymentChannel;

  constructor(data?: PartialMessage<UpdatePodRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.provisionPod.UpdatePodRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "pod", kind: "message", T: Pod },
    { no: 2, name: "payment", kind: "message", T: PaymentChannel },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): UpdatePodRequest {
    return new UpdatePodRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): UpdatePodRequest {
    return new UpdatePodRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): UpdatePodRequest {
    return new UpdatePodRequest().fromJsonString(jsonString, options);
  }

  static equals(a: UpdatePodRequest | PlainMessage<UpdatePodRequest> | undefined, b: UpdatePodRequest | PlainMessage<UpdatePodRequest> | undefined): boolean {
    return proto3.util.equals(UpdatePodRequest, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.provisionPod.UpdatePodResponse
 */
export class UpdatePodResponse extends Message<UpdatePodResponse> {
  /**
   * @generated from field: bool success = 1;
   */
  success = false;

  /**
   * @generated from field: string error = 2;
   */
  error = "";

  constructor(data?: PartialMessage<UpdatePodResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.provisionPod.UpdatePodResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "success", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 2, name: "error", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): UpdatePodResponse {
    return new UpdatePodResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): UpdatePodResponse {
    return new UpdatePodResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): UpdatePodResponse {
    return new UpdatePodResponse().fromJsonString(jsonString, options);
  }

  static equals(a: UpdatePodResponse | PlainMessage<UpdatePodResponse> | undefined, b: UpdatePodResponse | PlainMessage<UpdatePodResponse> | undefined): boolean {
    return proto3.util.equals(UpdatePodResponse, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.provisionPod.PaymentChannel
 */
export class PaymentChannel extends Message<PaymentChannel> {
  /**
   * @generated from field: bytes chainID = 1;
   */
  chainID = new Uint8Array(0);

  /**
   * @generated from field: bytes contractAddress = 2;
   */
  contractAddress = new Uint8Array(0);

  /**
   * @generated from field: bytes publisherAddress = 3;
   */
  publisherAddress = new Uint8Array(0);

  /**
   * @generated from field: bytes providerAddress = 4;
   */
  providerAddress = new Uint8Array(0);

  /**
   * @generated from field: bytes podID = 5;
   */
  podID = new Uint8Array(0);

  constructor(data?: PartialMessage<PaymentChannel>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.provisionPod.PaymentChannel";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "chainID", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 2, name: "contractAddress", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 3, name: "publisherAddress", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 4, name: "providerAddress", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 5, name: "podID", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PaymentChannel {
    return new PaymentChannel().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PaymentChannel {
    return new PaymentChannel().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PaymentChannel {
    return new PaymentChannel().fromJsonString(jsonString, options);
  }

  static equals(a: PaymentChannel | PlainMessage<PaymentChannel> | undefined, b: PaymentChannel | PlainMessage<PaymentChannel> | undefined): boolean {
    return proto3.util.equals(PaymentChannel, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.provisionPod.ProvisionPodResponse
 */
export class ProvisionPodResponse extends Message<ProvisionPodResponse> {
  /**
   * @generated from field: string error = 1;
   */
  error = "";

  /**
   * @generated from field: repeated apocryph.proto.v0.provisionPod.ProvisionPodResponse.ExposedHostPort addresses = 2;
   */
  addresses: ProvisionPodResponse_ExposedHostPort[] = [];

  /**
   * @generated from field: string namespace = 3;
   */
  namespace = "";

  constructor(data?: PartialMessage<ProvisionPodResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.provisionPod.ProvisionPodResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "error", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "addresses", kind: "message", T: ProvisionPodResponse_ExposedHostPort, repeated: true },
    { no: 3, name: "namespace", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ProvisionPodResponse {
    return new ProvisionPodResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ProvisionPodResponse {
    return new ProvisionPodResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ProvisionPodResponse {
    return new ProvisionPodResponse().fromJsonString(jsonString, options);
  }

  static equals(a: ProvisionPodResponse | PlainMessage<ProvisionPodResponse> | undefined, b: ProvisionPodResponse | PlainMessage<ProvisionPodResponse> | undefined): boolean {
    return proto3.util.equals(ProvisionPodResponse, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.provisionPod.ProvisionPodResponse.ExposedHostPort
 */
export class ProvisionPodResponse_ExposedHostPort extends Message<ProvisionPodResponse_ExposedHostPort> {
  /**
   * @generated from field: string multiaddr = 1;
   */
  multiaddr = "";

  /**
   * @generated from field: string containerName = 2;
   */
  containerName = "";

  /**
   * @generated from field: uint64 containerPort = 3;
   */
  containerPort = protoInt64.zero;

  constructor(data?: PartialMessage<ProvisionPodResponse_ExposedHostPort>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.provisionPod.ProvisionPodResponse.ExposedHostPort";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "multiaddr", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "containerName", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "containerPort", kind: "scalar", T: 4 /* ScalarType.UINT64 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ProvisionPodResponse_ExposedHostPort {
    return new ProvisionPodResponse_ExposedHostPort().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ProvisionPodResponse_ExposedHostPort {
    return new ProvisionPodResponse_ExposedHostPort().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ProvisionPodResponse_ExposedHostPort {
    return new ProvisionPodResponse_ExposedHostPort().fromJsonString(jsonString, options);
  }

  static equals(a: ProvisionPodResponse_ExposedHostPort | PlainMessage<ProvisionPodResponse_ExposedHostPort> | undefined, b: ProvisionPodResponse_ExposedHostPort | PlainMessage<ProvisionPodResponse_ExposedHostPort> | undefined): boolean {
    return proto3.util.equals(ProvisionPodResponse_ExposedHostPort, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.provisionPod.PodLogRequest
 */
export class PodLogRequest extends Message<PodLogRequest> {
  /**
   * @generated from field: string ContainerName = 1;
   */
  ContainerName = "";

  constructor(data?: PartialMessage<PodLogRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.provisionPod.PodLogRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "ContainerName", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PodLogRequest {
    return new PodLogRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PodLogRequest {
    return new PodLogRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PodLogRequest {
    return new PodLogRequest().fromJsonString(jsonString, options);
  }

  static equals(a: PodLogRequest | PlainMessage<PodLogRequest> | undefined, b: PodLogRequest | PlainMessage<PodLogRequest> | undefined): boolean {
    return proto3.util.equals(PodLogRequest, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.provisionPod.PodLogResponse
 */
export class PodLogResponse extends Message<PodLogResponse> {
  /**
   * @generated from field: apocryph.proto.v0.provisionPod.LogEntry logEntry = 1;
   */
  logEntry?: LogEntry;

  constructor(data?: PartialMessage<PodLogResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.provisionPod.PodLogResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "logEntry", kind: "message", T: LogEntry },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PodLogResponse {
    return new PodLogResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PodLogResponse {
    return new PodLogResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PodLogResponse {
    return new PodLogResponse().fromJsonString(jsonString, options);
  }

  static equals(a: PodLogResponse | PlainMessage<PodLogResponse> | undefined, b: PodLogResponse | PlainMessage<PodLogResponse> | undefined): boolean {
    return proto3.util.equals(PodLogResponse, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.provisionPod.LogEntry
 */
export class LogEntry extends Message<LogEntry> {
  /**
   * @generated from field: uint64 NanosecondsUnixEpoch = 1;
   */
  NanosecondsUnixEpoch = protoInt64.zero;

  /**
   * @generated from field: string line = 2;
   */
  line = "";

  constructor(data?: PartialMessage<LogEntry>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.provisionPod.LogEntry";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "NanosecondsUnixEpoch", kind: "scalar", T: 4 /* ScalarType.UINT64 */ },
    { no: 2, name: "line", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): LogEntry {
    return new LogEntry().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): LogEntry {
    return new LogEntry().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): LogEntry {
    return new LogEntry().fromJsonString(jsonString, options);
  }

  static equals(a: LogEntry | PlainMessage<LogEntry> | undefined, b: LogEntry | PlainMessage<LogEntry> | undefined): boolean {
    return proto3.util.equals(LogEntry, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.provisionPod.PodInfoRequest
 */
export class PodInfoRequest extends Message<PodInfoRequest> {
  /**
   * @generated from field: string namespace = 1;
   */
  namespace = "";

  constructor(data?: PartialMessage<PodInfoRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.provisionPod.PodInfoRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "namespace", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PodInfoRequest {
    return new PodInfoRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PodInfoRequest {
    return new PodInfoRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PodInfoRequest {
    return new PodInfoRequest().fromJsonString(jsonString, options);
  }

  static equals(a: PodInfoRequest | PlainMessage<PodInfoRequest> | undefined, b: PodInfoRequest | PlainMessage<PodInfoRequest> | undefined): boolean {
    return proto3.util.equals(PodInfoRequest, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.provisionPod.PodInfoResponse
 */
export class PodInfoResponse extends Message<PodInfoResponse> {
  /**
   * @generated from field: string info = 1;
   */
  info = "";

  constructor(data?: PartialMessage<PodInfoResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.provisionPod.PodInfoResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "info", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PodInfoResponse {
    return new PodInfoResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PodInfoResponse {
    return new PodInfoResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PodInfoResponse {
    return new PodInfoResponse().fromJsonString(jsonString, options);
  }

  static equals(a: PodInfoResponse | PlainMessage<PodInfoResponse> | undefined, b: PodInfoResponse | PlainMessage<PodInfoResponse> | undefined): boolean {
    return proto3.util.equals(PodInfoResponse, a, b);
  }
}

