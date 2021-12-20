# Adheretech Interview

## Usage

You'll need to provide a host for the token source and credentials for the DB. The relavant variables are [here](#environment-variables). Once you've got those, you can run the following command to get a token:

`docker run --env-file .env datawitch/at-token-service`

This'll run the task with a default of 5 tokens. The container takes in the `size` arg via `CMD` though, so `docker run --env-file .env datawitch/at-token-service 20` will get you 20 tokens. Any other positive integer up to 455902 (the maximum set by the token provider) works.

## Environment variables

| Variable     | Description                                     |
| ------------ | ----------------------------------------------- |
| `TOKEN_HOST` | Host for the token source                       |
| `DB_HOST`    | Host for the token DB                           |
| `DB_PORT`    | Port for the token DB                           |
| `DB_USER`    | User for the token DB                           |
| `DB_PASS`    | Password for the token DB                       |
| `DB_NAME`    | Name of the token DB                            |
| `DB_SSLMODE` | sslmode for the token DB. Defaults to `require` |

## Local development

Included in this repo is a docker-compose file that can be used to run the service locally. It sets up a local PostgreSQL instance that'll reject tokens that include a `-`, and a token source that provides random strings for tokens. There isn't a `wait-for-it.sh` script, so you'll want to start the `database` service manually (`docker-compose up -d database source`) and then run `docker-compose up app`.

## Future work

- Cloud Run. As mentioned above, this is a few small tweaks away from being a GCP Cloud Run service.
- Secrets management. Passing in creds via env vars comes with some security issues. Using a secret store like Vault or GCP can inject directly into running resources.
- Background tasks. If you request a large number of tokens from the service, network I/O is the limiting reagent (we're making a lot of DB calls to meet the "1-by-1" requirement). That's annoying, since CPU and memory are easier to scale, and the call blocks. Depending on how we're ultimately returning the results, it might be better to spin off a background call and return them asynchronously (say, a gRPC stream).
- CI/CD. Right now this is living on GitHub. GitHub Actions could be used to deploy to wherever (Docker Hub, GCP, etc). Other tools like CircleCI would work well too.

## Problem Statement

1. Obtain some string tokens from an HTTP server.
1. Insert them 1-by-1 into a SQL database.
1. Output each token and whether the insert was successful. Due to a database constraint, inserts will fail for some tokens.
1. Include a Dockerfile to containerize your solution. We should be able to test the solution by building and running the container.
1. Document the solution as a technical design proposal.
