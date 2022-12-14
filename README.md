# PostgreSQL server operator for Kubernetes

## Features

* Creates a database from a CRD
* Creates a role with random username and password from a CRD
* Update a database by creating a role

## Installation

Just create a Kubernetes Secret in the same namespace as operator itself.
Secret should contain these keys: POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASS, POSTGRES_URI_ARGS, POSTGRES_CLOUD_PROVIDER, POSTGRES_DEFAULT_DATABASE.
Example:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: ext-postgres-operator
  namespace: operators
type: Opaque
data:
  POSTGRES_HOST: ZGJwb3N0Z3Jlcwo=
  POSTGRES_USER: ZGJwb3N0Z3Jlcwo=
  POSTGRES_PASS: cG9zdGdyZXM=
  POSTGRES_URI_ARGS: IA==
  POSTGRES_CLOUD_PROVIDER: QVdT
  POSTGRES_DEFAULT_DATABASE: ZGJwb3N0Z3Jlcwo=
```

To install the operator, follow the steps below.

1. Configure Postgres credentials for the operator in `deploy/secret.yaml`
2. Create namespace if needed with\
   `kubectl apply -f deploy/namespace.yaml`
3. Apply the secret with\
   `kubectl apply -f deploy/secret.yaml`
4. Create the operator with either\
    `kubectl kustomize deploy/ | apply -f -`\
    or by using [kustomize](https://github.com/kubernetes-sigs/kustomize) directly\
    `kustomize build deploy/ | apply -f -`

## CRDs

### Postgres

```yaml
apiVersion: db.pgk8s.com/v1alpha1
kind: Postgres
metadata:
  name: my-db
  namespace: app
spec:
  database: test-db # Name of database created in PostgreSQL
  dropOnDelete: false # Set to true if you want the operator to drop the database and role when this CR is deleted (optional)
  masterRole: test-db-group (optional)
  schemas: # List of schemas the operator should create in database (optional)
  - stores
  - customers
  extensions: # List of extensions that should be created in the database (optional)
  - fuzzystrmatch
  - pgcrypto
```

This creates a database called `test-db` and a role `test-db-group` set as the owner of the database, with reader and writer roles also created.

### PostgresUser

```yaml
apiVersion: db.pgk8s.com/v1alpha1
kind: PostgresUser
metadata:
  name: my-db-user
  namespace: app
spec:
  role: username
  database: my-db       # This references the Postgres CR
  secretName: my-secret
  privileges: OWNER     # Can be OWNER/READ/WRITE
  annotations:          # Annotations to be propagated to the secrets metadata section (optional)
    foo: "bar"
```

This creates a user role `username-<hash>` and grants role `test-db-group`, `test-db-writer` or `test-db-reader` depending on `privileges` property. Its credentials are put in secret `my-secret-my-db-user`.

`PostgresUser` needs to reference a `Postgres` in the same namespace.

Two `Postgres` referencing the same database can exist in more than one namespace. The last CR referencing a database will drop the group role and transfer database ownership to the role used by the operator.
Every PostgresUser has a generated Kubernetes secret attached to it, which contains the following data (i.e.):

|  Key                 | Comment             |
|----------------------|---------------------|
| `DATABASE_NAME`      | Name of the database, same as in `Postgres` CR, copied for convenience |
| `HOST`               | PostgreSQL server host |
| `PASSWORD`           | Autogenerated password for user |
| `ROLE`               | Autogenerated role with login enabled (user) |
| `LOGIN`              | Same as `ROLE`. In case `POSTGRES_CLOUD_PROVIDER` is set to "Azure", `LOGIN` it will be set to `{role}@{serverName}`, serverName is extracted from `POSTGRES_USER` from operator's config. |
| `POSTGRES_URL`       | Connection string for Posgres, could be used for Go applications |
| `POSTGRES_JDBC_URL`  | JDBC compatible Postgres URI, formatter as `jdbc:postgresql://{POSTGRES_HOST}/{DATABASE_NAME}` |

#### Branching

`main` branch contains the latest source code with all the features. `vX.X.X` contains code for the specific major versions.
 i.e. `v0.4.x` contains the latest code for 0.4 version of the operator. See compatibility matrix below.

#### Tests

Please write tests and fix any broken tests before you open a PR. Tests should cover at least 80% of your code.

