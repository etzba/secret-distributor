# secret-distributor

A kubernetes operator for managing secrets inside a cluster. No password or secret management tool that requires authentication and external database, just an operator that spreads secrets inside a cluster.

### Why to use this operator

Unlike secret management tools, this operator helps setting up secrets in a namespace without external HTTP access to another api server, no database to maintain and no need to set up configuration.

### How does it work

Setup all secrets to distribute in the same namespace you deployed this operator. The operator will watch for CRD in another namespace:
```yaml
apiVersion: secdist.etzba.com/v1
kind: Distribution
metadata:
  labels:
    app.kubernetes.io/name: secret-distributor
    app.kubernetes.io/managed-by: kustomize
  name: docker-config-distribution
spec:
  secretName: docker-config
```
And set it up in the target namespace.

In this case for example, you are solving a problem of having no way to set docker config secrets in a new namespace.

#### Which problems it is solving

##### Automation of secret distribution

Assuming you'd like to run a test environment automatically and set your database or redis secrets automatically from a central location in the cluster - just add a crd with the relevant secret to setup.

##### Distribution of docker config json secrets

It is very easy to create a new docker secret to pull images from repository by adding a CRD to the new namespace
