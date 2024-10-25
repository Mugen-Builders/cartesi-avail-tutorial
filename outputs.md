# Inspecting and reading outputs

## Inspecting state

Inspecting the state of your dApp though `handle_inspect` function is done in the same way as using Cartesi Rollups standalone. You can refer to the [docs](https://docs.cartesi.io/cartesi-rollups/1.5/development/send-requests/#make-inspect-calls)

## Querying outputs

Querying outputs directly is the exact same as using Cartesi Rollups standalone. You can refer to the [docs](https://docs.cartesi.io/cartesi-rollups/1.5/rollups-apis/graphql/overview/)

To query outputs from a specific the process is very similar to using Cartesi Rollups Standalone. You can refer to the [docs](https://docs.cartesi.io/cartesi-rollups/1.5/rollups-apis/graphql/overview/) to read more.
The big difference is the output format. Instead of querying inputs through the `index` field, you query them through an `id` field.

This id field can come in two ways:

- It is a hex value returned from `/submit` endpoint when the input comes from and EIP-712 signed message
- It is string containing a scalar integer value that can be found inside the events emitted by the `inputBox` contract when sending the transaction through the layer 1.

### Example Queries

#### Listing inputs

```graphql
query {
  inputs(first: 30) {
    edges {
      node {
        id
        index
        status
        blockTimestamp
        msgSender
        payload
      }
    }
  }
}
```

#### Getting a specific input through its `id`

```graphql
query {
  input(id: "<input-id>") {
    id
    index
    status
    blockTimestamp
    msgSender
    payload
    notices {
      edges {
        node {
          payload
        }
      }
    }
    reports {
      edges {
        node {
          payload
        }
      }
    }
    vouchers {
      edges {
        node {
          payload
        }
      }
    }
  }
}
```
