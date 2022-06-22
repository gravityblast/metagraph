# Metagraph

Metagraph watches Ethereum smart contracts events emitting metadata pointers to IPFS files and pins them on an IPFS node.

## Summary

Many dapps stores files in decentralized storage systems and save a "metadata pointer" pointing to them in a smart contract.

Those files are often stored on IPFS and need to be pinned to avoid them to be garbage collected.

In most cases these dapps allow any user to store a new file on IPFS and save its metapointer on a smart contract,
and there's no easy way to pin these files directly.

One way could be to allow the dapp to call a web service asking to pin a specific file, but there are some arguments
against it:

* the dapp needs to point to a centralized server, making it basically a centralized app.
* the server needs to expose an API that anyone can spam even if an API key is used, since the API key would be in
written clear in the client side dapp.

The other way we propose with Metagraph is to allow anyone to run a `metagraph` node that watches
smart contracts events emitting a metadata pointer, and pin those files to a local or remote SAAS managed node,
making the pinning process decentralized and permissionless.

## Build dependencies

* go 1.18

## Build

```
git clone git@github.com:gitcoinco/metagraph.git
cd metagraph
make build
```

## Configuration

```
mkdir -p abis
cp config.json.example config.json
```

The `config.json` file has some global configurations and one event configuration for each event emitting a Metadata Pointer.

Example:

```
{
  "providerURL": "wss://eth-goerli.alchemyapi.io/v2/API_KEY",
  "events": [
    {
      "id": "project-metadata-updated",
      "description": "Project Metadata Updated With MetaPtr Struct",
      "contractName": "ProjectRegistry",
      "contractAddress": "0xa449048290cf2c68c387ecf32847C59D9746C438",
      "eventName": "MetadataUpdated",
      "metaPointerField": "metaPtr",
      "protocolField": "Protocol",
      "pointerField": "Pointer",
      "startingBlock": 7065923
    },
    ...
  ]
}
```

Each event configuration describes how the Metadata Pointer is specified in an event emitted by a smart contract.

The `contractName` field specifies the ABI file name stored in `abis/CONTRACT_NAME.json`.

In particular, the Metadata Pointer can be specified in 3 ways:

### With a Metadata Pointer struct

A dapp can emit an event with a MetaPtr struct that describe the Metadata Pointer with two fields

* protocol (1 for IPFS, etc...)
* pointer (a CIP for IPFS, or another file pointer for other storage systems)

```
struct MetaPtr {
    /// More info at https://github.com/gitcoinco/grants-round/tree/main/packages/contracts/docs/MetaPtrProtocol.md
    uint256 protocol;

    /// @notice Pointer to fetch metadata for the specified protocol
    string pointer;
}

contract ProjectRegistry {
    event MetadataUpdated(uint96 indexed projectID, MetaPtr metaPtr);

    function createProject(MetaPtr memory metadata) external {
      ...
      emit MetadataUpdated(projectID, metadata);
    }
}
```

In this case the event configuration requires three fields:

```
{
  ...
  "metaPointerField": "metaPtr",
  "protocolField": "Protocol",
  "pointerField": "Pointer",
  ...
}
```

## With the Protocol and Pointer fields outside of a struct

```
contract ProjectRegistry {
    event MetadataUpdated(uint96 indexed projectID, uint256 metadataProtocol, string metadataPointer);

    function createProject(uint256 metadataProtocol, string memory metadataPointer) external {
      ...
      emit MetadataUpdated(projectID, metadataProtocol, metadataPointer);
    }
}
```

config:

```
{
  ...
  // metaPointerField not defined
  "protocolField": "metadataProtocol",
  "pointerField": "metadataPointer",
  ...
}
```

## Only with a Pointer without specifing a protocol

Depending on the implementation the protocol field can default to
a specific storage system like IPFS.

```
contract ProjectRegistry {
    event MetadataUpdated(uint96 indexed projectID, string metadata);

    function createProject(string metadata) external {
      ...
      emit MetadataUpdated(projectID, metadata);
    }
}
```

config:

```
{
  ...
  // metaPointerField not defined
  // protocolField not defined
  "pointerField": "metadata",
  ...
}
```

## Run

```
make run
```

## Todo

- [X] Parse global configuration
- [X] Parse configuration for each event
- [X] Filter logs for each event/contract
- [X] Filter logs from a specific block
- [X] Parse event with MetaPointer struct
- [X] Parse event with metadata protocol and metadata pointer
- [X] Parse event with only metadata pointer
- [X] Log all Metadata Pointer found
- [ ] Call an external ipfs node/service to pin the file
- [ ] Watch for future events instead of only filtering for past events
- [ ] Save the last watched block for each event configuration
- [ ] Create a "fork event configuration" to allow to watch contracts created from a factory
