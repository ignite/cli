/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

/**
* `Any` contains an arbitrary serialized protocol buffer message along with a
URL that describes the type of the serialized message.

Protobuf library provides support to pack/unpack Any values in the form
of utility functions or additional generated methods of the Any type.

Example 1: Pack and unpack a message in C++.

    Foo foo = ...;
    Any any;
    any.PackFrom(foo);
    ...
    if (any.UnpackTo(&foo)) {
      ...
    }

Example 2: Pack and unpack a message in Java.

    Foo foo = ...;
    Any any = Any.pack(foo);
    ...
    if (any.is(Foo.class)) {
      foo = any.unpack(Foo.class);
    }

 Example 3: Pack and unpack a message in Python.

    foo = Foo(...)
    any = Any()
    any.Pack(foo)
    ...
    if any.Is(Foo.DESCRIPTOR):
      any.Unpack(foo)
      ...

 Example 4: Pack and unpack a message in Go

     foo := &pb.Foo{...}
     any, err := anypb.New(foo)
     if err != nil {
       ...
     }
     ...
     foo := &pb.Foo{}
     if err := any.UnmarshalTo(foo); err != nil {
       ...
     }

The pack methods provided by protobuf library will by default use
'type.googleapis.com/full.type.name' as the type URL and the unpack
methods only use the fully qualified type name after the last '/'
in the type URL, for example "foo.bar.com/x/y.z" will yield type
name "y.z".


JSON
====
The JSON representation of an `Any` value uses the regular
representation of the deserialized, embedded message, with an
additional field `@type` which contains the type URL. Example:

    package google.profile;
    message Person {
      string first_name = 1;
      string last_name = 2;
    }

    {
      "@type": "type.googleapis.com/google.profile.Person",
      "firstName": <string>,
      "lastName": <string>
    }

If the embedded message type is well-known and has a custom JSON
representation, that representation will be embedded adding a field
`value` which holds the custom JSON in addition to the `@type`
field. Example (for message [google.protobuf.Duration][]):

    {
      "@type": "type.googleapis.com/google.protobuf.Duration",
      "value": "1.212s"
    }
*/
export interface ProtobufAny {
  /**
   * A URL/resource name that uniquely identifies the type of the serialized
   * protocol buffer message. This string must contain at least
   * one "/" character. The last segment of the URL's path must represent
   * the fully qualified name of the type (as in
   * `path/google.protobuf.Duration`). The name should be in a canonical form
   * (e.g., leading "." is not accepted).
   *
   * In practice, teams usually precompile into the binary all types that they
   * expect it to use in the context of Any. However, for URLs which use the
   * scheme `http`, `https`, or no scheme, one can optionally set up a type
   * server that maps type URLs to message definitions as follows:
   * * If no scheme is provided, `https` is assumed.
   * * An HTTP GET on the URL must yield a [google.protobuf.Type][]
   *   value in binary format, or produce an error.
   * * Applications are allowed to cache lookup results based on the
   *   URL, or have them precompiled into a binary to avoid any
   *   lookup. Therefore, binary compatibility needs to be preserved
   *   on changes to types. (Use versioned type names to manage
   *   breaking changes.)
   * Note: this functionality is not currently available in the official
   * protobuf release, and it is not used for type URLs beginning with
   * type.googleapis.com.
   * Schemes other than `http`, `https` (or the empty scheme) might be
   * used with implementation specific semantics.
   */
  "@type"?: string;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

/**
* Exec defines modes of execution of a proposal on creation or on new vote.

 - EXEC_UNSPECIFIED: An empty value means that there should be a separate
MsgExec request for the proposal to execute.
 - EXEC_TRY: Try to execute the proposal immediately.
If the proposal is not allowed per the DecisionPolicy,
the proposal will still be open and could
be executed at a later point.
*/
export enum V1Exec {
  EXEC_UNSPECIFIED = "EXEC_UNSPECIFIED",
  EXEC_TRY = "EXEC_TRY",
}

/**
 * GroupInfo represents the high-level on-chain information for a group.
 */
export interface V1GroupInfo {
  /**
   * id is the unique ID of the group.
   * @format uint64
   */
  id?: string;

