# Google Cloud IAM — Session Notes

> Source: [Google Cloud IAM Documentation](https://docs.cloud.google.com/iam/docs/overview)

---

## 1. What is IAM?

IAM (Identity and Access Management) is the tool that manages fine-grained authorization for Google Cloud. It controls **who** can do **what** on **which resources**.

Every action in Google Cloud requires certain permissions. When someone tries to perform an action — for example, create a VM or view a dataset — IAM first checks if they have the required permissions. If they don't, IAM prevents them from performing the action.

---

## 2. The Core Model: Principal + Role + Resource

Granting permissions in IAM always involves three components:

| Component | Description |
|-----------|-------------|
| **Principal** | The identity of the person or system you want to give permissions to |
| **Role** | The collection of permissions you want to grant |
| **Resource** | The Google Cloud resource you want to let the principal access |

---

## 3. Principals — Who Can Act?

There are two broad categories of principals:

- **Human users**: Google Accounts, Google Groups, federated identities (workforce identity pools)
- **Workloads**: Service accounts, federated identities (workload identity pools)

> A **service account** is a principal in its own right — it represents a program or service, not a human.

---

## 4. Permissions and Roles

### Permission format

Permissions follow the format `service.resource.verb`. For example:
- `resourcemanager.projects.list` → list Resource Manager projects
- `storage.objects.get` → read an object from Cloud Storage

> **You cannot grant permissions directly to a principal.** You must go through roles.

### The 3 types of roles

| Type | Description | Use case |
|------|-------------|----------|
| **Predefined roles** | Managed by Google, scoped to a service | Production — e.g. `roles/pubsub.publisher` |
| **Custom roles** | Defined by you, exact permissions you choose | Least-privilege fine-tuning |
| **Basic roles** | Very broad (Owner, Editor, Viewer) | Testing only — never in production |

---

## 5. Allow Policies — the Technical Mechanism

You grant roles to principals using **allow policies**: a YAML or JSON object attached to a Google Cloud resource, containing a list of **role bindings**.

Example JSON structure:

```json
{
  "bindings": [
    {
      "role": "roles/storage.objectViewer",
      "members": [
        "user:alice@example.com",
        "serviceAccount:my-app@project.iam.gserviceaccount.com"
      ]
    },
    {
      "role": "roles/storage.objectAdmin",
      "members": [
        "group:devs@example.com"
      ]
    }
  ]
}
```

---

## 6. Where to Apply a Policy — Project vs. Resource Level?

**You do NOT need to apply a policy to each resource one by one.**

IAM uses **policy inheritance**: if you set a policy on a container resource (Organization, Folder, Project), it automatically applies to all resources inside that container.

### The 3 levels to choose from:

| Level | Scope | When to use |
|-------|-------|-------------|
| **Organization** | All resources in the org | Cross-team admin access, security baseline |
| **Folder** | All projects in the folder | Team or department-wide access |
| **Project** | All resources in the project | Most common — dev vs prod isolation |
| **Resource** | One specific resource | Fine-grained least-privilege (one bucket, one topic…) |

### Key rule

> To understand who can access a resource, you must look at **all** the allow policies that affect it — from the resource up to the organization root. The union of all these policies is called the **effective allow policy**.

---

## 7. Policy Inheritance — How It Works

```
Organization  ← policy set here
    ↓ inherited by
  Folder       ← policy set here overrides/adds to org policy
    ↓ inherited by
    Project    ← policy set here overrides/adds to folder policy
      ↓ inherited by
      Resource ← effective allow policy = union of all above
```

Three practical implications:

1. **One binding, many resources** — grant a role on a project to cover all resources inside it.
2. **Resources without their own policy** — not all resources accept allow policies directly (e.g. log buckets). Grant the role on the parent project instead.
3. **Effective allow policy** — always the union of all ancestor policies. A permission granted at the org level cannot be "removed" by a project-level policy alone — you need a deny policy for that.

---

## 8. IAM Decision Flow

When a principal makes a request, IAM evaluates in this order:

```
1. Principal makes a request
        ↓
2. IAM computes the effective allow policy
   (union of org + folder + project + resource policies)
        ↓
3. Does a role binding grant the required permission?
   → No  → ACCESS DENIED
   → Yes ↓
4. Does a deny policy block the permission?
   → Yes → ACCESS DENIED
   → No  ↓
5. Is the resource within the Principal Access Boundary?
   → No  → ACCESS DENIED
   → Yes ↓
6. Are all IAM Conditions satisfied?
   (time of day, resource type, tags…)
   → No  → ACCESS DENIED
   → Yes ↓
7. ACCESS GRANTED
```

> **Key insight**: a deny policy overrides an allow policy. Even if a role grants a permission, a deny policy can block it.

---

## 9. Advanced Access Control Mechanisms

| Mechanism | What it does |
|-----------|-------------|
| **Deny policies** | Block specific permissions even if a role grants them |
| **Principal Access Boundary (PAB)** | Restrict which resources a principal is eligible to access |
| **IAM Conditions** | Conditional access based on attributes (time, resource type, tags) |
| **Privileged Access Manager (PAM)** | Temporary, auditable, approval-based access to sensitive resources |

---

## 10. Technical Note: Eventual Consistency

The IAM API is **eventually consistent**. If you create a service account or modify a policy, the change may take a few seconds to propagate — IAM might temporarily return stale data or not find the newly created resource.

---

## Summary

> IAM controls access via the formula **Principal + Role + Resource**, enforced through **allow policies** that inherit down the resource hierarchy. Apply policies at the **project or folder level** for broad coverage, and at the **resource level** only for fine-grained least-privilege access. The **effective allow policy** is always the union of all ancestor policies — evaluated top-down at every request.