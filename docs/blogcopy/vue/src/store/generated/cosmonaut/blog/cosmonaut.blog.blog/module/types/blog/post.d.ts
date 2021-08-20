import { Writer, Reader } from 'protobufjs/minimal';
export declare const protobufPackage = "cosmonaut.blog.blog";
export interface Post {
    creator: string;
    id: number;
    title: string;
    body: string;
}
export declare const Post: {
    encode(message: Post, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Post;
    fromJSON(object: any): Post;
    toJSON(message: Post): unknown;
    fromPartial(object: DeepPartial<Post>): Post;
};
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
