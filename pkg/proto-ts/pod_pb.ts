// SPDX-License-Identifier: GPL-3.0

// @generated by protoc-gen-es v1.6.0 with parameter "target=ts"
// @generated from file pod.proto (package apocryph.proto.v0.pod, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3, protoInt64 } from "@bufbuild/protobuf";

/**
 * @generated from message apocryph.proto.v0.pod.Pod
 */
export class Pod extends Message<Pod> {
  /**
   * @generated from field: repeated apocryph.proto.v0.pod.Container containers = 1;
   */
  containers: Container[] = [];

  /**
   * @generated from field: repeated apocryph.proto.v0.pod.Volume volumes = 2;
   */
  volumes: Volume[] = [];

  /**
   * @generated from field: apocryph.proto.v0.pod.Replicas replicas = 3;
   */
  replicas?: Replicas;

  /**
   * @generated from field: apocryph.proto.v0.pod.KeyPair keyPair = 4;
   */
  keyPair?: KeyPair;

  /**
   * @generated from field: apocryph.proto.v0.pod.VerificationSettings verificationSettings = 5;
   */
  verificationSettings?: VerificationSettings;

  constructor(data?: PartialMessage<Pod>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.pod.Pod";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "containers", kind: "message", T: Container, repeated: true },
    { no: 2, name: "volumes", kind: "message", T: Volume, repeated: true },
    { no: 3, name: "replicas", kind: "message", T: Replicas },
    { no: 4, name: "keyPair", kind: "message", T: KeyPair },
    { no: 5, name: "verificationSettings", kind: "message", T: VerificationSettings },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Pod {
    return new Pod().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Pod {
    return new Pod().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Pod {
    return new Pod().fromJsonString(jsonString, options);
  }

  static equals(a: Pod | PlainMessage<Pod> | undefined, b: Pod | PlainMessage<Pod> | undefined): boolean {
    return proto3.util.equals(Pod, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.pod.Container
 */
export class Container extends Message<Container> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: apocryph.proto.v0.pod.Image image = 2;
   */
  image?: Image;

  /**
   * @generated from field: repeated string entrypoint = 3;
   */
  entrypoint: string[] = [];

  /**
   * @generated from field: repeated string command = 4;
   */
  command: string[] = [];

  /**
   * @generated from field: string workingDir = 5;
   */
  workingDir = "";

  /**
   * @generated from field: repeated apocryph.proto.v0.pod.Container.Port ports = 6;
   */
  ports: Container_Port[] = [];

  /**
   * @generated from field: map<string, string> env = 7;
   */
  env: { [key: string]: string } = {};

  /**
   * @generated from field: repeated apocryph.proto.v0.pod.Container.VolumeMount volumes = 8;
   */
  volumes: Container_VolumeMount[] = [];

  /**
   * "cpu", "memory", custom
   *
   * @generated from field: repeated apocryph.proto.v0.pod.Resource resourceRequests = 9;
   */
  resourceRequests: Resource[] = [];

  constructor(data?: PartialMessage<Container>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.pod.Container";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "image", kind: "message", T: Image },
    { no: 3, name: "entrypoint", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 4, name: "command", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 5, name: "workingDir", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "ports", kind: "message", T: Container_Port, repeated: true },
    { no: 7, name: "env", kind: "map", K: 9 /* ScalarType.STRING */, V: {kind: "scalar", T: 9 /* ScalarType.STRING */} },
    { no: 8, name: "volumes", kind: "message", T: Container_VolumeMount, repeated: true },
    { no: 9, name: "resourceRequests", kind: "message", T: Resource, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Container {
    return new Container().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Container {
    return new Container().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Container {
    return new Container().fromJsonString(jsonString, options);
  }

  static equals(a: Container | PlainMessage<Container> | undefined, b: Container | PlainMessage<Container> | undefined): boolean {
    return proto3.util.equals(Container, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.pod.Container.Port
 */
export class Container_Port extends Message<Container_Port> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: uint64 containerPort = 2;
   */
  containerPort = protoInt64.zero;

  /**
   * @generated from field: uint64 servicePort = 3;
   */
  servicePort = protoInt64.zero;

  /**
   * @generated from oneof apocryph.proto.v0.pod.Container.Port.exposedPort
   */
  exposedPort: {
    /**
     * @generated from field: string hostHttpHost = 4;
     */
    value: string;
    case: "hostHttpHost";
  } | {
    /**
     * string hostHttpsHost = 4;
     *
     * uint64 hostUdpPort = 6;
     * uint64 servicePort = 7;
     *
     * @generated from field: uint64 hostTcpPort = 5;
     */
    value: bigint;
    case: "hostTcpPort";
  } | { case: undefined; value?: undefined } = { case: undefined };

  constructor(data?: PartialMessage<Container_Port>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.pod.Container.Port";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "containerPort", kind: "scalar", T: 4 /* ScalarType.UINT64 */ },
    { no: 3, name: "servicePort", kind: "scalar", T: 4 /* ScalarType.UINT64 */ },
    { no: 4, name: "hostHttpHost", kind: "scalar", T: 9 /* ScalarType.STRING */, oneof: "exposedPort" },
    { no: 5, name: "hostTcpPort", kind: "scalar", T: 4 /* ScalarType.UINT64 */, oneof: "exposedPort" },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Container_Port {
    return new Container_Port().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Container_Port {
    return new Container_Port().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Container_Port {
    return new Container_Port().fromJsonString(jsonString, options);
  }

  static equals(a: Container_Port | PlainMessage<Container_Port> | undefined, b: Container_Port | PlainMessage<Container_Port> | undefined): boolean {
    return proto3.util.equals(Container_Port, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.pod.Container.VolumeMount
 */
export class Container_VolumeMount extends Message<Container_VolumeMount> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: string mountPath = 2;
   */
  mountPath = "";

  /**
   * @generated from field: bool readOnly = 3;
   */
  readOnly = false;

  constructor(data?: PartialMessage<Container_VolumeMount>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.pod.Container.VolumeMount";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "mountPath", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "readOnly", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Container_VolumeMount {
    return new Container_VolumeMount().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Container_VolumeMount {
    return new Container_VolumeMount().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Container_VolumeMount {
    return new Container_VolumeMount().fromJsonString(jsonString, options);
  }

  static equals(a: Container_VolumeMount | PlainMessage<Container_VolumeMount> | undefined, b: Container_VolumeMount | PlainMessage<Container_VolumeMount> | undefined): boolean {
    return proto3.util.equals(Container_VolumeMount, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.pod.Image
 */
export class Image extends Message<Image> {
  /**
   * @generated from field: bytes cid = 1;
   */
  cid = new Uint8Array(0);

  /**
   * @generated from field: apocryph.proto.v0.pod.Key key = 2;
   */
  key?: Key;

  /**
   * @generated from field: string url = 3;
   */
  url = "";

  /**
   * @generated from field: apocryph.proto.v0.pod.VerificationDetails verificationDetails = 4;
   */
  verificationDetails?: VerificationDetails;

  constructor(data?: PartialMessage<Image>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.pod.Image";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "cid", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 2, name: "key", kind: "message", T: Key },
    { no: 3, name: "url", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "verificationDetails", kind: "message", T: VerificationDetails },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Image {
    return new Image().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Image {
    return new Image().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Image {
    return new Image().fromJsonString(jsonString, options);
  }

  static equals(a: Image | PlainMessage<Image> | undefined, b: Image | PlainMessage<Image> | undefined): boolean {
    return proto3.util.equals(Image, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.pod.VerificationDetails
 */
export class VerificationDetails extends Message<VerificationDetails> {
  /**
   * @generated from field: string signature = 1;
   */
  signature = "";

  /**
   * @generated from field: string identity = 2;
   */
  identity = "";

  /**
   * @generated from field: string issuer = 3;
   */
  issuer = "";

  constructor(data?: PartialMessage<VerificationDetails>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.pod.VerificationDetails";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "signature", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "identity", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "issuer", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): VerificationDetails {
    return new VerificationDetails().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): VerificationDetails {
    return new VerificationDetails().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): VerificationDetails {
    return new VerificationDetails().fromJsonString(jsonString, options);
  }

  static equals(a: VerificationDetails | PlainMessage<VerificationDetails> | undefined, b: VerificationDetails | PlainMessage<VerificationDetails> | undefined): boolean {
    return proto3.util.equals(VerificationDetails, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.pod.Volume
 */
export class Volume extends Message<Volume> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: apocryph.proto.v0.pod.Volume.Type type = 2;
   */
  type = Volume_Type.VOLUME_EMPTY;

  /**
   * @generated from field: apocryph.proto.v0.pod.Volume.AccessMode accessMode = 3;
   */
  accessMode = Volume_AccessMode.VOLUME_RW_ONE;

  /**
   * @generated from oneof apocryph.proto.v0.pod.Volume.configuration
   */
  configuration: {
    /**
     * @generated from field: apocryph.proto.v0.pod.Volume.FilesystemConfig filesystem = 4;
     */
    value: Volume_FilesystemConfig;
    case: "filesystem";
  } | {
    /**
     * @generated from field: apocryph.proto.v0.pod.Volume.SecretConfig secret = 5;
     */
    value: Volume_SecretConfig;
    case: "secret";
  } | { case: undefined; value?: undefined } = { case: undefined };

  constructor(data?: PartialMessage<Volume>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.pod.Volume";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "type", kind: "enum", T: proto3.getEnumType(Volume_Type) },
    { no: 3, name: "accessMode", kind: "enum", T: proto3.getEnumType(Volume_AccessMode) },
    { no: 4, name: "filesystem", kind: "message", T: Volume_FilesystemConfig, oneof: "configuration" },
    { no: 5, name: "secret", kind: "message", T: Volume_SecretConfig, oneof: "configuration" },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Volume {
    return new Volume().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Volume {
    return new Volume().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Volume {
    return new Volume().fromJsonString(jsonString, options);
  }

  static equals(a: Volume | PlainMessage<Volume> | undefined, b: Volume | PlainMessage<Volume> | undefined): boolean {
    return proto3.util.equals(Volume, a, b);
  }
}

/**
 * @generated from enum apocryph.proto.v0.pod.Volume.Type
 */
export enum Volume_Type {
  /**
   * @generated from enum value: VOLUME_EMPTY = 0;
   */
  VOLUME_EMPTY = 0,

  /**
   * @generated from enum value: VOLUME_FILESYSTEM = 1;
   */
  VOLUME_FILESYSTEM = 1,

  /**
   * @generated from enum value: VOLUME_SECRET = 2;
   */
  VOLUME_SECRET = 2,
}
// Retrieve enum metadata with: proto3.getEnumType(Volume_Type)
proto3.util.setEnumType(Volume_Type, "apocryph.proto.v0.pod.Volume.Type", [
  { no: 0, name: "VOLUME_EMPTY" },
  { no: 1, name: "VOLUME_FILESYSTEM" },
  { no: 2, name: "VOLUME_SECRET" },
]);

/**
 * @generated from enum apocryph.proto.v0.pod.Volume.AccessMode
 */
export enum Volume_AccessMode {
  /**
   * @generated from enum value: VOLUME_RW_ONE = 0;
   */
  VOLUME_RW_ONE = 0,

  /**
   * VOLUME_RO_MANY = 1;
   *
   * @generated from enum value: VOLUME_RW_MANY = 2;
   */
  VOLUME_RW_MANY = 2,
}
// Retrieve enum metadata with: proto3.getEnumType(Volume_AccessMode)
proto3.util.setEnumType(Volume_AccessMode, "apocryph.proto.v0.pod.Volume.AccessMode", [
  { no: 0, name: "VOLUME_RW_ONE" },
  { no: 2, name: "VOLUME_RW_MANY" },
]);

/**
 * @generated from message apocryph.proto.v0.pod.Volume.FilesystemConfig
 */
export class Volume_FilesystemConfig extends Message<Volume_FilesystemConfig> {
  /**
   * "storage"
   *
   * @generated from field: repeated apocryph.proto.v0.pod.Resource resourceRequests = 1;
   */
  resourceRequests: Resource[] = [];

  constructor(data?: PartialMessage<Volume_FilesystemConfig>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.pod.Volume.FilesystemConfig";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "resourceRequests", kind: "message", T: Resource, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Volume_FilesystemConfig {
    return new Volume_FilesystemConfig().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Volume_FilesystemConfig {
    return new Volume_FilesystemConfig().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Volume_FilesystemConfig {
    return new Volume_FilesystemConfig().fromJsonString(jsonString, options);
  }

  static equals(a: Volume_FilesystemConfig | PlainMessage<Volume_FilesystemConfig> | undefined, b: Volume_FilesystemConfig | PlainMessage<Volume_FilesystemConfig> | undefined): boolean {
    return proto3.util.equals(Volume_FilesystemConfig, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.pod.Volume.SecretConfig
 */
export class Volume_SecretConfig extends Message<Volume_SecretConfig> {
  /**
   * @generated from field: bytes cid = 1;
   */
  cid = new Uint8Array(0);

  /**
   * @generated from field: apocryph.proto.v0.pod.Key key = 2;
   */
  key?: Key;

  /**
   * @generated from field: string file = 101;
   */
  file = "";

  /**
   * @generated from field: bytes contents = 102;
   */
  contents = new Uint8Array(0);

  /**
   * @generated from field: string contentsString = 103;
   */
  contentsString = "";

  constructor(data?: PartialMessage<Volume_SecretConfig>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.pod.Volume.SecretConfig";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "cid", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 2, name: "key", kind: "message", T: Key },
    { no: 101, name: "file", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 102, name: "contents", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 103, name: "contentsString", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Volume_SecretConfig {
    return new Volume_SecretConfig().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Volume_SecretConfig {
    return new Volume_SecretConfig().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Volume_SecretConfig {
    return new Volume_SecretConfig().fromJsonString(jsonString, options);
  }

  static equals(a: Volume_SecretConfig | PlainMessage<Volume_SecretConfig> | undefined, b: Volume_SecretConfig | PlainMessage<Volume_SecretConfig> | undefined): boolean {
    return proto3.util.equals(Volume_SecretConfig, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.pod.Replicas
 */
export class Replicas extends Message<Replicas> {
  /**
   * @generated from field: uint32 min = 1;
   */
  min = 0;

  /**
   * @generated from field: uint32 max = 2;
   */
  max = 0;

  /**
   * @generated from field: uint32 targetPendingRequests = 3;
   */
  targetPendingRequests = 0;

  constructor(data?: PartialMessage<Replicas>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.pod.Replicas";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "min", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 2, name: "max", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 3, name: "targetPendingRequests", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Replicas {
    return new Replicas().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Replicas {
    return new Replicas().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Replicas {
    return new Replicas().fromJsonString(jsonString, options);
  }

  static equals(a: Replicas | PlainMessage<Replicas> | undefined, b: Replicas | PlainMessage<Replicas> | undefined): boolean {
    return proto3.util.equals(Replicas, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.pod.Resource
 */
export class Resource extends Message<Resource> {
  /**
   * @generated from field: string resource = 1;
   */
  resource = "";

  /**
   * @generated from oneof apocryph.proto.v0.pod.Resource.quantity
   */
  quantity: {
    /**
     * @generated from field: uint64 amount = 2;
     */
    value: bigint;
    case: "amount";
  } | {
    /**
     * /1000
     *
     * @generated from field: uint64 amountMillis = 3;
     */
    value: bigint;
    case: "amountMillis";
  } | {
    /**
     * *1024
     *
     * @generated from field: uint64 amountKibi = 4;
     */
    value: bigint;
    case: "amountKibi";
  } | {
    /**
     * *1024^2
     *
     * @generated from field: uint64 amountMebi = 5;
     */
    value: bigint;
    case: "amountMebi";
  } | {
    /**
     * *1024^3
     *
     * @generated from field: uint64 amountGibi = 6;
     */
    value: bigint;
    case: "amountGibi";
  } | { case: undefined; value?: undefined } = { case: undefined };

  constructor(data?: PartialMessage<Resource>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.pod.Resource";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "resource", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "amount", kind: "scalar", T: 4 /* ScalarType.UINT64 */, oneof: "quantity" },
    { no: 3, name: "amountMillis", kind: "scalar", T: 4 /* ScalarType.UINT64 */, oneof: "quantity" },
    { no: 4, name: "amountKibi", kind: "scalar", T: 4 /* ScalarType.UINT64 */, oneof: "quantity" },
    { no: 5, name: "amountMebi", kind: "scalar", T: 4 /* ScalarType.UINT64 */, oneof: "quantity" },
    { no: 6, name: "amountGibi", kind: "scalar", T: 4 /* ScalarType.UINT64 */, oneof: "quantity" },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Resource {
    return new Resource().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Resource {
    return new Resource().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Resource {
    return new Resource().fromJsonString(jsonString, options);
  }

  static equals(a: Resource | PlainMessage<Resource> | undefined, b: Resource | PlainMessage<Resource> | undefined): boolean {
    return proto3.util.equals(Resource, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.pod.Key
 */
export class Key extends Message<Key> {
  /**
   * @generated from field: bytes data = 1;
   */
  data = new Uint8Array(0);

  constructor(data?: PartialMessage<Key>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.pod.Key";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "data", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Key {
    return new Key().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Key {
    return new Key().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Key {
    return new Key().fromJsonString(jsonString, options);
  }

  static equals(a: Key | PlainMessage<Key> | undefined, b: Key | PlainMessage<Key> | undefined): boolean {
    return proto3.util.equals(Key, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.pod.KeyPair
 */
export class KeyPair extends Message<KeyPair> {
  /**
   * @generated from field: string privateKey = 2;
   */
  privateKey = "";

  /**
   * @generated from field: string pubAddress = 3;
   */
  pubAddress = "";

  constructor(data?: PartialMessage<KeyPair>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.pod.KeyPair";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 2, name: "privateKey", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "pubAddress", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): KeyPair {
    return new KeyPair().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): KeyPair {
    return new KeyPair().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): KeyPair {
    return new KeyPair().fromJsonString(jsonString, options);
  }

  static equals(a: KeyPair | PlainMessage<KeyPair> | undefined, b: KeyPair | PlainMessage<KeyPair> | undefined): boolean {
    return proto3.util.equals(KeyPair, a, b);
  }
}

/**
 * @generated from message apocryph.proto.v0.pod.VerificationSettings
 */
export class VerificationSettings extends Message<VerificationSettings> {
  /**
   * @generated from field: bool ForcePolicy = 1;
   */
  ForcePolicy = false;

  /**
   * @generated from field: bool PublicVerifiability = 2;
   */
  PublicVerifiability = false;

  /**
   * @generated from field: string VerificationHost = 3;
   */
  VerificationHost = "";

  constructor(data?: PartialMessage<VerificationSettings>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "apocryph.proto.v0.pod.VerificationSettings";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "ForcePolicy", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 2, name: "PublicVerifiability", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
    { no: 3, name: "VerificationHost", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): VerificationSettings {
    return new VerificationSettings().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): VerificationSettings {
    return new VerificationSettings().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): VerificationSettings {
    return new VerificationSettings().fromJsonString(jsonString, options);
  }

  static equals(a: VerificationSettings | PlainMessage<VerificationSettings> | undefined, b: VerificationSettings | PlainMessage<VerificationSettings> | undefined): boolean {
    return proto3.util.equals(VerificationSettings, a, b);
  }
}

