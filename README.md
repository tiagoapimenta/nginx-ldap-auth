# Nginx LDAP Auth

Use this in order to provide a ingress authentication over LDAP for Kubernetes, change the Secret inside `config.sample.yaml` to match your LDAP server and run:

    kubectl create secret generic nginx-ldap-auth --from-file=config.yaml=config.sample.yaml

    kubectl apply -f k8s.yaml

For RBAC enabled cluster use the k8s-rbac.yaml manifest instead:

    kubectl apply -f k8s-rbac.yaml

Configure your ingress with annotation `nginx.ingress.kubernetes.io/auth-url: http://nginx-ldap-auth.default.svc.cluster.local:5555` as described on [nginx documentation](https://kubernetes.github.io/ingress-nginx/examples/auth/external-auth/).

## Config

The actual version choose a random server, in future version it is intended to have a pool of them, that is why it is a list, not a single one, but you can fill only one if you wish.

The prefix tell the program which protocol to use, if `ldaps://` it will try LDAP over SSL, if `ldap://` it will try plain LDAP with STARTTLS, case no prefix is given it will try to guess based on port, 636 for SSL and 389 for plain.

The actual version will fail if neither SSL or STARTTLS is possible, but next version will allow plain LDAP.

If the `user.requiredGroups` list is omited or empty all LDAP users will be allowed regardless the group, if not empty all groups will be required, the next version will have more flexible configuration.

If you are not sure what `filter`, `bindDN` or `baseDN` to use, here is a tip:

    ldapsearch -H ${servers[*]} -D ${auth.bindDN} -w ${auth.bindPW} -b ${user.baseDN|group.baseDN} ${user.filter|group.filter}

Replace the values between `${...}` to the ones on `config.yaml`, when you succeed you can fill the final configuration.

Timeouts are configurable, but it is recommended not to use values less than some seconds, it was planned to prevent several identical requests to LDAP servers.
