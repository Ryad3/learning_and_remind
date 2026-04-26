# Diffie-Hellman: Key Exchange

> A cryptographic protocol that allows two parties to establish a shared secret over an insecure channel, without ever transmitting it directly.

---

## Exchange Diagram

```
┌──────────────────────────────────────────────────────────────────┐
│           1 — Public parameters (known to all)                   │
│                  p (prime number), g (generator)                 │
└──────────────────────────────────────────────────────────────────┘

        ALICE                                         BOB
   ┌─────────────┐                             ┌─────────────┐
   │ Chooses a   │                             │ Chooses b   │
   │  (secret,   │                             │  (secret,   │
   │  unsent)    │                             │  unsent)    │
   └──────┬──────┘                             └──────┬──────┘
          │                                           │
          ▼                                           ▼
   ┌─────────────┐                             ┌─────────────┐
   │ A = gᵃ mod p│                             │ B = gᵇ mod p│
   │ (pub. key)  │                             │ (pub. key)  │
   └──────┬──────┘                             └──────┬──────┘
          │                                           │
          │   ──────────── A ──────────────►          │
          │                                           │
          │   ◄─────────── B ──────────────           │
          │       (cleartext exchange on              │
          │        the network — visible              │
          │        to an attacker)                    │
          │                                           │
          ▼                                           ▼
   ┌─────────────┐                             ┌─────────────┐
   │ S = Bᵃ mod p│                             │ S = Aᵇ mod p│
   └──────┬──────┘                             └──────┬──────┘
          │                                           │
          └──────────────────┬────────────────────────┘
                             ▼
              ┌──────────────────────────┐
              │   S = gᵃᵇ mod p          │
              │   Shared secret ✓        │
              │   (identical on both     │
              │    sides, never sent)    │
              └──────────────────────────┘

   ATTACKER (network): sees p, g, A, B
                       can NOT retrieve a or b
                       → discrete logarithm problem
```

---

## Step 1 — Public Parameters

Alice and Bob agree on two values, public and known to everyone:

- **p**: a large prime number (e.g., 2048 bits in practice)
- **g**: a generator of the multiplicative group modulo p (typically `g = 2` or `g = 5`)

These parameters can be standardized (RFC 3526, RFC 7919) or negotiated on the fly. Publishing them does not compromise security.

---

## Step 2 — Generating Private Secrets and Public Keys

### Alice

Alice chooses a random integer **a**, her private secret, which she never communicates. She calculates her **public key**:

```
A = gᵃ mod p
```

### Bob

Bob does the same with his own secret **b**:

```
B = gᵇ mod p
```

Alice sends `A` to Bob, Bob sends `B` to Alice. These two values travel **in cleartext** over the network: an attacker sees them, but this is not enough for them to reconstruct the secrets.

---

## Step 3 — Calculating the Shared Secret

### Alice receives B and calculates:

```
S = Bᵃ mod p
  = (gᵇ)ᵃ mod p
  = gᵃᵇ mod p
```

### Bob receives A and calculates:

```
S = Aᵇ mod p
  = (gᵃ)ᵇ mod p
  = gᵃᵇ mod p
```

Both arrive at the **same result** thanks to the properties of modular exponentiation. This secret `S = gᵃᵇ mod p` never transited over the network.

---

## Why an Attacker Cannot Reconstruct S?

The attacker has `p`, `g`, `A`, and `B`. To find `a` from `A = gᵃ mod p`, they would have to solve the **discrete logarithm problem**:

```
a = logᵍ(A) mod p
```

For small values, this is trivial. But with prime numbers of 2048 bits or more, no known classical algorithm can solve this problem in a reasonable amount of time — the best known attack (Number Field Sieve) remains exponentially expensive.

> **Note**: A sufficiently powerful quantum computer could solve this problem using Shor's algorithm. This is why post-quantum cryptography is exploring alternatives (lattices, etc.).

---

## Paint Colors Metaphor

An intuitive way to understand DH is the metaphor of mixing paint:

| Step | Paint |
|---|---|
| Public common color | `g` |
| Alice's secret | her private color (not transmitted) |
| Alice's public key `A` | mixture: common color + Alice's secret |
| Bob's public key `B` | mixture: common color + Bob's secret |
| Shared secret `S` | final mixture: common + Alice + Bob |

Mixing paints is easy in one direction (calculating `gᵃ mod p`), but "unmixing" a color to find the secret ingredient is practically impossible — exactly like the discrete logarithm.

---

## Numerical Example (Simplified Values)

```
p = 23        (prime number)
g = 5         (generator)

Alice chooses: a = 6
Bob chooses  : b = 15

A = 5⁶ mod 23  = 15625 mod 23  = 8
B = 5¹⁵ mod 23 = 30517578125 mod 23 = 19

Alice calculates: S = 19⁶ mod 23  = 47045881 mod 23 = 2
Bob calculates  : S = 8¹⁵ mod 23  = 35184372088832 mod 23 = 2

→ Shared secret S = 2  ✓
```

---

## Usage in SSH

In SSH (RFC 4253), the DH exchange takes place during the **initial handshake** (SSH Transport Layer):

1. The client proposes supported DH algorithms (`kex_algorithms`)
2. Both parties perform the exchange and obtain the shared secret `S`
3. This secret is used to derive the **symmetric session keys** (e.g., AES-256-CTR)
4. The entire rest of the session is encrypted with these derived keys

Modern variants used in SSH are **ECDH** (Elliptic Curve Diffie-Hellman, e.g., Curve25519), which offer equivalent security with much shorter keys and faster calculations.

---

## Limitations and Points of Attention

| Limitation | Detail |
|---|---|
| No native authentication | DH establishes a secret, but does not prove identity. Without verifying the server's key, one is vulnerable to a man-in-the-middle attack. SSH solves this via *host keys*. |
| Vulnerable to quantum computers | Shor's algorithm breaks DH in polynomial time on a quantum computer. |
| Logjam (2015) | Weak DH parameters (512 bits, precomputed) allowed breaking the exchange in real-time. Solved by forcing minimum 2048-bit groups. |
| Ephemeral DH (DHE/ECDHE) | To guarantee *forward secrecy*, DH keys must be regenerated for each session. A static DH (same a and b reused) loses this property. |
