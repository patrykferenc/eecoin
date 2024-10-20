# eecoin

EECoin is a decentralised cryptocurrency, made with love on the Faculty of Electrical Engineering at the WUT.

## Design

The design is stored mainly in the `docs/` folder as a VPP project.

_Note: make sure to check the diagrams using diagram navigator in case you can't find the nested models in the project model explorer window_

Please export the project from time to time to avoid corrupting the vpp file.

## Development

to run the **node** application run

```bash
go run cmd/node
```

to run the **wallet** application run

```bash
go run cmd/wallet
```

### Testing

To test the application you can use `go test`.

```bash
go test ./...
```
