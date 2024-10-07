// SPDX-License-Identifier: GPL-3.0

// @generated by protoc-gen-es v1.6.0 with parameter "target=ts"
// @generated from file deployment.proto (package apocryph.proto.v0.deployment, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import { ProvisionPodResponse } from "./provision-pod_pb.js";
import { Key, KeyPair, VerificationDetails } from "./pod_pb.js";

/**
 * @generated from message apocryph.proto.v0.deployment.Deployment
 */
export class Deployment extends Message<Deployment> {
  /**
   * @generated from field: string podManifestFile = 1;
   */
  podManifestFile = "";

  /**
   * @generated from field: apocryph.proto.v0.deployment.ProviderConfig provider = 2;
   */
  provider?: ProviderConfig;

  /**
   * @generated from field: apocryph.proto.v0.deployment.PaymentChannelConfig payment = 3;
   */
  payment?: PaymentChannelConfig;

  /**
   * @generated from field: repeated apocryph.proto.v0.deployment.UploadedImage images = 4;
   */
  images: UploadedImage[] = [];

  /**
   * @generated from field: repeated apocryph.proto.v0.deployment.UploadedSecret secrets = 5;
   */
  secrets: UploadedSecret[] = [];

  /**
   * @generated from field: apocryph.proto.v0.provisionPod.ProvisionPodResponse deployed = 6;
   */
  deployed?: ProvisionPodResponse;

  /**
   * @generated from field: apocryph.proto.v0.pod.KeyPair keyPair = 7;
   */
  keyPair?: KeyPair;

  constructor(data?: PartialMessage<Deployment>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.deployment.Deployment";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "podManifestFile", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "provider", kind: "message", T: ProviderConfig },
    { no: 3, name: "payment", kind: "message", T: PaymentChannelConfig },
    { no: 4, name: "images", kind: "message", T: UploadedImage, repeated: true },
    { no: 5, name: "secrets", kind: "message", T: UploadedSecret, repeated: true },
    { no: 6, name: "deployed", kind: "message", T: ProvisionPodResponse },
    { no: 7, name: "keyPair", kind: "message", T: KeyPair },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Deployment {
    return new Deployment().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Deployment {
    return new Deployment().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Deployment {
    return new Deployment().fromJsonString(jsonString, options);
  }

  static equals(a: Deployment | PlainMessage<Deployment> | undefined, b: Deployment | PlainMessage<Deployment> | undefined): boolean {
    return proto3.util.equals(Deployment, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.deployment.ProviderConfig
 */
export class ProviderConfig extends Message<ProviderConfig> {
  /**
   * @generated from field: bytes ethereumAddress = 1;
   */
  ethereumAddress = new Uint8Array(0);

  /**
   * @generated from field: string libp2pAddress = 2;
   */
  libp2pAddress = "";

  constructor(data?: PartialMessage<ProviderConfig>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.deployment.ProviderConfig";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "ethereumAddress", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 2, name: "libp2pAddress", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ProviderConfig {
    return new ProviderConfig().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ProviderConfig {
    return new ProviderConfig().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ProviderConfig {
    return new ProviderConfig().fromJsonString(jsonString, options);
  }

  static equals(a: ProviderConfig | PlainMessage<ProviderConfig> | undefined, b: ProviderConfig | PlainMessage<ProviderConfig> | undefined): boolean {
    return proto3.util.equals(ProviderConfig, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.deployment.PaymentChannelConfig
 */
export class PaymentChannelConfig extends Message<PaymentChannelConfig> {
  /**
   * @generated from field: bytes chainID = 1;
   */
  chainID = new Uint8Array(0);

  /**
   * @generated from field: bytes paymentContractAddress = 2;
   */
  paymentContractAddress = new Uint8Array(0);

  /**
   * @generated from field: bytes publisherAddress = 3;
   */
  publisherAddress = new Uint8Array(0);

  /**
   * @generated from field: bytes podID = 5;
   */
  podID = new Uint8Array(0);

  constructor(data?: PartialMessage<PaymentChannelConfig>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.deployment.PaymentChannelConfig";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "chainID", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 2, name: "paymentContractAddress", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 3, name: "publisherAddress", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 5, name: "podID", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PaymentChannelConfig {
    return new PaymentChannelConfig().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PaymentChannelConfig {
    return new PaymentChannelConfig().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PaymentChannelConfig {
    return new PaymentChannelConfig().fromJsonString(jsonString, options);
  }

  static equals(a: PaymentChannelConfig | PlainMessage<PaymentChannelConfig> | undefined, b: PaymentChannelConfig | PlainMessage<PaymentChannelConfig> | undefined): boolean {
    return proto3.util.equals(PaymentChannelConfig, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.deployment.UploadedImage
 */
export class UploadedImage extends Message<UploadedImage> {
  /**
   * @generated from field: string sourceUrl = 1;
   */
  sourceUrl = "";

  /**
   * @generated from field: string digest = 2;
   */
  digest = "";

  /**
   * @generated from field: bytes cid = 3;
   */
  cid = new Uint8Array(0);

  /**
   * @generated from field: apocryph.proto.v0.pod.Key key = 4;
   */
  key?: Key;

  /**
   * @generated from field: apocryph.proto.v0.pod.VerificationDetails verificationDetails = 5;
   */
  verificationDetails?: VerificationDetails;

  constructor(data?: PartialMessage<UploadedImage>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.deployment.UploadedImage";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "sourceUrl", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "digest", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "cid", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 4, name: "key", kind: "message", T: Key },
    { no: 5, name: "verificationDetails", kind: "message", T: VerificationDetails },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): UploadedImage {
    return new UploadedImage().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): UploadedImage {
    return new UploadedImage().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): UploadedImage {
    return new UploadedImage().fromJsonString(jsonString, options);
  }

  static equals(a: UploadedImage | PlainMessage<UploadedImage> | undefined, b: UploadedImage | PlainMessage<UploadedImage> | undefined): boolean {
    return proto3.util.equals(UploadedImage, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.deployment.UploadedSecret
 */
export class UploadedSecret extends Message<UploadedSecret> {
  /**
   * @generated from field: string volumeName = 1;
   */
  volumeName = "";

  /**
   * @generated from field: bytes sha256sum = 2;
   */
  sha256sum = new Uint8Array(0);

  /**
   * @generated from field: bytes cid = 3;
   */
  cid = new Uint8Array(0);

  /**
   * @generated from field: apocryph.proto.v0.pod.Key key = 4;
   */
  key?: Key;

  constructor(data?: PartialMessage<UploadedSecret>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.deployment.UploadedSecret";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "volumeName", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "sha256sum", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 3, name: "cid", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 4, name: "key", kind: "message", T: Key },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): UploadedSecret {
    return new UploadedSecret().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): UploadedSecret {
    return new UploadedSecret().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): UploadedSecret {
    return new UploadedSecret().fromJsonString(jsonString, options);
  }

  static equals(a: UploadedSecret | PlainMessage<UploadedSecret> | undefined, b: UploadedSecret | PlainMessage<UploadedSecret> | undefined): boolean {
    return proto3.util.equals(UploadedSecret, a, b);
  }
}

