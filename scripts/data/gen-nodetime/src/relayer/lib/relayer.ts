// cosmosjs related imports.
import {fromHex} from "@cosmjs/encoding";
import {DirectSecp256k1Wallet} from "@cosmjs/proto-signing";
import {GasPrice} from "@cosmjs/stargate";
import {Endpoint, IbcClient, Link} from "@confio/relayer/build";
import {buildCreateClientArgs, IbcClientOptions, prepareConnectionHandshake} from "@confio/relayer/build/lib/ibcclient";
import {orderFromJSON} from "cosmjs-types/ibc/core/channel/v1/channel";

// local imports.
import ConsoleLogger, { LogLevels } from './logger';

type Chain = {
    id: string;
    account: string,
    address_prefix: string;
    rpc_address: string;
    gas_price: string;
    client_id: string;
    estimated_block_time: number;
    estimated_indexer_time: number;
};

type Path = {
    id: string;
    ordering: string;
    src: PathEnd;
    dst: PathEnd;
};

type PathEnd = {
    chain_id: string;
    connection_id: string;
    channel_id: string;
    port_id: string;
    version: string;
    packet_height?: number;
    ack_height?: number;
};

const defaultEstimatedBlockTime = 400;
const defaultEstimatedIndexerTime = 80;

export default class Relayer {
    private defaultMaxAge = 86400;
    private logLevel = 2;

    constructor(logLevel: LogLevels=LogLevels.INFO) {
        if (logLevel) this.logLevel=logLevel;
    }
    public async link([
                          path,
                          srcChain,
                          dstChain,
                          srcKey,
                          dstKey,
                      ]: [Path, Chain, Chain, string, string]): Promise<Path> {
        const srcClient = await Relayer.getIBCClient(srcChain, srcKey);
        const dstClient = await Relayer.getIBCClient(dstChain, dstKey);
        const link = await Relayer.create(srcClient, dstClient, srcChain.client_id, dstChain.client_id, this.logLevel);

        const channels = await link.createChannel(
            'A',
            path.src.port_id,
            path.dst.port_id,
            orderFromJSON(path.ordering),
            path.dst.version
        );

        path.src.channel_id = channels.src.channelId;
        path.dst.channel_id = channels.dest.channelId;
        path.src.connection_id = link.endA.connectionID;
        path.dst.connection_id = link.endB.connectionID;

        return path;
    }
    
    public async start([
                           path,
                           srcChain,
                           dstChain,
                           srcKey,
                           dstKey
                       ]: [Path, Chain, Chain, string, string]): Promise<Path> {
        const srcClient = await Relayer.getIBCClient(srcChain, srcKey);
        const dstClient = await Relayer.getIBCClient(dstChain, dstKey);

        const link = await Link.createWithExistingConnections(
            srcClient,
            dstClient,
            path.src.connection_id,
            path.dst.connection_id,
            new ConsoleLogger(this.logLevel)
        );

        const heights = await link.checkAndRelayPacketsAndAcks(
            {
                packetHeightA: path.src.packet_height,
                packetHeightB: path.dst.packet_height,
                ackHeightA: path.src.ack_height,
                ackHeightB: path.dst.ack_height
            } ?? {},
            2,
            6
        );

        await link.updateClientIfStale('A', this.defaultMaxAge);
        await link.updateClientIfStale('B', this.defaultMaxAge);

        path.src.packet_height = heights.packetHeightA;
        path.dst.packet_height = heights.packetHeightB;
        path.src.ack_height = heights.ackHeightA;
        path.dst.ack_height = heights.ackHeightB;

        return path;
    }

    private static async getIBCClient(chain: Chain, key: string): Promise<IbcClient> {
        const chainGP = GasPrice.fromString(chain.gas_price);
        const signer = await DirectSecp256k1Wallet.fromKey(fromHex(key), chain.address_prefix);

        const [account] = await signer.getAccounts();
        const options: IbcClientOptions = {
            gasPrice: chainGP,
            estimatedBlockTime: chain.estimated_block_time ?? defaultEstimatedBlockTime,
            estimatedIndexerTime: chain.estimated_indexer_time ?? defaultEstimatedIndexerTime
        }

        return await IbcClient.connectWithSigner(
            chain.rpc_address,
            signer,
            account.address,
            options
        );
    }

    private static async create(
        nodeA: IbcClient,
        nodeB: IbcClient,
        clientA: string,
        clientB: string,
        logLevel:number
    ): Promise<Link> {
        let dstClientID = clientB;
        if (!clientB) {
            const args = await buildCreateClientArgs(nodeA);
            const {clientId: clientId} = await nodeB.createTendermintClient(
                args.clientState,
                args.consensusState
            );
            dstClientID = clientId;
        }

        let srcClientID = clientA;
        if (!clientA) {
            // client on A pointing to B
            const args2 = await buildCreateClientArgs(nodeB);
            const {clientId: clientId} = await nodeA.createTendermintClient(
                args2.clientState,
                args2.consensusState
            );
            srcClientID = clientId;
        }

        // wait a block to ensure we have proper proofs for creating a connection (this has failed on CI before)
        await Promise.all([nodeA.waitOneBlock(), nodeB.waitOneBlock()]);

        // connectionInit on nodeA
        const {connectionId: connIdA} = await nodeA.connOpenInit(
            srcClientID,
            dstClientID
        );

        // connectionTry on nodeB
        const proof = await prepareConnectionHandshake(
            nodeA,
            nodeB,
            srcClientID,
            dstClientID,
            connIdA
        );

        const {connectionId: connIdB} = await nodeB.connOpenTry(dstClientID, proof);

        // connectionAck on nodeA
        const proofAck = await prepareConnectionHandshake(
            nodeB,
            nodeA,
            dstClientID,
            srcClientID,
            connIdB
        );
        await nodeA.connOpenAck(connIdA, proofAck);

        // connectionConfirm on dest
        const proofConfirm = await prepareConnectionHandshake(
            nodeA,
            nodeB,
            srcClientID,
            dstClientID,
            connIdA
        );
        await nodeB.connOpenConfirm(connIdB, proofConfirm);

        const endA = new Endpoint(nodeA, srcClientID, connIdA);
        const endB = new Endpoint(nodeB, dstClientID, connIdB);

        return new Link(endA, endB, new ConsoleLogger(logLevel));
    }
}
