import Relayer from "../lib/relayer";

const relayer = new Relayer();

relayer.link(
    [
        {
            "id": "spn-1-orbit-1",
            "ordering": "ORDER_UNORDERED",
            "src": {
                "chain_id": "spn-1",
                "connection_id": "",
                "channel_id": "",
                "port_id": "monitoringp",
                "version": "spn-1",
                "packet_height": 0,
                "ack_height": 0
            },
            "dst": {
                "chain_id": "orbit-1",
                "connection_id": "",
                "channel_id": "",
                "port_id": "monitoringc",
                "version": "orbit-1",
                "packet_height": 0,
                "ack_height": 0
            }
        },
        {
            "id": "spn-1",
            "account": "alice",
            "address_prefix": "spn",
            "rpc_address": "http://0.0.0.0:26657",
            "gas_price": "0.0000025uspn",
            "gas_limit": 400000
        },
        {
            "id": "orbit-1",
            "account": "alice",
            "address_prefix": "cosmos",
            "rpc_address": "http://localhost:26659",
            "gas_price": "0.0000025stake",
            "gas_limit": 400000
        },
        "9347d4a8b6ab5251b4d30e1751d526713143049a7c8a07e93ff9a6f9fb604312",
        "9347d4a8b6ab5251b4d30e1751d526713143049a7c8a07e93ff9a6f9fb604312",
        "07-tendermint-0",
        "07-tendermint-0"
    ]
)