  /** admin is the account address of the group's admin. */
  admin?: string;

  /** metadata is any arbitrary metadata to attached to the group. */
  metadata?: string;

  /**
   * version is used to track changes to a group's membership structure that
   * would break existing proposals. Whenever any members weight is changed,
   * or any member is added or removed this version is incremented and will
   * cause proposals based on older versions of this group to fail
   * @format uint64
   */
  version?: string;

  /** total_weight is the sum of the group members' weights. */
  total_weight?: string;

  /**
   * created_at is a timestamp specifying when a group was created.
   * @format date-time
   */
  created_at?: string;
}

/**
 * GroupMember represents the relationship between a group and a member.
 */
export interface V1GroupMember {
  /**
   * group_id is the unique ID of the group.
   * @format uint64
   */
  group_id?: string;

  /**
   * Member represents a group member with an account address,
   * non-zero weight, metadata and added_at timestamp.
   */
  member?: V1Member;
}

/**
 * GroupPolicyInfo represents the high-level on-chain information for a group policy.
 */
export interface V1GroupPolicyInfo {
  /** address is the account address of group policy. */
  address?: string;

  /**
   * group_id is the unique ID of the group.
   * @format uint64
   */
  group_id?: string;

  /** admin is the account address of the group admin. */
  admin?: string;

  /** metadata is any arbitrary metadata to attached to the group policy. */
  metadata?: string;

  /**
   * version is used to track changes to a group's GroupPolicyInfo structure that
   * would create a different result on a running proposal.
   * @format uint64
   */
  version?: string;

  /**
   * `Any` contains an arbitrary serialized protocol buffer message along with a
   * URL that describes the type of the serialized message.
   *
   * Protobuf library provides support to pack/unpack Any values in the form
   * of utility functions or additional generated methods of the Any type.
   * Example 1: Pack and unpack a message in C++.
   *     Foo foo = ...;
   *     Any any;
   *     any.PackFrom(foo);
   *     ...
   *     if (any.UnpackTo(&foo)) {
   *       ...
   *     }
   * Example 2: Pack and unpack a message in Java.
   *     Any any = Any.pack(foo);
   *     if (any.is(Foo.class)) {
   *       foo = any.unpack(Foo.class);
   *  Example 3: Pack and unpack a message in Python.
   *     foo = Foo(...)
   *     any = Any()
   *     any.Pack(foo)
   *     if any.Is(Foo.DESCRIPTOR):
   *       any.Unpack(foo)
   *  Example 4: Pack and unpack a message in Go
   *      foo := &pb.Foo{...}
   *      any, err := anypb.New(foo)
   *      if err != nil {
   *        ...
   *      }
   *      ...
   *      foo := &pb.Foo{}
   *      if err := any.UnmarshalTo(foo); err != nil {
   * The pack methods provided by protobuf library will by default use
   * 'type.googleapis.com/full.type.name' as the type URL and the unpack
   * methods only use the fully qualified type name after the last '/'
   * in the type URL, for example "foo.bar.com/x/y.z" will yield type
   * name "y.z".
   * JSON
   * ====
   * The JSON representation of an `Any` value uses the regular
   * representation of the deserialized, embedded message, with an
   * additional field `@type` which contains the type URL. Example:
   *     package google.profile;
   *     message Person {
   *       string first_name = 1;
   *       string last_name = 2;
   *     {
   *       "@type": "type.googleapis.com/google.profile.Person",
   *       "firstName": <string>,
   *       "lastName": <string>
   * If the embedded message type is well-known and has a custom JSON
   * representation, that representation will be embedded adding a field
   * `value` which holds the custom JSON in addition to the `@type`
   * field. Example (for message [google.protobuf.Duration][]):
   *       "@type": "type.googleapis.com/google.protobuf.Duration",
   *       "value": "1.212s"
   */
  decision_policy?: ProtobufAny;

  /**
   * created_at is a timestamp specifying when a group policy was created.
   * @format date-time
   */
  created_at?: string;
}

/**
* Member represents a group member with an account address,
non-zero weight, metadata and added_at timestamp.
*/
export interface V1Member {
  /** address is the member's account address. */
  address?: string;

