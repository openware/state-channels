# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [channel.proto](#channel.proto)
    - [ApproveRequest](#.ApproveRequest)
    - [ApproveResponse](#.ApproveResponse)
    - [BoolResponse](#.BoolResponse)
    - [Channel](#.Channel)
    - [ChannelRequest](#.ChannelRequest)
    - [CheckHoldingsResponse](#.CheckHoldingsResponse)
    - [CheckSignatureRequest](#.CheckSignatureRequest)
    - [ConcludeRequest](#.ConcludeRequest)
    - [ConcludeRequest.SignaturesEntry](#.ConcludeRequest.SignaturesEntry)
    - [ConcludeResponse](#.ConcludeResponse)
    - [CurrentStateResponse](#.CurrentStateResponse)
    - [FundChannelRequest](#.FundChannelRequest)
    - [FundChannelResponse](#.FundChannelResponse)
    - [GasStation](#.GasStation)
    - [InitChannelRequest](#.InitChannelRequest)
    - [InitChannelResponse](#.InitChannelResponse)
    - [ProposeResponse](#.ProposeResponse)
    - [SignStateRequest](#.SignStateRequest)
    - [SignStateResponse](#.SignStateResponse)
    - [Signature](#.Signature)
  
    - [ChannelService](#.ChannelService)
  
- [init_proposal.proto](#init_proposal.proto)
    - [AddParticipantRequest](#.AddParticipantRequest)
    - [Contract](#.Contract)
    - [CreateProposalRequest](#.CreateProposalRequest)
    - [CreateProposalResponse](#.CreateProposalResponse)
    - [EmptyInitProposalResponse](#.EmptyInitProposalResponse)
    - [InitialProposal](#.InitialProposal)
    - [Participant](#.Participant)
  
    - [InitProposalService](#.InitProposalService)
  
- [state.proto](#state.proto)
    - [State](#.State)
  
- [state_proposal.proto](#state_proposal.proto)
    - [EmptyProposalResponse](#.EmptyProposalResponse)
    - [LiabilityRequest](#.LiabilityRequest)
    - [StateProposal](#.StateProposal)
    - [StateProposalRequest](#.StateProposalRequest)
  
    - [StateProposalService](#.StateProposalService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="channel.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## channel.proto



<a name=".ApproveRequest"></a>

### ApproveRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| private_key | [bytes](#bytes) |  |  |
| channel | [Channel](#Channel) |  |  |






<a name=".ApproveResponse"></a>

### ApproveResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signature | [Signature](#Signature) |  |  |






<a name=".BoolResponse"></a>

### BoolResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  |  |






<a name=".Channel"></a>

### Channel



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| initial_proposal | [InitialProposal](#InitialProposal) |  |  |
| state | [State](#State) |  |  |
| ch | [bytes](#bytes) |  | channel obj |






<a name=".ChannelRequest"></a>

### ChannelRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| channel | [Channel](#Channel) |  |  |






<a name=".CheckHoldingsResponse"></a>

### CheckHoldingsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| amount | [int64](#int64) |  |  |






<a name=".CheckSignatureRequest"></a>

### CheckSignatureRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signature | [Signature](#Signature) |  |  |
| state | [State](#State) |  |  |
| channel | [Channel](#Channel) |  |  |






<a name=".ConcludeRequest"></a>

### ConcludeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| participant | [Participant](#Participant) |  |  |
| private_key | [bytes](#bytes) |  |  |
| signatures | [ConcludeRequest.SignaturesEntry](#ConcludeRequest.SignaturesEntry) | repeated |  |
| gas_station | [GasStation](#GasStation) | optional |  |
| channel | [Channel](#Channel) |  |  |






<a name=".ConcludeRequest.SignaturesEntry"></a>

### ConcludeRequest.SignaturesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Signature](#Signature) |  |  |






<a name=".ConcludeResponse"></a>

### ConcludeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tx_id | [string](#string) |  |  |






<a name=".CurrentStateResponse"></a>

### CurrentStateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state | [State](#State) |  |  |






<a name=".FundChannelRequest"></a>

### FundChannelRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| participant | [Participant](#Participant) |  |  |
| private_key | [bytes](#bytes) |  |  |
| gas_station | [GasStation](#GasStation) | optional |  |
| channel | [Channel](#Channel) |  |  |






<a name=".FundChannelResponse"></a>

### FundChannelResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tx_id | [string](#string) |  |  |






<a name=".GasStation"></a>

### GasStation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| gas_price | [int64](#int64) |  |  |
| gas_limit | [uint64](#uint64) |  |  |






<a name=".InitChannelRequest"></a>

### InitChannelRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| initial_proposal | [InitialProposal](#InitialProposal) |  |  |
| participant_index | [uint32](#uint32) |  |  |






<a name=".InitChannelResponse"></a>

### InitChannelResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| channel | [Channel](#Channel) |  |  |






<a name=".ProposeResponse"></a>

### ProposeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state_proposal | [StateProposal](#StateProposal) |  |  |






<a name=".SignStateRequest"></a>

### SignStateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state_proposal | [StateProposal](#StateProposal) |  |  |
| private_key | [bytes](#bytes) |  |  |
| channel | [Channel](#Channel) |  |  |






<a name=".SignStateResponse"></a>

### SignStateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signature | [Signature](#Signature) |  |  |






<a name=".Signature"></a>

### Signature



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| r | [bytes](#bytes) |  |  |
| s | [bytes](#bytes) |  |  |
| v | [bytes](#bytes) |  |  |





 

 

 


<a name=".ChannelService"></a>

### ChannelService
channel?

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Init | [.InitChannelRequest](#InitChannelRequest) | [.InitChannelResponse](#InitChannelResponse) | Public |
| ApproveInit | [.ApproveRequest](#ApproveRequest) | [.ApproveResponse](#ApproveResponse) | Public |
| Fund | [.FundChannelRequest](#FundChannelRequest) | [.FundChannelResponse](#FundChannelResponse) | Public |
| ApproveFunding | [.ApproveRequest](#ApproveRequest) | [.ApproveResponse](#ApproveResponse) | Public |
| ProposeState | [.ChannelRequest](#ChannelRequest) | [.ProposeResponse](#ProposeResponse) | Public |
| SignState | [.SignStateRequest](#SignStateRequest) | [.SignStateResponse](#SignStateResponse) | Public |
| Conclude | [.ConcludeRequest](#ConcludeRequest) | [.ConcludeResponse](#ConcludeResponse) | Public |
| CheckSignature | [.CheckSignatureRequest](#CheckSignatureRequest) | [.BoolResponse](#BoolResponse) | Public |
| CurrentState | [.ChannelRequest](#ChannelRequest) | [.CurrentStateResponse](#CurrentStateResponse) | Public |
| CheckHoldings | [.ChannelRequest](#ChannelRequest) | [.CheckHoldingsResponse](#CheckHoldingsResponse) | Public |
| StateIsFinal | [.ChannelRequest](#ChannelRequest) | [.BoolResponse](#BoolResponse) | Public |

 



<a name="init_proposal.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## init_proposal.proto



<a name=".AddParticipantRequest"></a>

### AddParticipantRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| participant | [Participant](#Participant) |  |  |
| initial_proposal | [InitialProposal](#InitialProposal) |  |  |






<a name=".Contract"></a>

### Contract



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| nitro_client | [bytes](#bytes) |  |  |
| asset_address | [string](#string) |  |  |






<a name=".CreateProposalRequest"></a>

### CreateProposalRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| participant | [Participant](#Participant) |  |  |
| rpc_url | [string](#string) |  |  |
| contract_address | [string](#string) |  |  |
| asset_address | [string](#string) |  |  |






<a name=".CreateProposalResponse"></a>

### CreateProposalResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| initial_proposal | [InitialProposal](#InitialProposal) |  |  |






<a name=".EmptyInitProposalResponse"></a>

### EmptyInitProposalResponse







<a name=".InitialProposal"></a>

### InitialProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| participants | [Participant](#Participant) | repeated |  |
| state | [State](#State) |  |  |
| contract | [Contract](#Contract) |  |  |
| channel_nonce | [int64](#int64) |  |  |






<a name=".Participant"></a>

### Participant



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| address | [string](#string) |  |  |
| destination | [string](#string) |  |  |
| locked_amount | [int64](#int64) |  |  |
| index | [uint64](#uint64) |  |  |





 

 

 


<a name=".InitProposalService"></a>

### InitProposalService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateProposal | [.CreateProposalRequest](#CreateProposalRequest) | [.CreateProposalResponse](#CreateProposalResponse) | Public |
| AddParticipant | [.AddParticipantRequest](#AddParticipantRequest) | [.EmptyInitProposalResponse](#EmptyInitProposalResponse) | Public |

 



<a name="state.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## state.proto



<a name=".State"></a>

### State



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| chain_id | [uint64](#uint64) |  |  |
| participants | [string](#string) | repeated |  |
| channel_nonce | [int64](#int64) |  |  |
| app_definition | [string](#string) |  |  |
| challenge_duration | [uint64](#uint64) |  |  |
| app_data | [bytes](#bytes) |  |  |
| outcome | [bytes](#bytes) |  |  |
| turn_num | [uint64](#uint64) |  |  |
| is_final | [bool](#bool) |  |  |





 

 

 

 



<a name="state_proposal.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## state_proposal.proto



<a name=".EmptyProposalResponse"></a>

### EmptyProposalResponse







<a name=".LiabilityRequest"></a>

### LiabilityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| from | [uint64](#uint64) |  |  |
| to | [uint64](#uint64) |  |  |
| asset | [string](#string) |  |  |
| amount | [string](#string) |  |  |
| state_proposal | [StateProposal](#StateProposal) |  |  |






<a name=".StateProposal"></a>

### StateProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state | [State](#State) |  |  |
| liability_state | [bytes](#bytes) |  |  |






<a name=".StateProposalRequest"></a>

### StateProposalRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state_proposal | [StateProposal](#StateProposal) |  |  |





 

 

 


<a name=".StateProposalService"></a>

### StateProposalService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| SetFinal | [.StateProposalRequest](#StateProposalRequest) | [.EmptyProposalResponse](#EmptyProposalResponse) | Public |
| PendingLiability | [.LiabilityRequest](#LiabilityRequest) | [.EmptyProposalResponse](#EmptyProposalResponse) | Public |
| ExecutedLiability | [.LiabilityRequest](#LiabilityRequest) | [.EmptyProposalResponse](#EmptyProposalResponse) | Public |
| RevertLiability | [.LiabilityRequest](#LiabilityRequest) | [.EmptyProposalResponse](#EmptyProposalResponse) | Public |
| ApproveLiabilities | [.StateProposalRequest](#StateProposalRequest) | [.EmptyProposalResponse](#EmptyProposalResponse) | Public |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

