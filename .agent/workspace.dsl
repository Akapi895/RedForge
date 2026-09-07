workspace "Governed Security Analysis" "Reviewed V1 C4 architecture baseline for governed security analysis with modular-monolith packaging, interactive operator workflows, explicit logical data boundaries, replayable derived views, and deliberately deferred runtime extraction decisions." {
    !impliedRelationships false

    // V1 ARCHITECTURE DECISIONS:
    // D1. Server-side executable responsibilities run in one Backend Application
    //     as a modular monolith. Reasoning, Execution and Verification remain
    //     explicit internal subsystems with stable seams for future extraction.
    // D2. Web UI and CLI are separate access-client applications in C2.
    //     REST API and MCP HTTP Server remain interfaces/adapters provided by
    //     the Backend Application; no separate API/MCP service is introduced.
    // D3. Web UI is a deliberate V1 client boundary for browser-based operation.
    //     Framework, SPA implementation details and backend transport beyond the
    //     authenticated HTTP boundary remain unspecified.
    // D4. Transactional state, protected evidence and derived read models are
    //     logical persistence boundaries. Their physical engines, servers,
    //     hosting, replication and SQL/graph topology remain unspecified.
    // D5. Internal groups describe ownership/module boundaries, not additional
    //     network/process/service boundaries.
    // D6. STOP/Recovery is explicitly scoped to PHASE or MISSION. Phase STOP
    //     affects only the stopped phase and its dependants; Mission STOP freezes
    //     normal admission across the mission.
    // D7. Web UI is an interactive agent workbench, not a read-only dashboard.
    //     It carries conversational turns, scope/context input, clarification,
    //     plan confirmation, steering, mission controls, review decisions and
    //     authorized visualization queries. Backend remains authoritative.
    // D8. Derived SQL/knowledge-graph/observability views are replayed from
    //     committed canonical events. They may be stale/rebuildable and never
    //     authorize execution, grants, budgets or stop state.
    // D9. Reporting is template-driven. Template selection is operator/system
    //     controlled, frozen artifacts are durable before approval, and release
    //     authorization binds exact bytes, template version and destination.
    //
    // DELIBERATELY DEFERRED:
    // worker thread/subprocess/sandbox choices; independent Reasoning/Execution/
    // Verification runtimes; model-provider placement; identity provider;
    // target-system placement; export destination; AUTHORITY persistence;
    // physical placement/technology for durably frozen report artifacts;
    // physical database topology and deployment.
    //
    // No message broker, service discovery layer, technology vendor, deployment
    // topology or new workflow is introduced merely to prepare for future splits.
    // The scope remains analysis, approved diagnostics and bounded recovery.
    // There is no autonomous exploitation or intrusion workflow.

    model {
        properties {
            "structurizr.groupSeparator" "/"
        }

        operator = person "Operator" "Verified human principal who converses with the system, supplies mission scope/context, confirms plans, steers work, requests PHASE/MISSION stop or explicit resume, selects permitted report templates, and receives approved reports."
        reviewer = person "Authorized Reviewer" "Human approval role from HITL. May be held by the same person as Operator only when current permissions allow; cannot override hard boundaries."

        system = softwareSystem "Governed Security Analysis System" "Security analysis, approved diagnostics and bounded recovery under deterministic governance." {
            group "Access clients" {
                web = container "Web UI" "Interactive browser workbench for conversational mission intake, scope/context entry, clarification, plan confirmation, steering, scoped stop/resume, approval review, system/graph/evidence visualization, report-template selection and approved report delivery." "Browser application; framework unspecified" {
                    tags "Client,Web Client"
                    properties {
                        "source.nodes" "UI"
                        "packaging.v1" "Deliberate C2 decision: Web UI is a separate interactive browser client application."
                        "interaction" "Carries authenticated operator turns, clarification answers, plan confirmations, steering requests, review decisions, mission controls, report-template selection and authorized visualization/query requests."
                        "ownership" "Owns user interaction and presentation only. Backend owns trusted identity attachment, mission/phase state, authorization, approval state, task creation, verified facts and authoritative transitions."
                        "authority" "Interactive confirmation is not execution authority. UI cannot directly mutate canonical state, schedule work, call tools, bypass policy, query stores directly or publish verified facts."
                        "transport" "Authenticated HTTP interaction with Backend; framework, streaming mechanism and SPA implementation details remain unspecified."
                    }
                }
                cli = container "CLI Application" "Command-line entry for authenticated goals and mission controls." "CLI application; implementation and transport unspecified" {
                    tags "Client"
                    properties {
                        "source.nodes" "CLI"
                        "packaging.v1" "Deliberate C2 decision: CLI is a separate client application."
                        "transport" "Deliberately unspecified; this model does not assert REST, MCP or another binding."
                    }
                }
            }

            backend = container "Backend Application" "Interactive REST/MCP entry; deterministic mission interaction, governance and workflow; bounded reasoning, diagnostics, verification and recovery; projections and controlled reporting." "Application runtime unspecified" {
                tags "Backend"
                properties {
                    "packaging.v1" "Modular monolith: control, reasoning, execution and verification share one application lifecycle in V1. This is a deliberate packaging decision, not an inference from the Mermaid."
                    "future.extraction" "Reasoning (AGENTS), Execution (RUNNER) and Verification (VERIFY) keep explicit code/module seams so they can become separate runtimes later without rewriting domain ownership."
                    "webui.boundary" "Web UI is a separate interactive client in V1. Backend remains the sole owner of trusted identity attachment, governance and authoritative state changes."
                    "contract.TypedRequest" "INTENT is data: proposal plus server-owned principal, acting agent, delegation, mission and authz reference; request/operation identities and full audit digest."
                    "contract.TypedRequest.owner" "Entry owns trusted identity attachment; each producing module owns its typed proposal; recipients validate the applicable contract."
                    "contract.identity" "Stable logical operation identity across retries; separate attempt identity. A parameter digest is not an operation identity."
                    "contract.InteractionTurn" "Authenticated user turn classified as QUESTION, CLARIFICATION_RESPONSE, PLAN_CONFIRMATION, STEERING_REQUEST, STOP_REQUEST, RESUME_REQUEST, REPORT_REQUEST or APPROVAL_DECISION as applicable. Classification does not itself grant authority."
                    "contract.InteractionOutput" "Operator-facing reasoning/control outputs use typed contracts such as ClarificationRequest, MissionPlanProposal, SteeringProposal and UserQueryResponse; free-form rendering is presentation only."
                }

                group "Control/Access and trusted identity" {
                    entry = component "Entry Adapters and Authentication" "REST API and MCP HTTP adapters; authenticate human/service principals, classify interactive turns, enforce entry permissions, attach server-owned identity/delegation and route requests to their owning control/read modules. No work dispatch or policy decision at entry." "Deterministic module" {
                        properties {
                            "source.nodes" "API MCP AUTH"
                            "source.contract" "INTENT identity envelope plus authenticated InteractionTurn envelope"
                            "interaction.routing" "Routes conversational mission turns to Mission Controller, visualization/query requests to Authorized Read Layer, reviewer decisions to Approval Registry, and explicit human scoped STOP/RESUME requests to Policy."
                            "authority" "User text, UI state and client-supplied role/scope claims are untrusted inputs. Entry never converts them directly into execution permission or canonical facts."
                        }
                    }
                }

                group "Control/Governance - deterministic" {
                    policy = component "Authority and Policy" "Single current authority matrix and authz revision. Deterministically decide mission, phase, request, scoped stop/resume and export authorization; return decision, scope, reason and retryability. Policy decides only; routing/transition is owned elsewhere. Hard boundaries and authoritative dependency failures fail closed." "Deterministic module" {
                        properties {
                            "source.nodes" "AUTHORITY POLICY"
                            "ownership" "Owns policy semantics and the shared authority source. Call-site gates enforce the returned policy at the time of use."
                            "decision.boundary" "Policy decides ALLOW/DENY/PENDING-style outcomes; Authorized Transition Router performs allowed transitions. Policy does not directly advance workflow."
                            "denial.behavior" "Request denial or budget exhaustion blocks/fails the affected work but does not automatically trigger PHASE or MISSION STOP. A separate explicit scoped-stop rule/decision is required."
                            "availability" "If current authority, ledger, approval or another required authoritative dependency cannot be validated, the action fails closed."
                        }
                    }
                    ledger = component "Scope and Budget Ledger" "Sole accounting writer: atomic reservations and idempotent settlement; cumulative scope, attempts and recovery caps. Available tokens/cost = limit - consumed - reserved. Exhaustion blocks new work but does not itself trigger scoped STOP." "Transactional module" {
                        properties {
                            "source.nodes" "LEDGER"
                            "ownership" "Accounting rules and mutations; persistence in Transactional State Store. No authorization decisions or scheduling."
                            "exhaustion" "Budget/scope/attempt exhaustion yields a blocked or denied work outcome; it triggers neither PHASE nor MISSION STOP unless a separate explicit scoped-stop rule/decision says so."
                        }
                    }
                    approvals = component "Approval Registry" "Canonical frozen-contract digest, logical-operation deduplication, one pending review, recorded denials, expiry, revocation, reviewer and use limits. Approval binds the exact frozen operation; fixed batches enumerate exact operations rather than granting open-ended authority. Persists decisions before delivery." "Transactional module" {
                        properties {
                            "source.nodes" "APPROVAL"
                            "source.interaction" "HITL is the Authorized Reviewer interacting through authenticated UI; this module owns review state."
                            "approval.scope" "Approval covers the exact security-relevant frozen contract. Material tool/adapter/parameter/precondition/destination changes require a new authorization decision."
                            "batching" "A batch approval is permitted only for a fixed enumerated operation list; it is not an open-ended phase grant."
                            "interaction.boundary" "Clarification answers, ordinary conversational responses and plan confirmation are mission-interaction state, not Approval grants. Approval Registry is used only when current Policy requires governance approval."
                        }
                    }
                }

                group "Control/Workflow control - deterministic" {
                    mission = component "Mission Controller" "Owns mission goal, normalized scope, interactive intake/planning state and normal lifecycle. It can request bounded reasoning for clarification/planning, persists interaction state, and submits confirmed plans/steering to Governance. Scoped stop/fence semantics and recovery episodes remain owned by Scoped Stop and Recovery Controller." "Deterministic module" {
                        properties {
                            "source.nodes" "MISSION"
                            "interaction.lifecycle" "Mission intake may move through DRAFT, AWAITING_INPUT, PLANNING, AWAITING_CONFIRMATION and READY before RUNNING. These interaction states are durable and do not themselves grant execution authority."
                            "clarification" "Missing scope/context causes a typed ClarificationRequest rather than guessed values or ApprovalRequested. Clarification is ordinary operator interaction unless Policy separately requires governance approval."
                            "scope.intake" "Operator-supplied IPs, targets, exclusions and constraints are input claims until normalized/validated and checked by current Policy; confirmation does not expand authority."
                            "plan.confirmation" "Operator PlanConfirmation accepts a proposed mission/phase plan for further governed transition. It is not an approval grant for diagnostic operations."
                            "reasoning" "Mission Intake / Planning reasoning is delegated through Scheduler -> Context Broker -> LLM Call Gate; Mission Controller never calls a model directly."
                        }
                    }
                    phases = component "Phase Controllers" "Own lifecycle of intake, analysis and reporting phases, authorized task contracts and deterministic completion criteria. Delegate reasoning to the runtime. A single task failure does not automatically fail the phase; phase outcome follows explicit criticality/completion rules." "Deterministic modules" {
                        properties {
                            "source.nodes" "PHASE"
                            "completion" "Phase completion/failure is derived from explicit deterministic criteria over committed task/fact state, not from an agent saying the phase is done."
                            "task.failure" "Individual task failure is evaluated against phase criticality/completion criteria; it is not an automatic phase failure."
                            "steering" "Phase-specific SteeringProposal from reasoning or operator interaction is treated as a proposal; Phase Controller converts it to task/transition requests only through current workflow/governance rules."
                        }
                    }
                    transitions = component "Authorized Transition Router" "Validate request kind and consume the current authorization decision; coordinate domain-owner mutations in one transaction before routing committed outcomes. Explicit resume only; no permission expansion. It routes decisions but does not make policy decisions." "Deterministic module" {
                        properties {
                            "source.nodes" "ROUTER"
                            "ownership" "Owns authorized state transitions/routing only. It cannot convert a denial, exhausted budget or pending decision into permission."
                        }
                    }
                    scheduler = component "Bounded Scheduler" "Admit already-authorized task contracts with bounded concurrency, fairness and backpressure; consume committed task states, wait on blockers and keep continuations/retries inside the current grant. It schedules work but does not invent new tasks. Produce structured phase summaries." "Deterministic module" {
                        properties {
                            "source.nodes" "SCHEDULER"
                            "source.contract" "DIGEST is a summary value: verified results, unknowns, blockers, evidence references and remaining budget."
                            "task.creation" "Scheduler never creates a new task from its own reasoning. New work originates from Mission/Phase Control or a governed proposal path before scheduling."
                            "retry" "Retries are bounded by retry contract, attempts, budget, deadline, no-progress and current policy. Same logical work keeps its logical identity while a permitted retry receives a distinct attempt identity."
                            "restart" "Durable committed state is the recovery source after process restart; an in-memory queue is not authoritative."
                        }
                    }
                    progress = component "Task Progress Reducer" "Pure deterministic completion, no-progress and state-transition functions over committed events/facts; task revisions and consumed-event deduplication. Verification-dependent completion uses committed verifier facts. No planning, authorization, verification or LLM self-completion." "Deterministic module" {
                        properties {
                            "source.nodes" "PROGRESS"
                            "completion" "Task completion is based on deterministic done_when rules over committed state/facts; agent/LLM claims of completion are not authoritative. A task may require multiple specific valid facts, not a single generic success flag."
                            "verification.wait" "If required verification is delayed or INCONCLUSIVE, task remains pending/blocked according to deterministic rules; raw evidence is never used as a temporary bypass."
                            "no.progress" "No-progress is bounded and cannot be reset by wording changes or unrelated graph updates. Exhausted retry/no-progress limits produce blocked/failed state rather than infinite retry."
                        }
                    }
                }

                group "Control/Context and model admission" {
                    broker = component "Context Broker" "Check current task delegation and cache access; assemble relevant, fresh, token-bounded snapshots. Bind caches to principal, task scope and current authz revision, including report context. Context is granted to a task/profile; workers do not self-discover unrestricted state." "Deterministic module" {
                        properties {
                            "source.nodes" "BROKER"
                            "ownership" "Control-side context authorization, minimization and cache checks; not owned by Analysis Workers."
                            "context.scope" "Workers receive only task-relevant authorized context. No direct unrestricted SQL/graph/database access is granted to reasoning workers."
                            "cache.revocation" "Cache entries are bound to principal, scope and authz revision and must be invalidated/rechecked before reuse after authorization changes; TTL alone is insufficient."
                        }
                    }
                    llmGate = component "LLM Call Gate" "All model calls pass through one governed admission path with atomic reservation, durable call identity and valid admission fence. Own deadlines, cancellation, provider/model configuration, observed usage and settlement requests." "Deterministic runtime gate" {
                        properties {
                            "source.nodes" "LLM_GATE"
                            "ownership" "Authoritative model-call admission/accounting and provider/model selection remain control-side in V1."
                            "stop.scope" "PHASE STOP cancels/fences only model calls belonging to the affected phase/recovery scope; MISSION STOP cancels/fences normal model calls across the mission."
                            "provider.model" "Workers do not choose arbitrary providers/models. The system/operator configures the internal LLM endpoint/model; an explicitly selected alternate model remains a control-side configuration decision."
                            "credentials" "LLM API credentials are system-owned and are never exposed as worker authority or task context."
                            "unknown.usage" "Timeout/unknown provider outcome never implies zero usage; unresolved usage goes to Reconciliation."
                        }
                    }
                }

                group "Reasoning" {
                    analysis = component "Analysis Worker Pool" "Fixed named reasoning profiles with bounded instances, task-scoped skills/context and bounded attempts. Produce analysis, claims, typed proposals and report drafts only. No self-spawn, direct tools, unrestricted state access, model/provider selection or authority changes." "Bounded reasoning workers" {
                        tags "Reasoning"
                        properties {
                            "source.nodes" "AGENTS"
                            "packaging.v1" "Runs inside Backend Application; no independent Reasoning Runtime in V1."
                            "profiles" "Core reasoning roles/profiles are defined by the system and include Mission Intake / Planning, Inventory / Analysis, Evidence Review and Report. Worker instances are bounded executions of those profiles; workers do not dynamically invent/spawn new agent roles."
                            "skills" "Only skills relevant and authorized for the current task/profile are attached. Skills guide reasoning but never grant tool, state-write or authorization privileges."
                            "memory" "Task-scoped working memory/context is allowed. Private long-term worker memory is not authoritative; durable knowledge must enter governed canonical/derived state with provenance."
                            "context" "Workers do not directly query unrestricted databases/read models. Additional context must flow through the authorized read/context path."
                            "task.creation" "Workers may propose new work through a typed request, but cannot insert tasks directly into Scheduler."
                            "completion" "Workers may report claims/results but cannot authoritatively mark a task complete; deterministic Progress owns completion."
                            "outputs" "Any output that can influence control flow must use a typed contract/schema with provenance. Mission interaction uses typed ClarificationRequest, MissionPlanProposal, SteeringProposal and UserQueryResponse outputs; free-form text is not parsed as implicit authority."
                            "disagreement" "Multiple workers may disagree; majority vote or worker confidence does not create authority. Policy, Verification and Progress retain their respective decision ownership."
                            "model.selection" "Workers do not choose arbitrary provider/model. Model selection and credentials are controlled by system/operator configuration through LLM Call Gate."
                            "no.progress" "Workers may propose a changed strategy only within bounded attempts/budget/deadline/no-progress rules; they cannot continue reasoning indefinitely."
                            "future.extraction.seam" "Task input is scoped context/version plus deadline/cancellation; output is proposal/analysis/report content with provenance. Model invocation remains behind the governed LLM-call port."
                            "authority" "No direct tools, no authority changes, no unrestricted state access, no self-certified usage and no direct verified-fact publication."
                        }
                    }
                }

                group "Control/Capability resolution and execution control" {
                    capabilities = component "Capability Registry" "Approved diagnostics and registered compensations; schemas, exact adapter versions, effects and deterministic retry contracts." "Registry module" {
                        properties {
                            "source.nodes" "REGISTRY"
                            "ownership" "Control-side metadata and exact adapter/version contracts; execution loads only the committed implementation."
                            "selection" "Agents propose capability/target intent; deterministic resolver rules/configuration select the exact registered adapter/version. Unregistered tools are never an implicit fallback."
                            "retry.contract" "Carries retry/idempotency/effect metadata used by control-side retry eligibility; Runner does not invent retry policy."
                        }
                    }
                    resolver = component "Capability Resolver" "Deterministically resolve an agent/control capability-and-target proposal to the exact registered adapter, version and security-relevant parameters, then freeze that operation before approval. Unsupported or unavailable capabilities produce typed outcomes; no arbitrary executable/tool fallback." "Deterministic module" {
                        properties {
                            "source.nodes" "RESOLVER"
                            "authority" "Agents do not choose the final tool/adapter/version and cannot bypass the Registry."
                            "material.change" "After authorization, a material tool/adapter/parameter/precondition change requires a new authorization decision; the system does not use LLM judgement to label such a change insignificant."
                        }
                    }
                    executionGate = component "Execution Gate" "Final deterministic recheck immediately before dispatch: exact frozen contract, current permissions/grant, budget/limits, evidence/context freshness, material preconditions and current normal/recovery admission. Fail closed." "Deterministic gate" {
                        properties {
                            "source.nodes" "GATE"
                            "freshness" "Stale or materially changed evidence/context/preconditions block dispatch until refreshed and, when material, re-authorized."
                            "stop.recovery" "PHASE STOP freezes normal diagnostic admission only for the affected phase scope and blocks dependent work; independent phases may continue. MISSION STOP freezes normal diagnostic admission across the mission. Only explicitly admitted recovery operations may run inside an active stopped scope."
                        }
                    }
                    dispatch = component "Atomic Dispatch Transaction" "One transaction claims operation and grant use, reserves quota, and writes dispatch record/outbox. A commit acknowledgement is required before worker claim; no successful commit means no execution." "Transactional module" {
                        properties {
                            "source.nodes" "COMMIT"
                            "ownership" "Coordinates ledger and approval mutations in one transaction; does not duplicate their domain rules."
                            "execution.precondition" "Runner may begin only from a durably committed dispatch under the current fence. Transaction failure or rollback cannot authorize execution."
                        }
                    }
                    reconciliation = component "Outcome and Reservation Reconciliation" "Classify and resolve LLM calls, adapter attempts, export/release attempts, usage/effects and orphan holds through bounded, deadline-bound automatic investigation. Unknown outcomes never imply zero usage/effect/delivery and never trigger blind retry or blind re-release. Retry eligibility is deterministic and re-enters the appropriate governed admission path; unresolved ambiguity escalates only after bounded automatic handling is exhausted." "Deterministic module" {
                        properties {
                            "source.nodes" "RECONCILE"
                            "ownership" "Control-side outcome classification and supported settlement/retry-eligibility decisions; Ledger performs accounting mutation. Runner does not own settlement or retry authority."
                            "outcome.classification" "Distinguish proven/known outcomes from unknown outcomes and whether a retry is deterministically safe under the frozen operation contract; a timeout alone is not proof that execution did not occur."
                            "automatic.reconciliation" "Before MANUAL_REQUIRED/UNRESOLVED, perform bounded automatic investigation within reconciliation attempt/deadline/budget limits. Additional diagnostics produce a durable FollowUpRequested with reason RECONCILIATION_DIAGNOSTIC for the owning Mission/Phase Control or active Recovery episode; accepted proposals re-enter Resolver -> Policy -> Router -> Gate -> Dispatch."
                            "retry.eligibility" "Retry requires the frozen capability retry contract plus current policy, budget, deadline, no-progress and fence checks. Same operation_id is retained and a permitted retry gets a new attempt_id."
                            "llm.role" "LLM/Analysis may propose what to try next, but cannot decide that an external effect is safe to retry or authorize a retry."
                            "late.results" "Late observations/results remain valid evidence/reconciliation input, but a stale attempt cannot overwrite a newer authoritative outcome or settle the same reservation/release attempt twice."
                            "release.reconciliation" "Unknown, partial or unacknowledged export delivery is classified before any repeat release. A retry proposal can proceed only after deterministic eligibility and a fresh governed Policy/Export path; Reconciliation never sends an artifact itself."
                            "manual.last.resort" "MANUAL_REQUIRED/UNRESOLVED is used when bounded reconciliation cannot establish a safe automated continuation, retry/reconciliation limits are exhausted, a material contract change needs new authority, or policy explicitly requires human review."
                        }
                    }
                }

                group "Execution" {
                    adapters = component "Bounded Adapter Workers" "Claim committed work under the current admission fence; execute approved diagnostics or registered compensation only. Bound concurrency, deadlines and cancellation." "Bounded adapter workers" {
                        tags "Execution"
                        properties {
                            "source.nodes" "RUNNER"
                            "packaging.v1" "Runs inside Backend Application; no independent Execution Runtime in V1."
                            "execution.port" "Consumes only committed work through a narrow execution interface; never receives a mutable mission object, unrestricted repository/session or control service."
                            "committed.work.contract" "Stable operation_id, distinct attempt_id, committed dispatch reference, target/scope, pinned adapter/version/config, input/evidence refs, deadline and admission/fence reference."
                            "cancellation.fencing" "Cancellation is scoped and observable; stale admission claims cannot authorize new work. Cancellation cannot erase an already-sent external effect."
                            "result.evidence.contract" "Returns execution status/observations/evidence with operation/attempt identity and provenance; execution status is not a verified verdict."
                            "unknown.outcome" "Crash/timeout after a possible effect becomes durable unknown outcome for reconciliation; missing acknowledgement never triggers blind effect retry."
                            "redelivery" "Redelivery of the same attempt_id is deduplicated/idempotent and does not by itself execute the external effect again. A valid retry is a new attempt_id under the same logical operation_id."
                            "retry.authority" "Runner never creates its own retry attempt. It reports status/usage/evidence; deterministic control-side rules decide whether another attempt may be scheduled."
                            "state.access" "Use narrow work/admission/result/evidence ports. No direct ownership of approval, budget, workflow, reconciliation or verified-fact state."
                            "future.extraction.seam" "Control logic depends on an execution abstraction implemented by Runner. No RPC, broker, service discovery or remote-worker protocol is introduced in V1."
                        }
                    }
                }

                group "Control/Evidence intake" {
                    intake = component "Evidence Intake" "Validate evidence schema, provenance, task/operation/attempt identity and integrity metadata; persist immutable evidence. Agent claims, tool outputs and external content remain untrusted observations until Verification adjudicates them." "Deterministic module" {
                        properties {
                            "source.nodes" "INGEST"
                            "ownership" "Accepts observations as evidence only; never promotes reasoning/execution claims directly to verified facts or task completion."
                            "completion.boundary" "Raw evidence is never a substitute for a required verified FACT, even when it comes from a normally trusted adapter."
                            "verification.handoff" "Own durable intake receipts keyed by evidence_id before artifact writes. After bytes are durable, atomically commit evidence reference/hash and VerificationRequested job/outbox in Transactional State Store, then acknowledge acceptance. Resume pending receipts after restart; repeated intake is idempotent."
                        }
                    }
                }

                group "Verification" {
                    verifier = component "Verifier Registry and Bounded Workers" "Independently adjudicate scoped immutable evidence using verification rules; no direct target/tool calls. Produce VERIFIED, INVALIDATED or INCONCLUSIVE outcomes with provenance, observed time, validity and lifecycle. Commit FACT/lifecycle events before completion use." "Independent verification workers" {
                        tags "Verification"
                        properties {
                            "source.nodes" "VERIFY"
                            "source.contract" "FACT is a validated fact/lifecycle event, not an application or data store. Authorized fact publication belongs only to Verification."
                            "packaging.v1" "Runs inside Backend Application; logical/adjudication independence is required, but separate process/credentials are deliberately deferred."
                            "future.extraction.seam" "verification_job_id binds task/scope, immutable evidence refs/version and verification rule/version; verdict/provenance returns through a fact-publication interface. Execution has no valid business path to publish FACT."
                            "job.lifecycle" "Own bounded claims and restart recovery of committed verification jobs. Job completion and FACT/INCONCLUSIVE plus any FollowUpRequested commit atomically; job identity and claim revision prevent duplicate publication or stale-worker completion."
                            "authority" "Execution status, adapter success=true, agent confidence or LLM judgement is never sufficient to publish a verified fact."
                            "no.direct.probe" "Verifier has no direct external tool/target capability in V1. More evidence produces a durable FollowUpRequested with reason EVIDENCE_NEEDED for the owning Mission/Phase Control or active Recovery episode; it grants no execution authority."
                            "inconclusive" "INCONCLUSIVE is a first-class result. It does not become VERIFIED/INVALIDATED by force and does not automatically fail or retry the task."
                            "conflict" "Conflicting evidence remains explicit in provenance/verdict state. An LLM may assist analysis but cannot choose a winning source as verdict authority."
                            "history" "Verified-fact history is append/lifecycle based; newer evidence changes lifecycle/validity through new events rather than overwriting/deleting prior provenance."
                            "staleness" "A stale/expired fact remains historical evidence but cannot satisfy a current completion/authorization rule that requires fresh validity."
                            "retry" "INVALIDATED/INCONCLUSIVE does not directly trigger execution retry. Workflow/Policy/Scheduler decide any governed follow-up."
                        }
                    }
                }

                group "Control/Canonical state access and derived views" {
                    projections = component "Projection Workers" "Replay committed canonical events into idempotent SQL, knowledge-graph and optional observability projections with explicit lag/watermarks. Budget projections are read-only accounting views." "Bounded projection workers" {
                        properties {
                            "source.nodes" "PROJECT"
                            "source.of.truth" "Projection input is committed canonical events only. Projection state never becomes a substitute source of canonical truth."
                            "write.boundary" "Projection workers write derived representations only; they never mutate canonical events/state, Ledger, approvals, grants, workflow authority, stop state or verified facts."
                            "rebuild" "Derived representations are replayable/rebuildable from canonical history. Loss or corruption of a projection does not require inventing canonical state from the projection."
                            "ordering" "Require deterministic/idempotent processing and ordering for the relevant entity/stream/projection contract. A global total order is not required by default."
                        }
                    }
                    readAccess = component "Authorized Read Layer" "Enforce resource, row, field and subgraph access under current principal/scope and the shared authority source. Expose typed/versioned SQL, knowledge-graph and observability views with freshness metadata." "Deterministic query module" {
                        properties {
                            "source.nodes" "VIEWS - read enforcement"
                            "ownership" "Owns data-level read enforcement for both human/UI and reasoning consumers; Broker additionally owns task delegation, minimization and cache checks for reasoning."
                            "authorization" "Every protected read is checked against current principal/scope/authz revision; stale cached authorization cannot be used merely because its TTL has not expired."
                            "ui.boundary" "Web UI never queries SQL/graph stores directly. Interactive visualization/query requests enter through Entry and this layer."
                        }
                    }
                    telemetry = component "Observability Consumers and Views" "Consume non-critical events asynchronously with filtering, deduplication, retries and bounded queues. Derive usage, latency, queue/lag, recovery/holds/stop and approval-prompt/wait views." "Asynchronous consumers" {
                        properties {
                            "source.nodes" "OUTBOX AUDIT"
                            "ownership" "OUTBOX here means consumers, not a second transactional outbox. Audit views are outputs, not canonical state or workflow control."
                        }
                    }
                }

                group "Control/Scoped stop and bounded recovery" {
                    recovery = component "Scoped Stop and Recovery Controller" "Own PHASE- and MISSION-scoped stop epochs, cancellation/admission fences, unique bounded recovery episodes, deterministic escalation and terminal recovery status. Reconcile uncertain work before effect-based registered compensation. No automatic resume or reentry." "Deterministic state machines" {
                        properties {
                            "source.nodes" "STOP CANCEL CLEANUP"
                            "source.contract" "RECOVERED carries recovery_status (COMPLETE, PARTIAL or MANUAL_REQUIRED), outstanding work and terminal_reason; handover_episode_id identifies the mission recovery episode on takeover."
                            "ownership" "Owns scoped stop/fence, cancellation and recovery-episode semantics through distinct internal handlers; Mission/Phase Controllers own their lifecycle. Human STOP/RESUME joins the Router-coordinated transaction; no direct compensation execution."
                            "scope" "STOP scope is PHASE or MISSION. A phase stop freezes/cancels work in that phase only; a mission stop freezes normal work across the mission."
                            "dependencies" "Independent phases may continue during a phase stop. A phase whose required dependency is stopped/missing becomes BLOCKED rather than bypassing the dependency."
                            "recovery.scope" "A phase stop creates at most one bounded recovery episode for that phase/stop epoch. Automatic compensation is limited to confirmed effects inside the stopped scope unless a separately authorized request expands authority."
                            "scope.overlap" "MISSION STOP atomically fences phase recovery admission. After bounded reconciliation, atomically transfer remaining effects/history and close the phase episode as PARTIAL, terminal_reason=SUPERSEDED_BY_MISSION_STOP, handover_episode_id=mission_recovery_episode_id. Each effect has one compensation owner checked at dispatch; phase recovery stays fenced while mission STOP is active."
                            "escalation" "Phase STOP may escalate to Mission STOP only through an explicit deterministic escalation rule/decision; escalation is never implicit from ordinary failure."
                            "resume" "AuthorizedResume from Router binds the expected scope/stop epoch. Recovery validates current stop state and terminal recovery, then advances the fence in the same transaction as the owning controller's lifecycle update. Old claims remain invalid; counters are preserved. Phase resume is blocked by active mission STOP; mission resume preserves explicit phase stops."
                            "idempotency" "Repeated STOP for the same scope and stop epoch is idempotent and cannot create another recovery episode."
                            "hard.denial" "Hard denial ends that operation path, not automatically the episode. Never rewrite/resubmit materially equivalent compensation to bypass Policy. While the episode is active and within limits, a genuinely different compensation may be a new TypedRequest through Governance from the start."
                            "closure" "Recovery Controller closes an episode on completion, explicit episode-level termination (including mission takeover), episode-limit exhaustion, or no valid automatic continuation. Operation-level denial alone is not episode termination. COMPLETE requires no outstanding recovery work; PARTIAL/MANUAL_REQUIRED preserve outstanding work and any handover reference."
                            "manual.last.resort" "Recovery uses bounded automatic reconciliation/compensation first. MANUAL_REQUIRED is a terminal fallback only when safe automated continuation cannot proceed within authority, attempts, budget or deadline."
                            "late.settlement" "Closing recovery atomically invalidates outstanding recovery admission claims; late settlement/evidence remains processable without reopening the episode."
                        }
                    }
                }

                group "Control/Reporting and controlled export" {
                    reports = component "Report Assembly" "Assemble report content from authorized verified facts, explicit unknowns and reasoning drafts under a pinned operator/system-selected report template. Validate required sections/fields and provenance before producing the report artifact candidate. Draft generation and template conformance do not verify facts." "Deterministic assembly module" {
                        properties {
                            "source.nodes" "REPORT"
                            "template.selection" "Report template is selected by Operator, mission configuration or system policy; Report/LLM workers cannot choose or silently switch templates."
                            "template.contract" "Template is a presentation/schema contract, not a fact source or authority source. Required factual fields must be backed by authorized verified data or remain explicitly UNKNOWN/UNVERIFIED/INCONCLUSIVE as allowed by the template."
                            "template.version" "template_id and template_version are pinned for the report build and carried forward into artifact metadata and release records."
                            "template.validation" "Assembly validates required sections, field types, evidence/provenance references and template-specific constraints before redaction/freeze."
                            "severity.authority" "Authoritative severity/risk values must come from an authorized assessment contract or deterministic/domain rule over authorized facts. LLM/Report workers may propose a rating/rationale but cannot publish the authoritative severity."
                            "fact.boundary" "Assembly cannot create, upgrade, downgrade or overwrite verdicts/facts to satisfy a template."
                        }
                    }
                    redaction = component "Output Redaction and Validation" "Apply sensitive-field, secret, destination and final template-output validation rules; then freeze the exact artifact bytes and digest before requesting export authorization." "Deterministic module" {
                        properties {
                            "source.nodes" "REDACT"
                            "template.integrity" "Final rendered output must still conform to the pinned template/version after redaction. Any material content/template change produces a new artifact digest and requires fresh authorization."
                            "freeze.metadata" "Frozen artifact metadata includes template_id and template_version together with the exact artifact digest."
                            "freeze.durability" "The exact frozen artifact bytes and metadata must be durably recoverable before approval can become pending. Physical storage placement/technology remains deliberately deferred."
                        }
                    }
                    exports = component "Export Gate" "Recheck current permissions, pinned template metadata, artifact digest, destination and admission. Atomically claim grant use and commit the release record; release only the exact frozen bytes under the valid acknowledged release contract." "Transactional gate" {
                        properties {
                            "source.nodes" "EXPORT"
                            "source.contract" "OUTPUT is the released Approved Report artifact, not a service or additional data store."
                            "release.binding" "Release authorization/record binds artifact_digest + template_id + template_version + destination. A material byte/content/template/destination change requires a new authorization decision."
                            "release.idempotency" "The same release_attempt_id is deduplicated/idempotent. Unknown or unacknowledged delivery is reconciled before any new release attempt; Export Gate never blindly resends."
                        }
                    }
                }
            }

            group "Authoritative persistence" {
                transactions = container "Transactional State Store" "Canonical append-only events and transactional outbox, ledger holds/counters, approval/use claims, dispatch/release records, evidence intake receipts/references, verification jobs, admission fences and PHASE/MISSION-scoped stop/recovery identities and effect ownership share one transaction boundary." "ACID-capable store; engine unspecified" {
                    tags "Data Store,Authoritative"
                    properties {
                        "source.nodes" "EVENTS"
                        "source.persistence" "LEDGER APPROVAL COMMIT STOP CANCEL CLEANUP EXPORT LLM_GATE INGEST VERIFY"
                        "invariant" "Events stay append-only; associated current-state records are transactional. Durable state is authoritative; projections never authorize."
                    }
                }
                evidence = container "Protected Evidence Store" "Immutable artifacts and integrity hashes; authorized reads and secret references. Storage of evidence does not establish its truth." "Protected artifact store; implementation unspecified" {
                    tags "Data Store,Authoritative"
                    properties {
                        "source.nodes" "RAW"
                    }
                }
            }
            group "Derived state" {
                readModels = container "Derived Read Models" "Replayable typed/versioned derived representations built from committed canonical events, including SQL and knowledge-graph projections plus optional materialized observability/audit views, with explicit freshness/lag information." "Derived projections; physical engines unspecified" {
                    tags "Data Store,Derived"
                    properties {
                        "source.nodes" "VIEWS - materialized data; AUDIT only if materialized"
                        "boundary" "Logical derived-state boundary only. SQL and knowledge-graph representations are not asserted to share one physical database, server or engine."
                        "freshness" "Every projection exposes version/watermark/lag metadata. A stale projection may support analysis only when the task freshness contract permits it."
                        "graph.role" "Knowledge Graph projections are optimized for relationship traversal, asset/entity linkage, evidence/finding provenance traversal, dependency analysis and authorized task-scoped subgraph retrieval."
                        "authority" "Fresh or stale projections are never authoritative for execution, grants, reservations, budget mutation, current stop epoch or verified-fact publication."
                        "rebuild" "This boundary is replayable/rebuildable from canonical history and must not contain unique authoritative facts that cannot be reconstructed."
                    }
                }
            }
        }

        // C2 RELATIONSHIPS. Explicit aggregation; no inferred bypass paths.
        // Web UI and CLI are presentation/access clients. Backend remains the
        // authority boundary for identity, governance and state transitions.
        operator -> web "Converses with the system; supplies scope/context, confirms plans, steers work, requests scoped stop/resume, selects permitted report templates and receives approved reports" "User interaction" "C2"
        operator -> cli "Submits goals and mission controls" "User interaction" "C2"
        reviewer -> web "Reviews exact frozen contracts and material changes" "User interaction" "C2"
        web -> backend "Authenticated conversational turns, mission inputs/controls, review decisions, report-template selection, visualization queries and authorized results" "HTTP; backend API binding" "C2"
        cli -> backend "Authenticated requests and mission controls" "Client binding unspecified" "C2"
        backend -> transactions "Commits canonical state atomically; consumes typed committed events" "Transactional reads/writes and subscriptions; binding unspecified" "C2"
        backend -> evidence "Writes immutable evidence and performs authorized scoped reads" "Artifact store interface; binding unspecified" "C2"
        backend -> readModels "Builds replayable projections; reads through current authorization" "Derived read-model access; binding unspecified" "C2"

        // ENTRY. API and MCP are protocol adapters in one Backend Application,
        // not independent services. Web UI and CLI are client-side entry paths.
        web -> entry "Authenticated conversational turns, scope/context input, clarification answers, plan confirmations, steering requests, scoped stop/resume, report requests/template selection, visualization queries and review decisions" "HTTP; backend API binding" "Detail"
        cli -> entry "Submits requests for authentication" "Client binding unspecified" "Detail"
        policy -> entry "Required current entry and mission-control permissions" "In-process" "Check"
        entry -> mission "Authenticated mission creation or InteractionTurn: question, scope/context input, clarification response, plan confirmation, steering/report request; no work dispatched at entry" "In-process" "Detail"
        entry -> policy "Explicit human PHASE/MISSION STOP or RESUME TypedRequest with trusted envelope and expected scope/stop epoch" "In-process" "Detail"
        entry -> approvals "Authenticated reviewer decision for the frozen operation" "In-process" "Detail"
        mission -> entry "Typed operator-facing ClarificationRequest, MissionPlanProposal, Steering/UserQuery response or current interaction state; no hidden authority in rendered text" "In-process" "Detail"
        entry -> web "Authorized conversational/query response and current mission interaction state" "HTTP response/stream; binding unspecified" "Detail"
        mission -> transactions "Persist mission interaction lifecycle, normalized input references and operator confirmations; interaction state never grants execution authority by itself" "Durable append / transactional state" "Detail"

        // REQUESTS. INTENT is carried by these edges. Trusted envelopes are
        // always server attached, including proposals made by reasoning workers.
        mission -> policy "Confirmed mission/phase/steering TypedRequest after interaction; confirmation is not authority and current Policy still decides" "In-process" "Detail"
        phases -> policy "Phase change or governed task/steering TypedRequest" "In-process" "Detail"
        mission -> resolver "Diagnostic TypedRequest from an accepted mission-owned FollowUpRequested; preserve originating work limits" "In-process" "Detail"
        phases -> resolver "Diagnostic TypedRequest from an accepted phase-owned FollowUpRequested; preserve originating work limits" "In-process" "Detail"
        analysis -> resolver "Proposed capability/target TypedRequest only; Resolver owns exact registered tool/adapter/version selection" "In-process" "Detail"
        recovery -> resolver "Registered compensation or accepted follow-up diagnostic TypedRequest for the active scoped bounded recovery episode; no direct tool call" "In-process" "Detail"
        redaction -> policy "Export TypedRequest with frozen artifact digest, pinned template_id/version and destination" "In-process" "Detail"
        capabilities -> resolver "Required registered schema, deterministic adapter/version selection metadata, effects and retry contract" "In-process" "Check"
        resolver -> policy "Frozen exact operation contract/version; material change requires new authorization" "In-process" "Detail"
        resolver -> transactions "Unsupported or unavailable typed request outcome" "Durable append" "Detail"

        // POLICY AND REVIEW. Approval Registry owns review lifecycle. Web UI
        // presents the frozen review contract; Authorized Reviewer is the human,
        // not another backend component or authority source.
        ledger -> policy "Required authoritative cumulative limits" "In-process" "Check"
        policy -> approvals "Hard checks passed; obtain or request human approval" "In-process" "Detail"
        approvals -> web "Present one pending frozen-contract review with changes and decision context" "Authorized UI response/notification" "Detail"
        approvals -> policy "Bound recorded decision; denial/expiry cannot automatically re-prompt" "In-process" "Detail"
        approvals -> transactions "Commit approval lifecycle and grants before decision delivery" "Shared transaction" "Detail"
        policy -> transitions "Current decision allows the frozen request" "In-process" "Detail"
        policy -> transactions "Request denial, pending, budget-blocked or local dependency failure; these do not automatically trigger PHASE/MISSION STOP" "Durable append when available" "Detail"
        policy -> recovery "Scoped recovery-operation authorization outcome or separate recovery-episode termination decision; human STOP/RESUME goes through Router" "In-process; distinct typed outcomes" "Detail"
        transitions -> recovery "Authorized scoped STOP/RESUME transition; Recovery applies current stop/fence semantics inside the shared lifecycle transaction" "In-process / shared transaction" "Detail"
        transitions -> phases "Authorized phase lifecycle mutation, including scoped STOP/RESUME, joins transition transaction with Recovery fence update when applicable; notify only after commit" "In-process / shared transaction" "Detail"
        transitions -> mission "Authorized mission STOP/RESUME lifecycle mutation joins Recovery fence update in one transition transaction; notify only after commit" "In-process / shared transaction" "Detail"
        transitions -> executionGate "Approved frozen diagnostic or compensation operation" "In-process" "Detail"
        transitions -> exports "Approved frozen export" "In-process" "Detail"
        transitions -> transactions "Atomically commit domain-owner mutations and transition outcome before routing; scoped STOP/RESUME includes lifecycle and fence" "Shared transaction / durable append" "Detail"

        // REASONING AND CONTEXT. Data access checks and delegation/cache checks
        // are separate enforcement duties with one shared policy source.
        mission -> scheduler "Authorized bounded Mission Intake / Planning reasoning task for question/clarification/plan work; Scheduler does not invent it" "In-process" "Detail"
        phases -> scheduler "Authorized task contract inside the approved phase/scope; Scheduler does not invent tasks" "In-process" "Detail"
        scheduler -> broker "Admitted task or continuation within the current task grant" "In-process" "Detail"
        policy -> readAccess "Required current resource, row, field and subgraph policy" "In-process" "Check"
        entry -> readAccess "Authorized interactive visualization/query request under current principal, mission scope and authz revision" "In-process" "Detail"
        policy -> broker "Required current task delegation, read policy and cache validity" "In-process" "Check"
        readModels -> readAccess "Projected typed/versioned facts and lag watermarks; enforce access before return" "Read query result" "Detail"
        readAccess -> broker "Authorized records only, with projection freshness/watermark metadata" "In-process" "Detail"
        readAccess -> entry "Authorized UI query/visualization result with freshness/watermark metadata; no direct store exposure" "In-process" "Detail"
        broker -> llmGate "Relevant fresh bounded snapshot and authorization revision" "In-process" "Detail"
        broker -> transactions "Denied, stale or unavailable context" "Durable append" "Detail"
        llmGate -> analysis "Admit bounded profile execution after committed call reservation and fence check" "In-process" "Detail"
        analysis -> mission "Mission Intake / Planning output: typed ClarificationRequest, MissionPlanProposal, SteeringProposal or UserQueryResponse; proposal/response only, never direct state transition" "In-process" "Detail"
        analysis -> phases "Phase/task proposal or steering analysis only; owning Phase Controller decides whether to create a governed request" "In-process" "Detail"
        analysis -> llmGate "Runtime completion and provider-usage observations; LLM Call Gate remains authoritative for admission/accounting settlement" "In-process" "Detail"
        llmGate -> ledger "Atomically reserve and settle known call usage" "Shared transaction" "Detail"
        llmGate -> reconciliation "Unknown call outcome or usage" "In-process" "Detail"
        llmGate -> transactions "Durable call identity, lifecycle and typed admission failure" "Shared transaction / durable append" "Detail"

        // PROGRESS. DIGEST is produced by Scheduler and consumed as data.
        // Scheduler schedules authorized contracts only; it does not plan/invent
        // new tasks. Completion/no-progress/retry authority remains deterministic.
        analysis -> intake "Untrusted claims and task-result evidence" "In-process" "Detail"
        transactions -> progress "Relevant committed inputs only; exclude reducer output events" "Typed event subscription" "Detail"
        progress -> transactions "Idempotent TaskStateChanged with reason/severity and task revision check" "Durable append" "Detail"
        transactions -> scheduler "Committed TaskStateChanged/retry-eligibility outcomes only; schedule bounded continuation or retry when current rules allow" "Typed event subscription" "Detail"
        scheduler -> phases "Structured phase summary: verified results, unknowns, blockers, references, budget" "In-process value" "Detail"
        scheduler -> mission "Compact structured phase summary" "In-process value" "Detail"
        scheduler -> recovery "V1 internal deterministic PHASE/MISSION breaker safety path only; ordinary task failure or budget exhaustion is not automatically STOP" "In-process" "Detail"

        // EXECUTION. The same transaction spans claims, ledger, dispatch and
        // outbox; outbox delivery is not a substitute for that atomicity.
        // Runner never self-retries. Unknown outcomes go to bounded reconciliation;
        // a retry candidate is not execution authority and must re-enter governed
        // scheduling/admission with the same operation_id and a new attempt_id.
        policy -> executionGate "Required current authority, permission and policy decision" "In-process" "Check"
        approvals -> executionGate "Required exact grant, expiry, revocation and remaining uses" "In-process" "Check"
        ledger -> executionGate "Required authoritative quota and recovery limits" "In-process" "Check"
        broker -> executionGate "Required snapshot, evidence references, freshness and watermark" "In-process" "Check"
        executionGate -> dispatch "All exact-contract checks pass; no material change" "In-process" "Detail"
        executionGate -> transactions "Normal request blocked or material preconditions changed" "Durable append" "Detail"
        executionGate -> recovery "Scoped recovery-operation denial, blocker or limit outcome; Recovery evaluates episode continuation or closure" "In-process" "Detail"
        dispatch -> ledger "Operation reservation joins the dispatch transaction" "Shared transaction" "Detail"
        dispatch -> approvals "Exact grant-use claim joins the same dispatch transaction" "Shared transaction" "Detail"
        dispatch -> transactions "Atomically commit operation claim, dispatch record and outbox with all claims; compensation also checks current episode/effect ownership and fence" "Shared transaction" "Detail"
        dispatch -> adapters "Commit acknowledged; claim only committed work under the current fence" "In-process dispatch notification" "Detail"
        adapters -> transactions "Submit fenced work-claim and attempt lifecycle through the control-owned execution persistence port; no unrestricted repository access" "In-process port; transactional claim / durable append" "Detail"
        adapters -> ledger "Known actual usage and terminal attempt state" "In-process settlement request" "Detail"
        adapters -> reconciliation "Timeout, uncertain completion or cancellation; Runner does not self-retry" "In-process" "Detail"
        adapters -> intake "Raw diagnostic observations" "In-process" "Detail"
        reconciliation -> ledger "Confirmed idempotent settlement; consume known usage and release only proven unused reserved balance" "In-process settlement request" "Detail"
        reconciliation -> transactions "Commit classified outcome, retry eligibility, pending state, FollowUpRequested (RECONCILIATION_DIAGNOSTIC), or bounded MANUAL_REQUIRED/UNRESOLVED" "Durable append / transactional outbox" "Detail"
        ledger -> transactions "Commit accounting events with ledger mutations; preserve cumulative counters" "Shared transaction" "Detail"
        ledger -> reconciliation "Aged unresolved reservations for bounded investigation; TTL never proves zero usage or zero external effect" "In-process" "Detail"
        // A failed dispatch transaction causes NO dispatch. Its failure outcome
        // can be appended only after rollback and when event persistence works.

        // EVIDENCE. FACT is the verifier's typed durable output. The reasoning
        // Evidence Review profile cannot write verified facts or self-verify.
        intake -> evidence "Store immutable validated-schema artifacts and integrity hashes" "Artifact store interface" "Detail"
        intake -> transactions "Persist pending intake receipt; after durable bytes, atomically finalize evidence reference/hash and VerificationRequested job/outbox before acceptance" "Shared transaction / durable append" "Detail"
        policy -> evidence "Required current evidence and verifier read permissions" "Authorized store access contract" "Check"
        transactions -> verifier "Committed VerificationRequested jobs; claim/recover unfinished jobs by stable identity after restart" "Typed event subscription / durable job reads" "Detail"
        evidence -> verifier "Scoped verifier reads of untrusted evidence" "Authorized artifact reads" "Detail"
        verifier -> transactions "Persist job claims; atomically complete job with authorized FACT/INCONCLUSIVE and optional FollowUpRequested (EVIDENCE_NEEDED) through the verification port" "In-process port / shared transaction" "Detail"

        // FOLLOW-UP. One durable proposal contract; existing control owners decide.
        transactions -> mission "Mission-owned FollowUpRequested; idempotently record blocked/rejected outcome or durable diagnostic request" "Typed event subscription / shared transaction" "Detail"
        transactions -> phases "Phase-owned FollowUpRequested; idempotently record blocked/rejected outcome or durable diagnostic request" "Typed event subscription / shared transaction" "Detail"

        // STOP / CANCEL / CLEANUP are scoped internal state-machine handlers of
        // one owner. RECOVERED is a status, not a service. STOP scope is PHASE
        // or MISSION. Phase stop affects only work in that phase; dependent phases
        // become BLOCKED while independent phases may continue. No automatic resume.
        recovery -> transactions "Atomically persist scoped stop/resume fences, recovery entry/checkpoints, effect ownership/handover and terminal status; human STOP/RESUME joins lifecycle transaction" "Shared transaction / durable append" "Detail"
        recovery -> scheduler "Apply scoped fence: phase STOP freezes/cancels only that phase and blocks dependent phases; mission STOP freezes normal mission work; permit only admitted scoped recovery" "In-process" "Check"
        recovery -> llmGate "Apply current scoped stop/fence; cancel affected calls; mission takeover and terminal recovery invalidate old recovery claims" "In-process" "Check"
        recovery -> executionGate "Require current stop scope/epoch, active recovery identity and exclusive compensation effect ownership; mission STOP fences phase recovery" "In-process" "Check"
        recovery -> adapters "Fence/cancel affected work in the stopped scope; terminal recovery forbids reentry for that scope/epoch" "In-process" "Check"
        recovery -> exports "Enforce current scoped stop state on export admission; mission STOP freezes normal export, while phase-scoped effects must obey their dependencies/current authorization" "In-process" "Check"
        recovery -> ledger "Count scoped recovery entry atomically; preserve cumulative mission/phase counters and release only proven undispatchable holds" "Shared transaction" "Detail"
        recovery -> reconciliation "Reconcile in-flight/uncertain work inside the active stopped scope before compensation or closure" "In-process" "Detail"
        ledger -> recovery "Required persistent scoped recovery entry, attempt, cost and deadline caps; resume never resets cumulative counters" "In-process" "Check"
        transactions -> recovery "Committed outcomes and recovery-owned FollowUpRequested for the active episode; idempotent handling; terminal episodes, including PARTIAL episodes handed over to mission recovery, never restart" "Typed event subscription / shared transaction" "Detail"

        // DERIVED DATA AND NON-CRITICAL OBSERVABILITY.
        transactions -> projections "Committed projection-relevant canonical events only; canonical history is the source of truth" "Typed event subscription / replay" "Detail"
        projections -> readModels "Idempotent/replayable SQL/knowledge-graph updates with entity/stream ordering and lag watermarks" "Projection writes" "Detail"
        transactions -> telemetry "Filtered non-critical events; asynchronous bounded consumption" "Typed event subscription" "Detail"
        telemetry -> readModels "Optional materialized observability/audit views with lag metadata; never canonical workflow state" "Derived projection writes" "Detail"
        // No read-model, graph, audit/telemetry view or projection worker writes Ledger,
        // approvals, grants, verified facts, stop state or canonical workflow state.

        // REPORTING. Profile drafts, verified facts, assembly, redaction and
        // release authorization retain distinct responsibilities.
        analysis -> reports "Report draft/content suggestions only; worker cannot select/switch the pinned report template or create verified facts" "In-process" "Detail"
        broker -> reports "Restricted report context, authorized verified references and pinned operator/system-selected template metadata" "In-process" "Detail"
        reports -> redaction "Template-conformant assembly with verified facts, explicit unknowns, provenance and pinned template_id/version" "In-process" "Detail"
        policy -> exports "Required current export permission and authority" "In-process" "Check"
        approvals -> exports "Required artifact-bound grant, expiry, revocation and remaining uses" "In-process" "Check"
        exports -> approvals "Exact artifact-bound grant-use claim joins release transaction" "Shared transaction" "Detail"
        exports -> transactions "Atomic grant-use and durable release-attempt record binding artifact digest, template_id/version and destination; record denials/changed artifacts/release failures" "Shared transaction / durable append" "Detail"
        exports -> reconciliation "Unknown, partial or unacknowledged external release outcome; never blindly repeat delivery" "In-process" "Detail"
        reconciliation -> policy "Governed release-retry TypedRequest only after deterministic release eligibility; current Policy/approval/export checks still apply" "In-process" "Detail"
        exports -> web "Approved Report: release frozen bytes only after acknowledged commit and valid release claim" "Authorized UI output channel" "Detail"
    }

    views {
        container system "C2-Containers" "V1 modular monolith: interactive Web UI and CLI clients, one Backend Application, and three logical persistence containers. Reasoning, Execution and Verification remain internal backend subsystems with future extraction seams." {
            include operator reviewer web cli backend transactions evidence readModels
            autoLayout tb 260 180
            default
        }

        component backend "C3-BackendOverview" "Orientation view only: shows the interactive Web UI -> mission/reasoning loop and the main governed execution/verification path. It is intentionally non-exhaustive; detailed C3 views carry the full model. Control, Reasoning, Execution and Verification are not separate runtime/process boundaries in V1." {
            // INTERACTIVE ENTRY + CONTROL.
            include web entry mission scheduler broker llmGate policy resolver executionGate dispatch intake evidence transactions

            // SEPARABLE-LATER SUBSYSTEMS.
            include analysis adapters verifier

            autoLayout lr 380 260
        }

        component backend "C3-Governance" "Reviewed V1 governance: Policy decides while Router transitions; clarification/plan confirmation is distinct from Approval; approval binds exact frozen operations; required authority dependencies fail closed; denial or budget exhaustion does not automatically trigger PHASE/MISSION STOP." {
            include reviewer web cli entry mission policy ledger approvals resolver transitions transactions
            autoLayout tb 240 160
        }
        component backend "C3-Workflow" "Reviewed V1 workflow: durable interactive mission intake/plan-confirmation before governed execution, bounded parallel scheduling of already-authorized tasks, deterministic completion/no-progress, bounded retries with distinct attempts, durable restart recovery and no automatic scoped STOP on ordinary task failure." {
            include mission phases policy transitions scheduler progress broker recovery transactions analysis resolver
            autoLayout tb 240 160
        }
        component backend "C3-Reasoning" "Reviewed V1 reasoning: interactive mission questions/planning use a fixed Mission Intake / Planning profile; context is granted/minimized, every model call is gated, profiles/skills are task-scoped and bounded, workers cannot self-spawn/use tools/choose models/create tasks/decide completion, and control-affecting outputs are typed proposals rather than free-form authority." {
            include web entry mission scheduler policy readAccess readModels broker llmGate analysis ledger reconciliation transactions
            autoLayout tb 240 160
        }
        component backend "C3-OperationAuthorization" "Reviewed V1 authorization: agents propose capability/target, Resolver deterministically freezes the exact registered operation, material changes require re-authorization, Gate performs the final current-state recheck, and no committed dispatch means no execution." {
            include analysis recovery capabilities resolver policy approvals transitions executionGate broker ledger dispatch
            autoLayout tb 240 160
        }
        component backend "C3-DispatchAndSettlement" "Reviewed V1 dispatch/settlement: committed work only, deduplicated attempts, conservative settlement, bounded automatic reconciliation and deterministic safe-retry eligibility. Retry re-enters governed control; HITL/manual handling is a last resort, not the default for uncertainty." {
            include executionGate dispatch approvals ledger transactions adapters reconciliation intake recovery
            autoLayout tb 240 160
        }
        component backend "C3-EvidenceAndProgress" "Reviewed V1 evidence path: raw observations remain untrusted, only Verification can publish FACT lifecycle outcomes, INCONCLUSIVE/conflicts remain explicit, fact history is append-only, and deterministic done_when rules over committed valid facts decide completion without bypassing a slow verifier." {
            include analysis adapters intake evidence policy verifier transactions progress scheduler
            autoLayout tb 240 160
        }
        component backend "C3-StopAndRecovery" "Reviewed V1 scoped recovery: PHASE or MISSION stop, one bounded recovery episode per scope/epoch, scoped fences/cancellation, dependency blocking, deterministic phase-to-mission escalation, governed compensation and explicit human resume only." {
            include entry policy transitions mission phases scheduler recovery llmGate executionGate adapters exports ledger reconciliation transactions resolver
            autoLayout tb 240 160
        }
        component backend "C3-DerivedViews" "Reviewed V1 derived state: committed canonical events replay into rebuildable SQL/knowledge-graph/observability views with watermarks; graph supports relationship traversal and task-scoped subgraphs; UI and agents read only through authorized layers; projections never own authority or accounting." {
            include web entry transactions projections readModels readAccess policy broker telemetry
            autoLayout tb 240 160
        }
        component backend "C3-Reporting" "Reviewed V1 reporting: operator/system-selected pinned template drives structure; LLM produces draft content only; authoritative severity comes from governed assessment rules/contracts; verified facts/explicit unknowns populate the template; redaction precedes durable freeze; export authorization binds exact bytes/template/destination and uncertain delivery is reconciled before retry." {
            include operator web analysis broker reports redaction policy approvals transitions exports recovery reconciliation transactions
            autoLayout tb 240 160
        }

        styles {
            element "Element" {
                color "#172B4D"
                background "#EAF1FA"
                stroke "#57769D"
                strokeWidth 2
                fontSize 22
            }
            element "Person" {
                shape Person
                background "#17365D"
                color "#FFFFFF"
            }
            element "Container" {
                shape RoundedBox
                width 400
            }
            element "Component" {
                shape RoundedBox
                width 360
                fontSize 20
            }
            element "Backend" {
                background "#DDEBDD"
                stroke "#527652"
            }
            element "Client" {
                background "#E7EFFA"
            }
            element "Web Client" {
                shape WebBrowser
            }
            element "Execution" {
                background "#E8F1E8"
                stroke "#5F7F5F"
            }
            element "Verification" {
                background "#E8EEF8"
                stroke "#5E7396"
            }
            element "Data Store" {
                shape Cylinder
                background "#FFF0D5"
                stroke "#A98536"
            }
            element "Derived" {
                background "#F2E8F7"
                stroke "#9471AA"
            }
            element "Reasoning" {
                background "#EDE4FA"
                stroke "#9576B5"
            }
            element "Group" {
                color "#52657D"
                stroke "#AAB8C8"
                fontSize 24
            }
            relationship "Relationship" {
                color "#576A82"
                thickness 2
                fontSize 18
                dashed false
            }
            relationship "Check" {
                color "#B27929"
                dashed true
            }
        }
    }

    // IMPLEMENTATION CONTRACTS - APPLY IN THE SOURCE'S ORDER.
    // These are behavioral requirements, not additional C4 containers.
    // A diagram documents them; it does not implement or prove them.

    // 1. SCOPED STOP AND RECOVERY
    // STOP scope is exactly PHASE or MISSION.
    // Human scoped STOP/RESUME uses Entry -> Policy -> Router exclusively.
    // Router coordinates Recovery's stop/fence mutation and the owning Mission/Phase
    // Controller's lifecycle mutation in one transaction before success/notification.
    // Scheduler -> Recovery remains the V1 internal deterministic breaker safety path;
    // it cannot accept human requests or authorize resume/compensation.
    // Request denial, ordinary task failure and budget exhaustion are not
    // automatically STOP; a separate explicit stop rule/decision is required.
    // Stop epoch and recovery-entry uniqueness are enforced atomically per scope.
    // Repeated STOP for the same scope/epoch is idempotent and cannot create
    // another recovery episode.
    // PHASE STOP freezes/cancels normal work in that phase only. Independent phases
    // may continue; phases with required dependencies on the stopped phase become
    // BLOCKED and cannot bypass the missing dependency.
    // MISSION STOP freezes normal admission across the mission.
    // Each active scoped stop may own at most one bounded recovery episode.
    // V1 overlap: Mission STOP atomically fences admission/claims of phase recovery.
    // Reconcile running phase attempts within existing bounds, then atomically
    // transfer remaining effects, unresolved attempts and their history to the
    // mission episode and close phase episodes with recovery_status=PARTIAL,
    // terminal_reason=SUPERSEDED_BY_MISSION_STOP and
    // handover_episode_id=mission_recovery_episode_id. SUPERSEDED is not a status.
    // Unresolved effects/attempts remain in reconciliation and block new compensation
    // for that effect; handover never proves safety. Stable effect_id is preserved
    // across episodes with one compensation owner, checked atomically at dispatch.
    // Handover preserves completed compensation, denials and cumulative limits;
    // it cannot repeat a completed effect or revive a denied operation path.
    // No phase recovery is admitted while mission STOP remains active, including
    // after a later phase STOP or a terminal mission recovery episode.
    // Recovery first reconciles unresolved/in-flight work; uncertain effect is not
    // blindly retried or compensated.
    // Compensation is allowed only for confirmed effects and registered
    // compensation inside the stopped scope, subject to current policy, limits and
    // the governed Resolver -> Policy -> Router -> Gate -> Dispatch -> Runner path.
    // Recovery Controller never directly calls cleanup tools.
    // A hard-denied operation ends that automatic path; materially equivalent
    // compensation may not be rewritten/resubmitted to bypass Policy. The episode
    // may consider genuinely different compensation as a new TypedRequest through
    // Governance only while active and within its existing cumulative limits.
    // Recovery Controller closes on completion, explicit episode-level termination
    // (including mission takeover),
    // episode-limit exhaustion, or no valid automatic continuation; operation-level
    // denial alone does not close an episode. COMPLETE requires no outstanding work.
    // Recovery retries/investigation are bounded by attempts, budget, deadline,
    // current policy and retry contract. MANUAL_REQUIRED is a last-resort terminal
    // fallback, not the default response to uncertainty.
    // PHASE STOP may escalate to MISSION STOP only through an explicit deterministic
    // escalation rule/decision; escalation is never inferred from ordinary failure.
    // Closing recovery atomically invalidates outstanding recovery admission claims
    // for that scope/epoch. Terminal events cannot create a new episode.
    // Terminal recovery status is COMPLETE, PARTIAL or MANUAL_REQUIRED with
    // outstanding work preserved. Late evidence/settlement remains processable.
    // Explicit human resume is AuthorizedResume bound to the expected scope/stop
    // epoch. Recovery checks current canonical state, requires terminal scoped recovery
    // and no unfinished handover, and rejects phase resume while Mission STOP is active.
    // These checks serialize with concurrent STOP on both phase and mission fences.
    // Resume advances the fence generation; old claims never become valid again.
    // Replayed/stale resume cannot clear a newer STOP. Mission resume does not clear
    // explicit phase stops; future admission still checks dependencies/current policy.
    // Resume never resets consumed scope, attempts or recovery accounting, and
    // recovery completion alone never resumes normal work.

    // 2. RESERVATION SETTLEMENT
    // Dispatch claims, reservation and grant-use claims share one transaction.
    // Cancellation and dispatch race through the same authoritative admission fence.
    // Release before dispatch requires proof that late dispatch is impossible.
    // Known usage is consumed; only the proven unused reserved balance is released.
    // Unknown usage/effect remains held or is conservatively accounted; timeout/TTL
    // alone never proves zero usage or that an external effect did not occur.
    // Settlement is idempotent per reservation; late results cannot settle twice.
    // Reconciliation performs bounded automatic investigation before manual fallback.
    // It classifies the outcome and deterministic retry eligibility; it never blindly
    // retries an unknown external effect.
    // A permitted retry keeps operation_id and creates a distinct attempt_id, then
    // re-enters the governed scheduling/policy/gate/dispatch path.
    // Runner never creates retry attempts itself. LLM/Analysis may propose a next
    // action, but cannot authorize or declare an external-effect retry safe.
    // Additional reconciliation diagnostics, if needed, are governed proposals and
    // cannot be called directly by Reconciliation; use FollowUpRequested below.
    // MANUAL_REQUIRED/UNRESOLVED is a last resort after bounded attempts/deadline/
    // budget are exhausted, safe automated continuation cannot be established,
    // a material contract change needs new authority, or policy requires HITL.
    // Consumed scope and attempts are cumulative and are not refunded by cleanup.

    // 3. TRUSTED IDENTITY
    // Backend attaches principal, acting agent, delegation and authz context.
    // Agent-supplied identity or role claims never grant permissions.
    // Requests include targets, capability, parameters, evidence and done_when.
    // Operation identity is stable across retries; attempt identity is distinct.
    // A parameter digest alone is not an operation identity.

    // 4. PROGRESS CONTRACT
    // Separate pure completion, no-progress and state-transition functions.
    // Completion/phase completion uses deterministic criteria over committed state;
    // agent or LLM claims that work is "done" are never authoritative.
    // Completion requiring verification uses committed verifier facts.
    // Scheduler admits already-authorized task contracts and never invents new tasks.
    // Multiple tasks may run concurrently only within bounded concurrency/backpressure.
    // A single task failure does not automatically fail its phase; explicit phase
    // criticality/completion criteria decide the phase outcome.
    // Retry is bounded by retry contract, attempt count, budget, deadline, current
    // policy and no-progress rules; there is no infinite retry loop.
    // Same logical work keeps its logical identity; each permitted attempt is distinct.
    // Unrelated graph updates and wording changes cannot reset no-progress counters.
    // Reducer uses task revision checks and deduplicates consumed events.
    // Durable committed state is the restart source; memory queues are not authority.
    // Scheduler and recovery consumers filter event types and terminal states.
    // Repeated events cannot restart work or retry hard denial.
    // FollowUpRequested is the shared durable proposal from Reconciliation or
    // Verification, with reason RECONCILIATION_DIAGNOSTIC or EVIDENCE_NEEDED,
    // proposal_id, scope, originating task/case and operation/attempt/evidence refs.
    // Existing outbox delivery resolves one current owner from canonical state:
    // Phase Control for phase work, Mission Control for mission work, or Recovery
    // for work assigned to an active episode. Normal work in a stopped scope stays
    // blocked; a proposal cannot grant recovery admission or reopen a terminal episode.
    // The owner atomically records proposal consumption with a blocked/rejected
    // outcome or durable new request; redelivery cannot lose or duplicate that request.
    // Accepted diagnostics follow Resolver -> Policy -> Router -> Gate -> Dispatch.
    // A new request retains originating task/case/episode limits and no-progress
    // history; it cannot reset attempts, budgets or deadlines. New diagnostic work
    // has its own operation_id; an eligible retry retains the existing operation_id.
    // Denial, ordinary task failure and budget exhaustion do not automatically
    // trigger PHASE or MISSION STOP; a separate explicit scoped-stop rule decides so.

    // 5. EVIDENCE AND VERIFICATION
    // Raw observations/evidence never become verified facts or task completion
    // merely because they were produced by an Agent or Runner.
    // Evidence Intake validates schema, identity, provenance and integrity but does
    // not adjudicate truth.
    // Intake first persists a pending receipt with stable evidence_id, expected
    // artifact reference/hash and verification contract, then writes immutable bytes.
    // After durable bytes are confirmed, one state-store transaction finalizes the
    // evidence reference/hash and VerificationRequested job/outbox; only then ACK.
    // Restart resumes pending receipts by checking their referenced bytes. Missing
    // bytes await evidence redelivery under the same identity, never tool re-execution;
    // the same evidence_id with a different hash is rejected. No cross-store ACID
    // transaction or separate queue is required.
    // verification_job_id binds task/scope, immutable evidence refs and rule/version.
    // Verifier owns bounded durable job claims and recovery of unfinished jobs.
    // Terminal job state, FACT/INCONCLUSIVE and any FollowUpRequested commit together
    // before completion use; identity/revision checks deduplicate redelivery and
    // reject stale workers. A changed evidence/rule version requires a new job.
    // Only Verification can publish FACT lifecycle outcomes.
    // Verifier has no direct target/tool capability in V1; additional evidence needs
    // a governed proposal/execution path.
    // VERIFIED, INVALIDATED and INCONCLUSIVE are explicit outcomes.
    // Conflicting evidence remains explicit; LLM confidence or majority vote is not
    // verdict authority.
    // Fact history is append/lifecycle based; newer evidence does not overwrite or
    // delete historical provenance. Stale facts remain historical but cannot satisfy
    // a current rule that requires valid/fresh facts.
    // Deterministic done_when may require multiple committed valid facts.
    // Slow/backlogged verification causes pending/blocked progress as defined by
    // workflow rules; raw evidence is never a temporary verification bypass.
    // INVALIDATED/INCONCLUSIVE does not self-trigger execution retry.

    // 6. REASONING RUNTIME CONTRACT
    // Reasoning workers receive granted, minimal, task-scoped context only; they do
    // not query unrestricted canonical/derived stores directly.
    // All model calls pass through LLM Call Gate. Provider/model selection and API
    // credentials are system/operator controlled; workers cannot choose arbitrary
    // models/providers or access the credential as task authority.
    // Worker profiles are system-defined and bounded; workers cannot self-spawn.
    // Skills are attached only for the current authorized task/profile. Skills guide
    // reasoning but never grant tools, state writes or authorization.
    // Private worker long-term memory is not authoritative. Durable knowledge enters
    // canonical/derived state only through governed contracts with provenance.
    // Workers may propose new tasks but cannot schedule them directly.
    // Workers may report claims/results but cannot decide completion or verdict.
    // Control-affecting outputs are typed contracts; free-form text is not implicit
    // authority. Multiple workers disagreeing or voting does not create authority.
    // Reasoning/no-progress attempts are bounded by budget, deadline and workflow
    // rules; workers cannot continue indefinitely.
    // Mission Intake / Planning is a fixed system profile for questions,
    // clarification and plan proposals. Its outputs are typed interaction contracts,
    // not direct mission/task transitions or execution authority.

    // 7. INTERACTIVE OPERATOR WORKFLOW
    // Web UI is an interactive client, not an authority boundary.
    // Authenticated user turns are classified and routed by Entry; user-supplied
    // role/scope/identity claims remain untrusted until Backend validates them.
    // Mission Controller owns durable DRAFT/AWAITING_INPUT/PLANNING/
    // AWAITING_CONFIRMATION/READY interaction lifecycle before RUNNING as applicable.
    // Missing required input yields ClarificationRequest; the system does not guess.
    // Clarification answers and PlanConfirmation are ordinary mission interaction,
    // not Approval grants and not permission expansion.
    // PlanConfirmation permits Control to submit the confirmed plan to current
    // Policy/transition checks; it never authorizes diagnostic execution by itself.
    // Mid-mission questions and SteeringRequests may invoke bounded reasoning, but
    // resulting SteeringProposal is only a proposal for Mission/Phase Control.
    // UI visualization/query access goes Entry -> Authorized Read Layer -> derived
    // views; UI never connects directly to SQL/Knowledge Graph/authoritative stores.
    // Operator-facing reasoning/control output is typed; free-form chat rendering
    // cannot be parsed back as hidden authority.

    // 8. READ AUTHORIZATION, DERIVED VIEWS AND CACHE
    // Views enforce data access; Broker checks task delegation and minimizes context.
    // Both use one policy source with current authorization revision.
    // Projection Workers consume committed canonical events only and write derived
    // representations only; no projection may mutate canonical/authority state.
    // SQL, Knowledge Graph and optional observability projections are replayable and
    // rebuildable from canonical history.
    // Deterministic/idempotent ordering is required for the relevant entity/stream;
    // global total ordering is not required unless a projection contract needs it.
    // Every projection exposes freshness/version/watermark information.
    // A stale projection may be used for analysis only when task freshness rules
    // permit it; fresh or stale projections never authorize execution, grants,
    // reservations, budget mutation, current stop epoch or verified-fact publication.
    // Knowledge Graph is explicitly used for relationship traversal, asset/entity
    // linkage, evidence/finding provenance, dependency analysis and authorized
    // task-scoped subgraph retrieval.
    // Cache revocation/authz revision is checked before reuse, including report
    // context; TTL alone does not preserve access after permission changes.
    // UI and reasoning consumers never query derived stores directly.

    // 9. APPROVAL DIGEST AND DEDUPLICATION
    // Audit digest records the full request version.
    // Approval digest covers the canonical security-relevant frozen contract.
    // Bind principal, delegation, mission, operation, targets, exact adapter,
    // effectful parameters, limits, preconditions and destination as applicable.
    // Exclude only schema-declared non-executing metadata.
    // Never use LLM judgments to decide that a change is insignificant.
    // A changed material contract requires a new authorization decision.
    // Equal digests do not bypass expiry, revocation or operation-use limits.
    // Pending reviews deduplicate within the same logical operation and context.
    // Fixed batches enumerate exact approved operations; they are not open-ended
    // phase grants. Reviewer and Operator may be the same human only when the
    // current authority model permits both roles; HITL never overrides hard bounds.

    // 10. REVIEWED V1 OPERATION AUTHORIZATION
    // Analysis/LLM proposes capability and target; it does not choose the final
    // tool/adapter/version or obtain execution authority.
    // Capability Resolver deterministically selects only registered adapters and
    // freezes the exact operation before approval; unsupported/unavailable is typed.
    // There is no fallback to an arbitrary or unregistered executable.
    // Policy decides; Authorized Transition Router commits/routes allowed transitions.
    // Required authoritative dependency failures fail closed.
    // Execution Gate rechecks current permission/grant, budget/limits, freshness,
    // material preconditions and current normal/recovery admission immediately
    // before dispatch.
    // Stale/materially changed inputs are blocked; material contract changes require
    // new authorization rather than silent tool/parameter substitution.
    // Hard denial cannot be bypassed by automatically rewriting a materially same
    // request until it passes policy.
    // PHASE STOP freezes normal execution admission only in the affected phase
    // scope and blocks dependants while independent phases may continue.
    // MISSION STOP freezes normal execution admission across the mission.
    // Only explicitly admitted recovery work may proceed inside an active stop scope.
    // No durably committed dispatch acknowledgement means no Runner execution.

    // 11. TEMPLATE-DRIVEN REPORTING
    // Report template selection is owned by Operator, mission configuration or
    // system policy. Report/LLM workers cannot choose, switch or relax the template.
    // The selected template_id/template_version is pinned for one report build.
    // A template defines presentation/schema requirements only; it is never a fact,
    // verdict, severity or authorization source.
    // Required factual fields must be backed by authorized verified facts or remain
    // explicitly UNKNOWN/UNVERIFIED/INCONCLUSIVE when the template permits it.
    // Report Assembly cannot invent facts or modify verifier verdicts to fill fields.
    // Authoritative severity/risk values originate from an authorized assessment
    // contract or deterministic/domain rule over authorized facts. LLM may propose
    // severity/rationale but cannot publish the authoritative rating.
    // Assembly validates required sections/fields/provenance against the pinned
    // template before Output Redaction and Validation.
    // Redaction occurs before artifact freeze. The final rendered output must still
    // conform to the pinned template/version after redaction.
    // Exact frozen bytes plus template_id/template_version/digest must be durably
    // recoverable before approval can become pending; physical storage is deferred.
    // Export authorization and durable release-attempt records bind the exact artifact
    // digest, pinned template_id/version and destination.
    // The same release_attempt_id is deduplicated/idempotent. Unknown, partial or
    // unacknowledged delivery goes to Reconciliation; the system never blindly sends
    // the same external release again. A new release attempt requires deterministic
    // eligibility plus current Policy/approval/export checks.
    // Any material content, template or destination change requires new authorization.
    // Released bytes must exactly match the authorized frozen artifact.
}
