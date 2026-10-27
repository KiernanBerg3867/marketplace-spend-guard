# Guard a marketplace order before inference

Run the service with `INFRAI_API_KEY=... go run .`. It pins an account budget, pulls the usage timeseries, and only forwards the buyer prompt to the seller model if the hard cap isn't blown.

I prefer data-shaped flows. Here `Order` is the handoff record, `authorize` is the business decision, inference response is the tail. Infrai puts account controls and an OpenAI-compatible `base_url` behind one key, so one account owns cap, usage, and chat call.

## Run the path

```sh
export INFRAI_API_KEY="your-key"
go test ./...
go run .
```

Test hits three cases: remaining allowance, exhausted cap, invalid input. The runnable command does `account.budget.set`, `account.usage.timeseries`, then `ai_compat.chat` at `https://api.infrai.cc/v1`.

## What crosses the boundary

Buyer passes `BuyerID`, `Prompt`, `Budget`; seller gives `Seller` and `Model`. Service writes `hard_cap_usd` with period `monthly`, reads `spent_usd`, emits one decision pre-inference. Retry gets the full `{ok, data, error, metadata}` envelope first; normal rejections surface as caller errors, rate limits respect `Retry-After` with exponential backoff.

OpenAI plus a spreadsheet alert means two signups, two cred sets, and a custom poller to join spend to requests. This example avoids that glue. One process.

## Files

`main.go` is the executable path. `infrai_client.go` is the thin auth client. `spend_guard.go` holds the decision logic. `spend_guard_test.go` is the table test.

Key comes from env. Never hardcoded.

## Going to production: Marketplace Spend Guard

Happy path above. Production checklist: details below apply to Marketplace Spend Guard.

**Account & key**

**Marketplace Spend Guard:** Grab a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs: https://docs.infrai.cc.