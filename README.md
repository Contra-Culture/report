# report

Report is a structured reports/logs generation library for Golang with no dependencies.

For more details see [report.go](./report.go). The code is very simple.

**Your issues and PRs with implevements are welcome.**

## FAQ

### Why don't you use tests?
The code is quite simple. Testing will make code more complex and slower. For example, for testing purposes I need to add DI for time.Now() to create predictable time objects. It was implemented that way before, but for now I decided to make `report` even simpler.
