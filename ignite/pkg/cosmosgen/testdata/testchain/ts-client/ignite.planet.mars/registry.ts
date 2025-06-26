import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgMyMessageRequest } from "./types/ignite/planet/mars/mars";
import { MsgBarRequest } from "./types/ignite/planet/mars/mars";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/ignite.planet.mars.MsgMyMessageRequest", MsgMyMessageRequest],
    ["/ignite.planet.mars.MsgBarRequest", MsgBarRequest],
    
];

export { msgTypes }