  /** weight is the member's voting weight that should be greater than 0. */
  weight?: string;

  /** metadata is any arbitrary metadata attached to the member. */
  metadata?: string;

  /**
   * added_at is a timestamp specifying when a member was added.
   * @format date-time
   */
  added_at?: string;
}

/**
* MemberRequest represents a group member to be used in Msg server requests.
Contrary to `Member`, it doesn't have any `added_at` field
since this field cannot be set as part of requests.
*/
export interface V1MemberRequest {
  /** address is the member's account address. */
  address?: string;

  /** weight is the member's voting weight that should be greater than 0. */
  weight?: string;

  /** metadata is any arbitrary metadata attached to the member. */
  metadata?: string;
}

/**
 * MsgCreateGroupPolicyResponse is the Msg/CreateGroupPolicy response type.
 */
export interface V1MsgCreateGroupPolicyResponse {
  /** address is the account address of the newly created group policy. */
  address?: string;
}

/**
 * MsgCreateGroupResponse is the Msg/CreateGroup response type.
 */
export interface V1MsgCreateGroupResponse {
  /**
   * group_id is the unique ID of the newly created group.
   * @format uint64
   */
  group_id?: string;
}

/**
 * MsgCreateGroupWithPolicyResponse is the Msg/CreateGroupWithPolicy response type.
 */
export interface V1MsgCreateGroupWithPolicyResponse {
  /**
   * group_id is the unique ID of the newly created group with policy.
   * @format uint64
   */
  group_id?: string;

  /** group_policy_address is the account address of the newly created group policy. */
  group_policy_address?: string;
}

/**
 * MsgExecResponse is the Msg/Exec request type.
 */
export interface V1MsgExecResponse {
  /**
   * ProposalExecutorResult defines types of proposal executor results.
   *
   *  - PROPOSAL_EXECUTOR_RESULT_UNSPECIFIED: An empty value is not allowed.
   *  - PROPOSAL_EXECUTOR_RESULT_NOT_RUN: We have not yet run the executor.
   *  - PROPOSAL_EXECUTOR_RESULT_SUCCESS: The executor was successful and proposed action updated state.
   *  - PROPOSAL_EXECUTOR_RESULT_FAILURE: The executor returned an error and proposed action didn't update state.
   */
  result?: V1ProposalExecutorResult;
}

/**
 * MsgLeaveGroupResponse is the Msg/LeaveGroup response type.
 */
export type V1MsgLeaveGroupResponse = object;

/**
 * MsgSubmitProposalResponse is the Msg/SubmitProposal response type.
 */
export interface V1MsgSubmitProposalResponse {
  /**
   * proposal is the unique ID of the proposal.
   * @format uint64
   */
  proposal_id?: string;
}

/**
 * MsgUpdateGroupAdminResponse is the Msg/UpdateGroupAdmin response type.
 */
export type V1MsgUpdateGroupAdminResponse = object;

/**
 * MsgUpdateGroupMembersResponse is the Msg/UpdateGroupMembers response type.
 */
export type V1MsgUpdateGroupMembersResponse = object;

/**
 * MsgUpdateGroupMetadataResponse is the Msg/UpdateGroupMetadata response type.
 */
export type V1MsgUpdateGroupMetadataResponse = object;

/**
 * MsgUpdateGroupPolicyAdminResponse is the Msg/UpdateGroupPolicyAdmin response type.
 */
export type V1MsgUpdateGroupPolicyAdminResponse = object;

/**
 * MsgUpdateGroupPolicyDecisionPolicyResponse is the Msg/UpdateGroupPolicyDecisionPolicy response type.
 */
export type V1MsgUpdateGroupPolicyDecisionPolicyResponse = object;

/**
 * MsgUpdateGroupPolicyMetadataResponse is the Msg/UpdateGroupPolicyMetadata response type.
 */
export type V1MsgUpdateGroupPolicyMetadataResponse = object;

/**
 * MsgVoteResponse is the Msg/Vote response type.
 */
export type V1MsgVoteResponse = object;

/**
 * MsgWithdrawProposalResponse is the Msg/WithdrawProposal response type.
 */
export type V1MsgWithdrawProposalResponse = object;

/**
* Proposal defines a group proposal. Any member of a group can submit a proposal
for a group policy to decide upon.
A proposal consists of a set of `sdk.Msg`s that will be executed if the proposal
passes as well as some optional metadata associated with the proposal.
*/
export interface V1Proposal {
  /**
   * id is the unique id of the proposal.
   * @format uint64
   */
  id?: string;

