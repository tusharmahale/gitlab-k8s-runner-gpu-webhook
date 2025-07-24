# Gitlab K8s Runner Container Admission Webhook with Kubernetes Mutation

This repository contains a Kubernetes admission webhook designed to mutate pod specifications by adding GPU resources to any container named "build". The webhook operates as part of a GitLab CI/CD pipeline and is intended to ensure that specified containers within your Kubernetes cluster are equipped with GPU resource limits and requests.


## Key Features

- **Automated Resource Mutation**: Automatically adds GPU resources (requests and limits) to designated containers to facilitate workloads requiring GPU access.

- **CI/CD Integration**: Seamlessly integrates with GitLab CI/CD for continuous delivery to your Kubernetes environment.

## Deployment

Before deploying the webhook, ensure you have TLS certificates configured for secure communication, as Kubernetes requires HTTPS for webhooks.

1. **Generate Certificates**:
   Run the certificate script in the `deploy/` directory to create the necessary certificates.

    ```bash
    cd deploy
    ./certificate.sh
    ```

2. **Deploy to Kubernetes**:
   Apply the Kubernetes manifests in the correct order to deploy the service, create necessary deployment configurations, and register the admission webhook.

## Usage

After deployment, any new pod creation request to your Kubernetes API server will invoke the admission webhook. The webhook evaluates each container in a pod and modifies those with the name "build" to include `nvidia.com/gpu` resource limits and requests.

## License

This project is licensed under the MIT License.

## Contributing

We welcome contributions! Please submit pull requests or open issues for any changes or suggestions.