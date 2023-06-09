# Byzantine Shield

Byzantine Shield is a middleware proxy for blockchain clients designed to mitigate the risks associated with communicating with byzantine nodes. It forwards JSON RPC requests from clients, like Geth, to a pre-configured list of blockchain nodes, and aggregates their responses to detect inconsistencies and return a single, consistent answer.

## Features

- Aggregates responses from multiple blockchain nodes
- Detects inconsistencies between node responses
- Reduces the risks associated with byzantine nodes
- Compatible with blockchains following the JSON RPC API format
- Written in Golang for performance and concurrency

## Security

Byzantine Shield enhances the security of blockchain clients by:

1. **Mitigating byzantine node risks:** By aggregating responses from multiple nodes and detecting inconsistencies, Byzantine Shield ensures that clients receive accurate and consistent data, even in the presence of malicious nodes.

2. **Reducing single point of failure:** Clients are no longer dependent on a single blockchain node for their data. Byzantine Shield contacts multiple nodes, increasing the chances of receiving correct data.

3. **Enhancing trust:** By cross-verifying the responses from multiple nodes, Byzantine Shield helps build trust in the data received by the client, ensuring its validity and accuracy.

## Usage

byzantine-shield can be easily integrated into your existing blockchain client application. Follow the steps below to set up and configure byzantine-shield:

1. **Build the tool:** Build the byzantine-shield binary by running the following command:

    ```bash
    go build -o byzantine-shield cmd/main.go
    ```

    This will create an executable named `byzantine-shield`.

2. **Configuration:** Create a configuration file for byzantine-shield using the example provided in the `example/` folder. This file should include the list of blockchain nodes you want to use, as well as any other relevant settings.

3. **Starting byzantine-shield:** Start the byzantine-shield proxy with the following command:

    ```bash
    ./byzantine-shield --config path/to/your/config.yml --addr 127.0.0.1 --port 8080
    ```

    Replace `path/to/your/config.yml` with the path to your configuration file, and adjust the `--addr` and `--port` options as needed.

4. **Connecting your client application:** Now, your client application can connect to byzantine-shield by sending JSON RPC requests to the specified address and port (e.g., `http://127.0.0.1:8080`).

By following these steps, you can easily integrate byzantine-shield into your existing blockchain client application and benefit from its security features.

## Contributing

We welcome contributions to Byzantine Shield! If you would like to contribute, please follow our [contributing guidelines](CONTRIBUTING.md) to get started.

## License

Byzantine Shield is licensed under the [MIT License](LICENSE)
