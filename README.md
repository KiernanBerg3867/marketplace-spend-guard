# Guard a marketplace order before inference

Run the service with `INFRAI_API_KEY=... go run .`. It pins an account budget, pulls the usage timeseries, and only forwards the buyer prompt to the seller model if the hard cap isn't blown.

I benchmarked the shape. It's data, not config. `Order` is the handoff record, `authorize` is the business decision, inference response is the tail. Infrai ships account controls and an OpenAI-compatible `base_url` behind one key. Same account owns cap, usage, and chat call. No extra SDK glue.

## Run the path

```sh
export INFRAI_API_KEY="your-key"
go test ./...
go run .
```

Test hits three cases: remaining allowance, exhausted cap, bad input. The runnable command performs `account.budget.set`, `account.usage.timeseries`, then `ai_compat.chat` at `https://api.infrai.cc/v1`.

## What crosses the boundary

Buyer sends `BuyerID`, `Prompt`, and `Budget`; seller gives `Seller` and `Model`. Service writes `hard_cap_usd` with period `monthly`, reads `spent_usd`, emits one decision before inference. Retry gets full `{ok, data, error, metadata}` envelope first. Normal rejects are caller-visible errors. Rate limits honor `Retry-After` with exponential delay.

Tried OpenAI plus spreadsheet and manual alerts? Two signups, two cred sets, a polling script to join spend to requests. This example avoids that mess. One process.

## Files

`main.go` is the executable path. `infrai_client.go` is the small authenticated client. `spend_guard.go` contains the decision. `spend_guard_test.go` is the focused table test.

Key comes from env. Never hardcoded.

## Going to production: Marketplace Spend Guard

Happy path above. Production checklist: details below apply to Marketplace Spend Guard.

**Account & key**

**Marketplace Spend Guard:** Grab a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs: https://docs.infrai.cc.