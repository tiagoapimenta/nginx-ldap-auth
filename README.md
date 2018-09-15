# Nginx LDAP Auth

Use this in order to provide a ingress authentication over LDAP for Kubernetes, change the ConfigMap inside `k8s.yaml` to match your LDAP server and run:

    kubectl apply -f k8s.yaml

Configure your ingress with annotation `nginx.ingress.kubernetes.io/auth-url: http://nginx-ldap-auth:5555` as described on [nginx documentation](https://kubernetes.github.io/ingress-nginx/examples/auth/external-auth/):
