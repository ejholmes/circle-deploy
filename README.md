# Circle Deploy

GitHub Deployments -> CircleCI glue.

This is a small webhook that can be used to trigger CircleCI builds from GitHub Deployments.

## Usage

Add a webhook to the GitHub repository with the following:

```
https://circle-deploy.herokuapp.com?circle-token=<token>
```

When this app receives a GitHub deployment request for a repo, it will use the [CircleCI builds API](https://circleci.com/docs/nightly-builds) to trigger a build for the repo, passing the following build parameters:

Parameter | Description
----------|------------
`GITHUB_DEPLOYMENT` | Set to the id of the GitHub deployment. You can use this to determine if you should run tests or do a deployment
`GITHUB_DEPLOYMENT_ENVIRONMENT` | The environment that was requested to be deployed to
