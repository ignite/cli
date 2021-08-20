/* eslint-disable */
import { Reader, Writer } from 'protobufjs/minimal'

export const protobufPackage = 'cosmonaut.blog.blog'

/** this line is used by starport scaffolding # proto/tx/message */
export interface MsgCreatePost {
  creator: string
  title: string
  body: string
}

export interface MsgCreatePostResponse {}

const baseMsgCreatePost: object = { creator: '', title: '', body: '' }

export const MsgCreatePost = {
  encode(message: MsgCreatePost, writer: Writer = Writer.create()): Writer {
    if (message.creator !== '') {
      writer.uint32(10).string(message.creator)
    }
    if (message.title !== '') {
      writer.uint32(18).string(message.title)
    }
    if (message.body !== '') {
      writer.uint32(26).string(message.body)
    }
    return writer
  },

  decode(input: Reader | Uint8Array, length?: number): MsgCreatePost {
    const reader = input instanceof Uint8Array ? new Reader(input) : input
    let end = length === undefined ? reader.len : reader.pos + length
    const message = { ...baseMsgCreatePost } as MsgCreatePost
    while (reader.pos < end) {
      const tag = reader.uint32()
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string()
          break
        case 2:
          message.title = reader.string()
          break
        case 3:
          message.body = reader.string()
          break
        default:
          reader.skipType(tag & 7)
          break
      }
    }
    return message
  },

  fromJSON(object: any): MsgCreatePost {
    const message = { ...baseMsgCreatePost } as MsgCreatePost
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator)
    } else {
      message.creator = ''
    }
    if (object.title !== undefined && object.title !== null) {
      message.title = String(object.title)
    } else {
      message.title = ''
    }
    if (object.body !== undefined && object.body !== null) {
      message.body = String(object.body)
    } else {
      message.body = ''
    }
    return message
  },

  toJSON(message: MsgCreatePost): unknown {
    const obj: any = {}
    message.creator !== undefined && (obj.creator = message.creator)
    message.title !== undefined && (obj.title = message.title)
    message.body !== undefined && (obj.body = message.body)
    return obj
  },

  fromPartial(object: DeepPartial<MsgCreatePost>): MsgCreatePost {
    const message = { ...baseMsgCreatePost } as MsgCreatePost
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator
    } else {
      message.creator = ''
    }
    if (object.title !== undefined && object.title !== null) {
      message.title = object.title
    } else {
      message.title = ''
    }
    if (object.body !== undefined && object.body !== null) {
      message.body = object.body
    } else {
      message.body = ''
    }
    return message
  }
}

const baseMsgCreatePostResponse: object = {}

export const MsgCreatePostResponse = {
  encode(_: MsgCreatePostResponse, writer: Writer = Writer.create()): Writer {
    return writer
  },

  decode(input: Reader | Uint8Array, length?: number): MsgCreatePostResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input
    let end = length === undefined ? reader.len : reader.pos + length
    const message = { ...baseMsgCreatePostResponse } as MsgCreatePostResponse
    while (reader.pos < end) {
      const tag = reader.uint32()
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7)
          break
      }
    }
    return message
  },

  fromJSON(_: any): MsgCreatePostResponse {
    const message = { ...baseMsgCreatePostResponse } as MsgCreatePostResponse
    return message
  },

  toJSON(_: MsgCreatePostResponse): unknown {
    const obj: any = {}
    return obj
  },

  fromPartial(_: DeepPartial<MsgCreatePostResponse>): MsgCreatePostResponse {
    const message = { ...baseMsgCreatePostResponse } as MsgCreatePostResponse
    return message
  }
}

/** Msg defines the Msg service. */
export interface Msg {
  /** this line is used by starport scaffolding # proto/tx/rpc */
  CreatePost(request: MsgCreatePost): Promise<MsgCreatePostResponse>
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc
  constructor(rpc: Rpc) {
    this.rpc = rpc
  }
  CreatePost(request: MsgCreatePost): Promise<MsgCreatePostResponse> {
    const data = MsgCreatePost.encode(request).finish()
    const promise = this.rpc.request('cosmonaut.blog.blog.Msg', 'CreatePost', data)
    return promise.then((data) => MsgCreatePostResponse.decode(new Reader(data)))
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>
}

type Builtin = Date | Function | Uint8Array | string | number | undefined
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>
