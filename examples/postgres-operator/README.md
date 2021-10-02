# Zarf Postgres Operator Example

This example demonstrates deploying a performant and highly available PostgreSQL database to a Zarf airgap cluster. It uses Zalando's [postgres-operator](https://github.com/zalando/postgres-operator) and provides the Postgres Operator UI and a deployment of PGAdmin for demo purposes.

## Tool Choice

After looking at several alternatives, Zalando's postgres operator felt like the best choice. Other tools that were close runners-up were the postgres-operator by [CrunchyData](https://github.com/CrunchyData/postgres-operator) and [KubeDB](https://github.com/kubedb/operator).

## Instructions

1. Create a Zarf cluster as outlined in the main [README](../../README.md#2-create-the-zarf-cluster)
2. Follow [step 3](../../README.md#3-add-resources-to-the-zarf-cluster) using the `zarf.yaml` config in this folder
3. The Postgres Operator UI will be available at [https://postgres-operator-ui.localhost:8443](https://postgres-operator-ui.localhost:8443) and PGAdmin will be available at [https://pgadmin.localhost:8443](https://pgadmin.localhost:8443).
4. Set up a server in PGAdmin:
    - General // Name: `acid-zarf-test`
    - General // Server group: `Servers`
    - Connection // Host: (the URL in the table below)
    - Connection // Port: `5432`
    - Connection // Maintenance database: `postgres`
    - Connection // Username: `zarf`
    - Connection // Password: (run the command in the table below)
    - SSL // SSL mode: `Require`
5. Create the backups bucket in MinIO (TODO: Figure out how to create the bucket automatically)
   1. Navigate to [https://minio-console.localhost:8443](https://minio-console.localhost:8443)
   2. Log in - Username: `minio` - Password: `minio123`
   3. Buckets -> Create Bucket
      - Bucket Name: `postgres-operator-backups`

## Logins

| Service                   | URL                                                                                        | Username             | Password                                                                                                                                                   |
| ------------------------- | ------------------------------------------------------------------------------------------ | -------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Postgres Operator UI      | [https://postgres-operator-ui.localhost:8443](https://postgres-operator-ui.localhost:8443) | N/A                  | N/A                                                                                                                                                        |
| PGAdmin                   | [https://pgadmin.localhost:8443](https://pgadmin.localhost:8443)                           | `zarf@example.local` | Run: `zarf tools get-admin-password`                                                                                                                       |
| Example Postgres Database | `acid-zarf-test.postgres-operator.svc.cluster.local`                                       | `zarf`               | Run: `echo $(kubectl get secret zarf.acid-zarf-test.credentials.postgresql.acid.zalan.do -n postgres-operator --template={{.data.password}} \| base64 -d)` |
| Minio Console             | [https://minio-console.localhost:8443](https://minio-console.localhost:8443)               | `minio`              | `minio123`                                                                                                                                                 |

## References
- https://blog.flant.com/comparing-kubernetes-operators-for-postgresql/
- https://blog.flant.com/our-experience-with-postgres-operator-for-kubernetes-by-zalando/