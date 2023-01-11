# Set up actions runner + app

1. Setup env
   - `export GHA_RUNNER_TOKEN=<token>`
   - `export GHA_RUNNER_REPO=<repo you want to run actions from>`
     - message @sipsma for adding the repo to the test gh app
   - `export GH_APP_ID=278907` - "sipsma-test" app
   - `export GH_WEBHOOK_SECRET=<secret>` - get from @sipsma
   - `export GH_PRIVATE_KEY=<key>` - get from @sipsma
1. ngrok
   - Install ngrok https://dashboard.ngrok.com/get-started/setup
   - run `ngrok http 45363`
   - (TODO: this sucks, need better way) ask @sipsma to change the github app to use the temporary ngrok URL
1. `go run ./cmd/runner/main.go`
1. Start a workflow job in the github repo configured above using `runs-on: dagger-runner`
