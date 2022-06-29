# k8spot
WARNING: THIS IS NOT PRODUCTION READY!

## Talk
- k8spot has been featured at Troopers 22 (see troopers-22.pdf for the slides)

## To Dos / Areas of Improvements
- Add client cert fingerprint logging
- Secure json paths (gate to third_party dir)
- Spoof x-kubernetes headers dynamically
- HTTPS (w/ real certs)
- More endpoints
- Autogen for workloads
    - Try parsing OpenAPI
- Add more k8s flavors


## Usage
`go run main.go`

Then point your kubectl to the supplied config, e.g.:
`kubectl --kubeconfig=kubeconfig get nodes`

To pull corresponding API JSONs for your kubernetes:
`curl --cacert /your/ca.crt --cert ~/your/client.crt --key /your/client.key https://your.k8s.fqdn:8443/...`
... or run `get_json.sh`
