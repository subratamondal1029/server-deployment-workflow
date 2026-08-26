# server-deployment-workflow
Creating a server deployment workflow. with go, docker, github CI.

```text
1. Developer merges code into main
                ↓
2. GitHub Actions starts
                ↓
3. GitHub CI runs tests
                ↓
4. GitHub CI builds the Go binary
                ↓
5. GitHub CI creates a GitHub Release
   containing the built binary
                ↓
6. Production server downloads
   the release artifact only
                ↓
7. Server puts the binary into
   the Docker build directory
                ↓
8. Docker builds a new image
   FROM scratch
   + copies the binary
                ↓
9. Stop/remove the old container (more then that)
                ↓
10. Start the new Docker container
                ↓
11. Go server runs inside the
    minimal scratch container
```
