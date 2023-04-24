# byzantine-shield

byzantine-shield is a middleware proxy for blockchain clients designed to mitigate the risks associated with communicating with byzantine nodes. It forwards JSON RPC requests from clients, like Geth, to a pre-configured list of blockchain nodes, and aggregates their responses to detect inconsistencies and return a single, consistent answer.

## Features

- Aggregates responses from multiple blockchain nodes
- Detects inconsistencies between node responses
- Reduces the risks associated with byzantine nodes
- Follows the JSON RPC API format
- Written in Golang for performance and concurrency

## Security

byzantine-shield enhances the security of blockchain clients by:

1. **Mitigating byzantine node risks:** By aggregating responses from multiple nodes and detecting inconsistencies, byzantine-shield ensures that clients receive accurate and consistent data, even in the presence of malicious nodes.

2. **Reducing single point of failure:** Clients are no longer dependent on a single blockchain node for their data. byzantine-shield contacts at least F+1 nodes (where F is the number of byzantine nodes the blockchain tolerates), increasing the chances of receiving correct data.

3. **Enhancing trust:** By cross-verifying the responses from multiple nodes, byzantine-shield helps build trust in the data received by the client, ensuring its validity and accuracy.

## Usage

byzantine-shield can be easily integrated into your existing blockchain client application. Detailed instructions on setting up and configuring byzantine-shield will be provided in the documentation.

## Roadmap

- MVP
- Utilize Goroutines for concurrent requests to nodes
- Enable dynamic configuration for adding, removing, or updating nodes
- Improve error handling with retries and timeouts
- Integrate logging and monitoring for performance tracking and troubleshooting
- Add authentication and encryption for secure communication

## Contributing

We welcome contributions to byzantine-shield! If you would like to contribute, please follow our [contributing guidelines](CONTRIBUTING.md) to get started.

## License

byzantine-shield is licensed under the [MIT License](LICENSE)
