# Wordpress

This example demonstrates how to use Zarf to deploy a Wordpress blog into a cluster.  It is used as a part of the [Creating a Zarf Package](../../docs/6-zarf-tutorials/0-creating-a-zarf-package.md) and [Deploying a Zarf Package](../../docs/6-zarf-tutorials/2-deploying-zarf-packages.md) tutorials.

:::info

To view the example source code, select the `Edit this page` link below the article and select the parent folder.

:::

```yaml
kind: ZarfPackageConfig # ZarfPackageConfig is the package kind for most normal zarf packages
metadata:
  name: wordpress       # specifies the name of our package and should be unique and unchanging through updates
  version: 16.0.4       # (optional) a version we can track as we release updates or publish to a registry
  description: |        # (optional) a human-readable description of the package that you are creating
    "A Zarf Package that deploys the Wordpress blogging and content management platform"

variables:
    # The unique name of the variable corresponding to the ###ZARF_VAR_### template
  - name: WORDPRESS_USERNAME
    # A human-readable description of the variable shown during prompting
    description: The username that is used to login to the Wordpress admin account
    # A default value to take if --confirm is used or the user chooses the default prompt
    default: zarf
    # Whether to prompt for this value interactively if it is not --set on the CLI
    prompt: true
  - name: WORDPRESS_PASSWORD
    description: The password that is used to login to the Wordpress admin account
    prompt: true
    # Whether to treat this value as sensitive to keep it out of Zarf logs
    sensitive: true
  - name: WORDPRESS_EMAIL
    description: The email that is used for the Wordpress admin account
    default: hello@defenseunicorns.com
    prompt: true
  - name: WORDPRESS_FIRST_NAME
    description: The first name that is used for the Wordpress admin account
    default: Zarf
    prompt: true
  - name: WORDPRESS_LAST_NAME
    description: The last name that is used for the Wordpress admin account
    default: The Axolotl
    prompt: true
  - name: WORDPRESS_BLOG_NAME
    description: The blog name that is used for the Wordpress admin account
    default: The Zarf Blog
    prompt: true

components:
  - name: wordpress  # specifies the name of our component and should be unique and unchanging through updates
    description: |   # (optional) a human-readable description of the component you are defining
      "Deploys the Bitnami-packaged Wordpress chart into the cluster"
    required: true   # (optional) sets the component as 'required' so that it is always deployed
    charts:
      - name: wordpress
        url: oci://registry-1.docker.io/bitnamicharts/wordpress
        version: 16.0.4
        namespace: wordpress
        valuesFiles:
          - wordpress-values.yaml
    images:
      - docker.io/bitnami/apache-exporter:0.13.3-debian-11-r2
      - docker.io/bitnami/mariadb:10.11.2-debian-11-r21
      - docker.io/bitnami/wordpress:6.2.0-debian-11-r18
    manifests:
      - name: connect-services
        files:
          - connect-services.yaml
```