  /** group_policy_address is the account address of group policy. */
  group_policy_address?: string;

  /** metadata is any arbitrary metadata to attached to the proposal. */
  metadata?: string;

  /** proposers are the account addresses of the proposers. */
  proposers?: string[];

  /**
   * submit_time is a timestamp specifying when a proposal was submitted.
   * @format date-time
   */
  submit_time?: string;

  /**
   * group_version tracks the version of the group at proposal submission.
   * This field is here for informational purposes only.
   * @format uint64
   */
  group_version?: string;

  /**
   * group_policy_version tracks the version of the group policy at proposal submission.
   * When a decision policy is changed, existing proposals from previous policy
   * versions will become invalid with the `ABORTED` status.
   * This field is here for informational purposes only.
   * @format uint64
   */
  group_policy_version?: string;

  /**
   * ProposalStatus defines proposal statuses.
   *
   *  - PROPOSAL_STATUS_UNSPECIFIED: An empty value is invalid and not allowed.
   *  - PROPOSAL_STATUS_SUBMITTED: Initial status of a proposal when submitted.
   *  - PROPOSAL_STATUS_ACCEPTED: Final status of a proposal when the final tally is done and the outcome
   * passes the group policy's decision policy.
   *  - PROPOSAL_STATUS_REJECTED: Final status of a proposal when the final tally is done and the outcome
   * is rejected by the group policy's decision policy.
   *  - PROPOSAL_STATUS_ABORTED: Final status of a proposal when the group policy is modified before the
   * final tally.
   *  - PROPOSAL_STATUS_WITHDRAWN: A proposal can be withdrawn before the voting start time by the owner.
   * When this happens the final status is Withdrawn.
   */
  status?: V1ProposalStatus;

  /** TallyResult represents the sum of weighted votes for each vote option. */
  final_tally_result?: V1TallyResult;

  /**
   * voting_period_end is the timestamp before which voting must be done.
   * Unless a successfull MsgExec is called before (to execute a proposal whose
   * tally is successful before the voting period ends), tallying will be done
   * at this point, and the `final_tally_result`and `status` fields will be
   * accordingly updated.
   * @format date-time
   */
  voting_period_end?: string;

  /**
   * ProposalExecutorResult defines types of proposal executor results.
   *
   *  - PROPOSAL_EXECUTOR_RESULT_UNSPECIFIED: An empty value is not allowed.
   *  - PROPOSAL_EXECUTOR_RESULT_NOT_RUN: We have not yet run the executor.
   *  - PROPOSAL_EXECUTOR_RESULT_SUCCESS: The executor was successful and proposed action updated state.
   *  - PROPOSAL_EXECUTOR_RESULT_FAILURE: The executor returned an error and proposed action didn't update state.
   */
  executor_result?: V1ProposalExecutorResult;

