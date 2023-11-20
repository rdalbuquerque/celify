# celify
CLI to run CEL based validations agaisnt yaml or json

## Get started
### Install on Linux
```bash
# Downloads the CLI based on your OS/arch and puts it in /usr/local/bin
curl -fsSL https://raw.githubusercontent.com/rdalbuquerque/celify/master/scripts/install.sh | sh
```

### Install on Windows
```powershell
Invoke-RestMethod "https://raw.githubusercontent.com/rdalbuquerque/celify/master/scripts/install.ps1" | Invoke-Expression
```

### Example usage
```bash
validations=$(cat <<EOF
validations:
- expression: "object.spec.template.spec.containers.all(c, c.resources.limits.memory != null && c.resources.requests.memory != null)"
  messageExpression: "'all containers must specify memory resource'"
- expression: "object.metadata.name == 'my-deployment'"
  messageExpression: "'expected deployment name to be different than my-deployment, got ' + object.metadata.name"
EOF
)
target=$(cat <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - name: my-container
        image: nginx
        resources:
          limits:
            memory: 1Gi
          requests:
            memory: 1Gi
EOF
)
celify validate --validations "$validations" --target "$target"
```
Output:
![Alt text](image.png)

In the case a validation fails, the output will show the message expression result and the evaluated object.

In this example we are validating a Kubernetes deployment against a set of rules. The first rule checks if all containers have memory limits and requests defined. The second rule checks if the deployment name is different than `my-deployment`.