# Proposal: Token Service

Author(s): Jessica Ellis

Last Updated: 2021-12-18

Discussion at [https://github.com/jessdwitch/adheretech-interview/issues/1](https://github.com/jessdwitch/adheretech-interview/issues/1)

## Abstract

This document proposes a docker app to provide and store tokens.

## Background

Currently, an external API is able to provide tokens, and a PostgreSQL server is already stood up for persisting those tokens.

Deployment is via a docker application.

### Samples Request to Token provider

`POST /be-interview-env-datasource-84c27f85?size=10`

```json
{}
```

Response:

```
HTTP/2 200
content-type: text/plain; charset=utf-8

UNFvwishxzYsXUUchI4NBw
pOQMl-eKw9Y1bUy3STwJzQ
Jg535-x8sYQnnRPc3NiPOA
YYg7IdoPwiFjadFeDp6rtg
DAHQkEMMejUVuPSMW3Y7Pg
K-CiXx3xIY2jpEhkjwViRQ
Ej8G4vet4IMbnHM8SckUrA
GC_oqK8FMvN1gwb-byrzYA
VW4so8LIdLmzXkSQqQCGcA
JZKOZ1wvV4OJpMoRLDhRZg
```

---

*Note*: The token provider does not always return exactly the requested number of tokens. See [Open Questions](#token-provider-does-not-always-return-the-requested-number-of-tokens).

### Requirements

- Docker
- SQL

## Proposal

1. To retrieve tokens from Token Source;
1. To store tokens in the `secret_tokens.data`;
1. Output tokens with an indication of whether they were successfully stored.

### Advantages

- Does not modify existing implementations of Token Source or DB
- Portable
- Easy to test
- Can be deployed to Google Cloud Run without much work

### Disadvantages

- No way to retrieve tokens after they have been stored
- No monitoring
- Token limit set by an external API

## Compatibility

Token Service is built around existing resources. It should be fully compatible with them.

## Implementation

### Tokens Service v1

Implemented as a docker application. Simple to run:

`docker run token_service <SIZE>`

## Open Questions

Normally, I'd have a requirements meetings to ask these questions and get customer sign-off at the end, but since it's the weekend and this isn't a real production service, for expediency's sake, I'm going to make up some resolutions and note them below.

### Token Provider does not always return the requested number of tokens

It is unknown whether this is a bug, or an undocumented restriction on the token service.

This behavior is shown in the following sample request:

```
POST /be-interview-env-datasource-84c27f85?size=5

{}
```

Response:

```
HTTP/2 200
content-type: text/plain; charset=utf-8

5YutqDANtWY2Rm8LKk2Cgg
DZNe3TVX_H9IYcAvKR8jSQ
HkOlYpMaqwIbutrNvuu8GA
gIrb_Dqbi_RFtX2HLEDNBw
```

Note that there are only 4 tokens provided in spite of the `size` parameter being set to 5.

**Resolution:** Bug will be fixed in the future. Token service must be resilient to this behavior until then.

### Must tokens be processed in order?

Technical documentation sets the requirement that tokens must be inserted into the DB "1-by-1". This reads to me as each insertion being an individual transaction, but does not suggest that their order matters. Closely related to this question is whether the API response should also maintain the Token Source's order.

**Resolution:** The order of tokens does not matter.