  /** messages is a list of `sdk.Msg`s that will be executed if the proposal passes. */
  messages?: ProtobufAny[];
}

/**
* ProposalExecutorResult defines types of proposal executor results.

 - PROPOSAL_EXECUTOR_RESULT_UNSPECIFIED: An empty value is not allowed.
 - PROPOSAL_EXECUTOR_RESULT_NOT_RUN: We have not yet run the executor.
 - PROPOSAL_EXECUTOR_RESULT_SUCCESS: The executor was successful and proposed action updated state.
 - PROPOSAL_EXECUTOR_RESULT_FAILURE: The executor returned an error and proposed action didn't update state.
*/
export enum V1ProposalExecutorResult {
  PROPOSAL_EXECUTOR_RESULT_UNSPECIFIED = "PROPOSAL_EXECUTOR_RESULT_UNSPECIFIED",
  PROPOSAL_EXECUTOR_RESULT_NOT_RUN = "PROPOSAL_EXECUTOR_RESULT_NOT_RUN",
  PROPOSAL_EXECUTOR_RESULT_SUCCESS = "PROPOSAL_EXECUTOR_RESULT_SUCCESS",
  PROPOSAL_EXECUTOR_RESULT_FAILURE = "PROPOSAL_EXECUTOR_RESULT_FAILURE",
}

/**
* ProposalStatus defines proposal statuses.

 - PROPOSAL_STATUS_UNSPECIFIED: An empty value is invalid and not allowed.
 - PROPOSAL_STATUS_SUBMITTED: Initial status of a proposal when submitted.
 - PROPOSAL_STATUS_ACCEPTED: Final status of a proposal when the final tally is done and the outcome
passes the group policy's decision policy.
 - PROPOSAL_STATUS_REJECTED: Final status of a proposal when the final tally is done and the outcome
is rejected by the group policy's decision policy.
 - PROPOSAL_STATUS_ABORTED: Final status of a proposal when the group policy is modified before the
final tally.
 - PROPOSAL_STATUS_WITHDRAWN: A proposal can be withdrawn before the voting start time by the owner.
When this happens the final status is Withdrawn.
*/
export enum V1ProposalStatus {
  PROPOSAL_STATUS_UNSPECIFIED = "PROPOSAL_STATUS_UNSPECIFIED",
  PROPOSAL_STATUS_SUBMITTED = "PROPOSAL_STATUS_SUBMITTED",
  PROPOSAL_STATUS_ACCEPTED = "PROPOSAL_STATUS_ACCEPTED",
  PROPOSAL_STATUS_REJECTED = "PROPOSAL_STATUS_REJECTED",
  PROPOSAL_STATUS_ABORTED = "PROPOSAL_STATUS_ABORTED",
  PROPOSAL_STATUS_WITHDRAWN = "PROPOSAL_STATUS_WITHDRAWN",
}

/**
 * QueryGroupInfoResponse is the Query/GroupInfo response type.
 */
export interface V1QueryGroupInfoResponse {
  /** GroupInfo represents the high-level on-chain information for a group. */
  info?: V1GroupInfo;
}

/**
 * QueryGroupMembersResponse is the Query/GroupMembersResponse response type.
 */
export interface V1QueryGroupMembersResponse {
  /** members are the members of the group with given group_id. */
  members?: V1GroupMember[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

/**
 * QueryGroupPoliciesByAdminResponse is the Query/GroupPoliciesByAdmin response type.
 */
export interface V1QueryGroupPoliciesByAdminResponse {
  /** group_policies are the group policies info with provided admin. */
  group_policies?: V1GroupPolicyInfo[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

/**
 * QueryGroupPoliciesByGroupResponse is the Query/GroupPoliciesByGroup response type.
 */
export interface V1QueryGroupPoliciesByGroupResponse {
  /** group_policies are the group policies info associated with the provided group. */
  group_policies?: V1GroupPolicyInfo[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

/**
 * QueryGroupPolicyInfoResponse is the Query/GroupPolicyInfo response type.
 */
export interface V1QueryGroupPolicyInfoResponse {
  /** GroupPolicyInfo represents the high-level on-chain information for a group policy. */
  info?: V1GroupPolicyInfo;
}

/**
 * QueryGroupsByAdminResponse is the Query/GroupsByAdminResponse response type.
 */
export interface V1QueryGroupsByAdminResponse {
  /** groups are the groups info with the provided admin. */
  groups?: V1GroupInfo[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

/**
 * QueryGroupsByMemberResponse is the Query/GroupsByMember response type.
 */
export interface V1QueryGroupsByMemberResponse {
  /** groups are the groups info with the provided group member. */
  groups?: V1GroupInfo[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

/**
 * QueryProposalResponse is the Query/Proposal response type.
 */
export interface V1QueryProposalResponse {
  /**
   * Proposal defines a group proposal. Any member of a group can submit a proposal
   * for a group policy to decide upon.
   * A proposal consists of a set of `sdk.Msg`s that will be executed if the proposal
   * passes as well as some optional metadata associated with the proposal.
   */
  proposal?: V1Proposal;
}

/**
 * QueryProposalsByGroupPolicyResponse is the Query/ProposalByGroupPolicy response type.
 */
export interface V1QueryProposalsByGroupPolicyResponse {
  /** proposals are the proposals with given group policy. */
  proposals?: V1Proposal[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

/**
 * QueryTallyResultResponse is the Query/TallyResult response type.
 */
export interface V1QueryTallyResultResponse {
  /** TallyResult represents the sum of weighted votes for each vote option. */
  tally?: V1TallyResult;
}

/**
 * QueryVoteByProposalVoterResponse is the Query/VoteByProposalVoter response type.
 */
export interface V1QueryVoteByProposalVoterResponse {
  /** Vote represents a vote for a proposal. */
  vote?: V1Vote;
}

/**
 * QueryVotesByProposalResponse is the Query/VotesByProposal response type.
 */
export interface V1QueryVotesByProposalResponse {
  /** votes are the list of votes for given proposal_id. */
  votes?: V1Vote[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

/**
 * QueryVotesByVoterResponse is the Query/VotesByVoter response type.
 */
export interface V1QueryVotesByVoterResponse {
  /** votes are the list of votes by given voter. */
  votes?: V1Vote[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

/**
 * TallyResult represents the sum of weighted votes for each vote option.
 */
export interface V1TallyResult {
  /** yes_count is the weighted sum of yes votes. */
  yes_count?: string;

  /** abstain_count is the weighted sum of abstainers. */
  abstain_count?: string;

  /** no_count is the weighted sum of no votes. */
  no_count?: string;

  /** no_with_veto_count is the weighted sum of veto. */
  no_with_veto_count?: string;
}

/**
 * Vote represents a vote for a proposal.
 */
export interface V1Vote {
  /**
   * proposal is the unique ID of the proposal.
   * @format uint64
   */
  proposal_id?: string;

  /** voter is the account address of the voter. */
  voter?: string;

  /**
   * VoteOption enumerates the valid vote options for a given proposal.
   *
   *  - VOTE_OPTION_UNSPECIFIED: VOTE_OPTION_UNSPECIFIED defines an unspecified vote option which will
   * return an error.
   *  - VOTE_OPTION_YES: VOTE_OPTION_YES defines a yes vote option.
   *  - VOTE_OPTION_ABSTAIN: VOTE_OPTION_ABSTAIN defines an abstain vote option.
   *  - VOTE_OPTION_NO: VOTE_OPTION_NO defines a no vote option.
   *  - VOTE_OPTION_NO_WITH_VETO: VOTE_OPTION_NO_WITH_VETO defines a no with veto vote option.
   */
  option?: V1VoteOption;

  /** metadata is any arbitrary metadata to attached to the vote. */
  metadata?: string;

  /**
   * submit_time is the timestamp when the vote was submitted.
   * @format date-time
   */
  submit_time?: string;
}

/**
* VoteOption enumerates the valid vote options for a given proposal.

 - VOTE_OPTION_UNSPECIFIED: VOTE_OPTION_UNSPECIFIED defines an unspecified vote option which will
return an error.
 - VOTE_OPTION_YES: VOTE_OPTION_YES defines a yes vote option.
 - VOTE_OPTION_ABSTAIN: VOTE_OPTION_ABSTAIN defines an abstain vote option.
 - VOTE_OPTION_NO: VOTE_OPTION_NO defines a no vote option.
 - VOTE_OPTION_NO_WITH_VETO: VOTE_OPTION_NO_WITH_VETO defines a no with veto vote option.
*/
export enum V1VoteOption {
  VOTE_OPTION_UNSPECIFIED = "VOTE_OPTION_UNSPECIFIED",
  VOTE_OPTION_YES = "VOTE_OPTION_YES",
  VOTE_OPTION_ABSTAIN = "VOTE_OPTION_ABSTAIN",
  VOTE_OPTION_NO = "VOTE_OPTION_NO",
  VOTE_OPTION_NO_WITH_VETO = "VOTE_OPTION_NO_WITH_VETO",
}

/**
* message SomeRequest {
         Foo some_parameter = 1;
         PageRequest pagination = 2;
 }
*/
export interface V1Beta1PageRequest {
  /**
   * key is a value returned in PageResponse.next_key to begin
   * querying the next page most efficiently. Only one of offset or key
   * should be set.
   * @format byte
   */
  key?: string;

  /**
   * offset is a numeric offset that can be used when key is unavailable.
   * It is less efficient than using key. Only one of offset or key should
   * be set.
   * @format uint64
   */
  offset?: string;

  /**
   * limit is the total number of results to be returned in the result page.
   * If left empty it will default to a value to be set by each app.
   * @format uint64
   */
  limit?: string;

  /**
   * count_total is set to true  to indicate that the result set should include
   * a count of the total number of items available for pagination in UIs.
   * count_total is only respected when offset is used. It is ignored when key
   * is set.
   */
  count_total?: boolean;

  /**
   * reverse is set to true if results are to be returned in the descending order.
   *
   * Since: cosmos-sdk 0.43
   */
  reverse?: boolean;
}

/**
* PageResponse is to be embedded in gRPC response messages where the
corresponding request message has used PageRequest.

 message SomeResponse {
         repeated Bar results = 1;
         PageResponse page = 2;
 }
*/
export interface V1Beta1PageResponse {
  /**
   * next_key is the key to be passed to PageRequest.key to
   * query the next page most efficiently. It will be empty if
   * there are no more results.
   * @format byte
   */
  next_key?: string;

  /**
   * total is total number of results available if PageRequest.count_total
   * was set, its value is undefined otherwise
   * @format uint64
   */
  total?: string;
}

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, ResponseType } from "axios";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
}

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({ securityWorker, secure, format, ...axiosConfig }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({ ...axiosConfig, baseURL: axiosConfig.baseURL || "" });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  private mergeRequestParams(params1: AxiosRequestConfig, params2?: AxiosRequestConfig): AxiosRequestConfig {
    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.instance.defaults.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  private createFormData(input: Record<string, unknown>): FormData {
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      formData.append(
        key,
        property instanceof Blob
          ? property
          : typeof property === "object" && property !== null
          ? JSON.stringify(property)
          : `${property}`,
      );
      return formData;
    }, new FormData());
  }

  public request = async <T = any, _E = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = (format && this.format) || void 0;

    if (type === ContentType.FormData && body && body !== null && typeof body === "object") {
      requestParams.headers.common = { Accept: "*/*" };
      requestParams.headers.post = {};
      requestParams.headers.put = {};

      body = this.createFormData(body as Record<string, unknown>);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
        ...(requestParams.headers || {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
    });
  };
}

/**
 * @title cosmos/group/v1/events.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryGroupInfo
   * @summary GroupInfo queries group info based on group id.
   * @request GET:/cosmos/group/v1/group_info/{group_id}
   */
  queryGroupInfo = (groupId: string, params: RequestParams = {}) =>
    this.request<V1QueryGroupInfoResponse, RpcStatus>({
      path: `/cosmos/group/v1/group_info/${groupId}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryGroupMembers
   * @summary GroupMembers queries members of a group
   * @request GET:/cosmos/group/v1/group_members/{group_id}
   */
  queryGroupMembers = (
    groupId: string,
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<V1QueryGroupMembersResponse, RpcStatus>({
      path: `/cosmos/group/v1/group_members/${groupId}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryGroupPoliciesByAdmin
   * @summary GroupsByAdmin queries group policies by admin address.
   * @request GET:/cosmos/group/v1/group_policies_by_admin/{admin}
   */
  queryGroupPoliciesByAdmin = (
    admin: string,
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<V1QueryGroupPoliciesByAdminResponse, RpcStatus>({
      path: `/cosmos/group/v1/group_policies_by_admin/${admin}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryGroupPoliciesByGroup
   * @summary GroupPoliciesByGroup queries group policies by group id.
   * @request GET:/cosmos/group/v1/group_policies_by_group/{group_id}
   */
  queryGroupPoliciesByGroup = (
    groupId: string,
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<V1QueryGroupPoliciesByGroupResponse, RpcStatus>({
      path: `/cosmos/group/v1/group_policies_by_group/${groupId}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryGroupPolicyInfo
   * @summary GroupPolicyInfo queries group policy info based on account address of group policy.
   * @request GET:/cosmos/group/v1/group_policy_info/{address}
   */
  queryGroupPolicyInfo = (address: string, params: RequestParams = {}) =>
    this.request<V1QueryGroupPolicyInfoResponse, RpcStatus>({
      path: `/cosmos/group/v1/group_policy_info/${address}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryGroupsByAdmin
   * @summary GroupsByAdmin queries groups by admin address.
   * @request GET:/cosmos/group/v1/groups_by_admin/{admin}
   */
  queryGroupsByAdmin = (
    admin: string,
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<V1QueryGroupsByAdminResponse, RpcStatus>({
      path: `/cosmos/group/v1/groups_by_admin/${admin}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryGroupsByMember
   * @summary GroupsByMember queries groups by member address.
   * @request GET:/cosmos/group/v1/groups_by_member/{address}
   */
  queryGroupsByMember = (
    address: string,
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<V1QueryGroupsByMemberResponse, RpcStatus>({
      path: `/cosmos/group/v1/groups_by_member/${address}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryProposal
   * @summary Proposal queries a proposal based on proposal id.
   * @request GET:/cosmos/group/v1/proposal/{proposal_id}
   */
  queryProposal = (proposalId: string, params: RequestParams = {}) =>
    this.request<V1QueryProposalResponse, RpcStatus>({
      path: `/cosmos/group/v1/proposal/${proposalId}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
 * No description
 * 
 * @tags Query
 * @name QueryTallyResult
 * @summary TallyResult returns the tally result of a proposal. If the proposal is
still in voting period, then this query computes the current tally state,
which might not be final. On the other hand, if the proposal is final,
then it simply returns the `final_tally_result` state stored in the
proposal itself.
 * @request GET:/cosmos/group/v1/proposals/{proposal_id}/tally
 */
  queryTallyResult = (proposalId: string, params: RequestParams = {}) =>
    this.request<V1QueryTallyResultResponse, RpcStatus>({
      path: `/cosmos/group/v1/proposals/${proposalId}/tally`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryProposalsByGroupPolicy
   * @summary ProposalsByGroupPolicy queries proposals based on account address of group policy.
   * @request GET:/cosmos/group/v1/proposals_by_group_policy/{address}
   */
  queryProposalsByGroupPolicy = (
    address: string,
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<V1QueryProposalsByGroupPolicyResponse, RpcStatus>({
      path: `/cosmos/group/v1/proposals_by_group_policy/${address}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryVoteByProposalVoter
   * @summary VoteByProposalVoter queries a vote by proposal id and voter.
   * @request GET:/cosmos/group/v1/vote_by_proposal_voter/{proposal_id}/{voter}
   */
  queryVoteByProposalVoter = (proposalId: string, voter: string, params: RequestParams = {}) =>
    this.request<V1QueryVoteByProposalVoterResponse, RpcStatus>({
      path: `/cosmos/group/v1/vote_by_proposal_voter/${proposalId}/${voter}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryVotesByProposal
   * @summary VotesByProposal queries a vote by proposal.
   * @request GET:/cosmos/group/v1/votes_by_proposal/{proposal_id}
   */
  queryVotesByProposal = (
    proposalId: string,
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<V1QueryVotesByProposalResponse, RpcStatus>({
      path: `/cosmos/group/v1/votes_by_proposal/${proposalId}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryVotesByVoter
   * @summary VotesByVoter queries a vote by voter.
   * @request GET:/cosmos/group/v1/votes_by_voter/{voter}
   */
  queryVotesByVoter = (
    voter: string,
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<V1QueryVotesByVoterResponse, RpcStatus>({
      path: `/cosmos/group/v1/votes_by_voter/${voter}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });
}
