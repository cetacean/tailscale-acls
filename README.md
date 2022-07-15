# `cetacean` Tailscale ACLs

These are the live Tailscale ACLs for [cetacean](https://github.com/cetacean).
This is automatically deployed using
[tailscale/gitops-acl-action](https://github.com/tailscale/gitops-acl-action).
The API key should be updated monthly.

This repository is intended to be an example of using `gitops-pusher` for
documentation reasons. These are the rules for my personal tailnet.

## Contributions

If you are not a part of `cetacean`, please do not contribute to this
repository. It is intended to demonstrate `gitops-pusher` and monitor it to make
sure it continues working.

If you are a part of `cetacean`, make changes to the ACL file and open a pull
request to get them tested. If tests pass, your changes will usually be mergable
without too much ceremony.
