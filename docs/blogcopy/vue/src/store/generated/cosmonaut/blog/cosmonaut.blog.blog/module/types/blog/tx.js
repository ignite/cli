/* eslint-disable */
import { Reader, Writer } from 'protobufjs/minimal';
export const protobufPackage = 'cosmonaut.blog.blog';
const baseMsgCreatePost = { creator: '', title: '', body: '' };
export const MsgCreatePost = {
    encode(message, writer = Writer.create()) {
        if (message.creator !== '') {
            writer.uint32(10).string(message.creator);
        }
        if (message.title !== '') {
            writer.uint32(18).string(message.title);
        }
        if (message.body !== '') {
            writer.uint32(26).string(message.body);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMsgCreatePost };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.creator = reader.string();
                    break;
                case 2:
                    message.title = reader.string();
                    break;
                case 3:
                    message.body = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseMsgCreatePost };
        if (object.creator !== undefined && object.creator !== null) {
            message.creator = String(object.creator);
        }
        else {
            message.creator = '';
        }
        if (object.title !== undefined && object.title !== null) {
            message.title = String(object.title);
        }
        else {
            message.title = '';
        }
        if (object.body !== undefined && object.body !== null) {
            message.body = String(object.body);
        }
        else {
            message.body = '';
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.creator !== undefined && (obj.creator = message.creator);
        message.title !== undefined && (obj.title = message.title);
        message.body !== undefined && (obj.body = message.body);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseMsgCreatePost };
        if (object.creator !== undefined && object.creator !== null) {
            message.creator = object.creator;
        }
        else {
            message.creator = '';
        }
        if (object.title !== undefined && object.title !== null) {
            message.title = object.title;
        }
        else {
            message.title = '';
        }
        if (object.body !== undefined && object.body !== null) {
            message.body = object.body;
        }
        else {
            message.body = '';
        }
        return message;
    }
};
const baseMsgCreatePostResponse = {};
export const MsgCreatePostResponse = {
    encode(_, writer = Writer.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMsgCreatePostResponse };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(_) {
        const message = { ...baseMsgCreatePostResponse };
        return message;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = { ...baseMsgCreatePostResponse };
        return message;
    }
};
export class MsgClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
    }
    CreatePost(request) {
        const data = MsgCreatePost.encode(request).finish();
        const promise = this.rpc.request('cosmonaut.blog.blog.Msg', 'CreatePost', data);
        return promise.then((data) => MsgCreatePostResponse.decode(new Reader(data)));
    }
}
