# Freezer Service

``Freezer is handling everything related to user's ice tokenomics and statistics about it.``

### Development

These are the crucial/critical operations you will need when developing `Freezer`:

1. If you need to generate a new Authorization Token & UserID for testing locally:
    1. run `make print-token-XXX`, where `XXX` is the role you want for the user.
2. If you need to seed your local database, or even a remote one:
    1. run `make start-seeding`
    2. it requires an .env entry: `MASTER_DB_INSTANCE_ADDRESS=admin:pass@127.0.0.1:3301`
3. `make run-freezer`
    1. This runs the actual read service.
    2. It will feed off of the properties in `./application.yaml`
    3. By default, https://localhost:2443/tokenomics/r runs the Open API (Swagger) entrypoint.
4. `make run-freezer-refrigerant`
    1. This runs the actual write service.
    2. It will feed off of the properties in `./application.yaml`
    3. By default, https://localhost:3443/tokenomics/w runs the Open API (Swagger) entrypoint.
5. `make start-test-environment`
    1. This bootstraps a local test environment with **Freezer**'s dependencies using your `docker` and `docker-compose` daemons.
    2. It is a blocking operation, SIGTERM or SIGINT will kill it.
    3. It will feed off of the properties in `./application.yaml`
        1. MessageBroker GUIs
            1. https://www.conduktor.io
            2. https://www.kafkatool.com
            3. (CLI) https://vectorized.io/redpanda
        2. DB GUIs
            1. https://github.com/tarantool/awesome-tarantool#gui-clients
            2. (CLI) `docker exec -t -i mytarantool console` where `mytarantool` is the container name
6. `make all`
    1. This runs the CI pipeline, locally -- the same pipeline that PR checks run.
    2. Run it before you commit to save time & not wait for PR check to fail remotely.
7. `make local`
    1. This runs the CI pipeline, in a descriptive/debug mode. Run it before you run the "real" one.
8. `make lint`
    1. This runs the linters. It is a part of the other pipelines, so you can run this separately to fix lint issues.
9. `make test`
    1. This runs all tests.
10. `make benchmark`
    1. This runs all benchmarks.
