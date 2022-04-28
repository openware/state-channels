# Liabilities channel

## Trade

### State represenation

```yaml
# Participant liabilities
Participant:
    # Liabilities per asset
    Asset:
        # Liability amount to other participant
        Participant: amount
```

**Note:**  
Here we store particpant liabilities per asset, in order to easily calculate the total amount of liabilities for each asset.  
Allowing us to use this informations for collateralization calculations.

**Example state**:  
```yaml
Broker A:
    BTC:
        Broker B: 1
Broker B:
    USDT:
        Broker A: 40000
```

### Workflow

```mermaid
sequenceDiagram
    participant Alice
    participant Broker A
    participant Broker B
    participant Bob
    Bob->>Broker B: Sell limit order: 1 BTC @ 40,000 USDT
    Alice->>Broker A: Buy market order: 1 BTC
    alt Best liquitidy is on Broker B order book: 1 BTC @ 40,000 USDT<br>and Broker A has enough collateral
        Broker A-->>Broker A: Create liabilities update
        Note right of Broker A: Broker A -> Borker B: +1 BTC<br>Broker B -> Broker A: +40,000 USDT
        Broker A-->>Broker A: Sign liabilities update
        Broker A->>Broker B: Liabilities update
        Broker B-->>Broker B: Infer liquitidy request from received liabilities update
        Note right of Broker B: 1 BTC @ 40,000 USDT
        Broker B-->>Broker B: Check Broker A collateralization
        alt Broker A has enough collateral<br>and Broker B order book (still) contains requested liquidity
            Broker B->>Bob: Matched order
            Broker B-->>Broker B: Sign liabilities update
            Broker B->>Broker A: Liabilities update
            Broker A->>Alice: Matched order
        else Broker A does not have enough collateral<br>or liquidity request could not be matched at given price
            Broker B->>Broker A: Liquidity request error
            Broker A-->>Broker A: Revert liabilities update
            Broker A-->>Broker A: Try to match order again
        end
    else Best liquidity is on Broker A order book<br>or Broker A collateral is not sufficient
        Broker A->>Alice: Matched order
    end
```

### Notes

**Interesing concept:**  
One broker could take liquidity from another, without a market order from one of his traders.  
This could be a good technique for said broker to **steal liquidity** from others brokers (to put it on his local order book).  
*Knowing that the liquidity on others brokers is very interesting to take now, compared to current markets conditions.*

## Settlement

### One way

```mermaid
sequenceDiagram
    participant Broker A
    participant Broker B
    Note over Broker A,Broker B: Broker A -> Broker B: 1 BTC<br>Broker B -> Broker A: 40,000 USDT
    Broker A-->>Broker A: Send 1 BTC to Broker B
    Broker A-->>Broker A: Create liabilities update
    Note right of Broker A: Broker A -> Borker B: -1 BTC
    Broker A-->>Broker A: Sign liabilities update
    Broker A->>Broker B: Tx Id + Liabilities update
    Broker B-->>Broker B: Wait for Tx Id confirmation
    Broker B-->>Broker B: Sign liabilities update
    Broker B->>Broker A: Liabilities update
```

### Dispute resolution

**What happen if receiving broker does not sign the liability update ?**
