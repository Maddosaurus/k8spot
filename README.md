# k8spot

## Usage
`go run main.go`

Then point your kubectl to the supplied config, e.g.:
`kubectl --kubeconfig=kubeconfig get nodes`

To pull corresponding API JSONs for your kubernetes:
`curl --cacert /your/ca.crt --cert ~/your/client.crt --key /your/client.key https://your.k8s.fqdn:8443/...`
... or run `get_json.sh`
