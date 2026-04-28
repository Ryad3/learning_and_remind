# Google Cloud — Session Notes

> Source: [Google Cloud Official Documentation](https://docs.cloud.google.com/docs)

---

## 1. What is Google Cloud?

Google Cloud is a **cloud platform**: rather than buying and managing physical hardware (servers, disks, networking), you access these resources remotely over the internet, in the form of **services** managed by Google.

---

## 2. The Physical Infrastructure: Universe → Regions → Zones

The platform is built on a three-level geographical hierarchy.

### The Universe
A fully self-contained cloud, with its own networking that is separate from the public internet and other universes. Google Cloud is the original universe, with resources in datacenters all over the world.

### Regions
Geographic areas — Google Cloud has regions in Asia, Australia, Europe, Africa, the Middle East, North America, and South America.

### Zones
Regions are subdivided into zones. For example, zone `a` in the East Asia region is named `asia-east1-a`. Zones have high-bandwidth, low-latency network connections to other zones in the same region.

### Resource Scope Summary

| Scope | Type | Examples |
|-------|------|---------|
| Global | Accessible anywhere in the universe | Preconfigured disk images, snapshots, networks |
| Regional | Accessible within the same region | Static external IP addresses |
| Zonal | Accessible within the same zone | VMs, VM types, disks |

> ⚠️ A disk can only be attached to a machine in the **same zone**. Attaching a disk from one region to a machine in another would introduce too much latency — Google Cloud prevents this.

---

## 3. The "Everything is a Service" Model

In the cloud, what you might think of as software and hardware products become **services**. These services provide access to the underlying resources — from managed Kubernetes to data storage. You build your application by mixing and matching these services.

---

## 4. How to Interact with Google Cloud

There are four main ways:

- **Web console**: graphical interface accessible from a browser
- **gcloud CLI**: command line, e.g. `gcloud compute instances create`
- **Client libraries**: SDKs in Python, Java, Go, Node.js, C#, PHP, Ruby, C++
- **Infrastructure as Code**: Terraform with the Google Cloud provider

---

## 5. The Resource Hierarchy (Organization, Folders, Projects)

### Overview

There are two broad categories of resources:
- **Container resources** (Organization, Folders, Projects) — used to organize and control access
- **Service resources** (VMs, GKE clusters, Pub/Sub topics, etc.) — the actual components of Google Cloud products

### Hierarchy Diagram

```
Organization: mycompany.com              ← single root node
│
├── Folder: Backend Team
│   ├── Project: backend-dev   (ID: bk-dev-001)
│   │   ├── VM (Compute Engine)
│   │   └── Bucket (Cloud Storage)
│   └── Project: backend-prod  (ID: bk-prod-001)
│       └── resources...
│
└── Folder: Data Team
    ├── Project: data-ingestion (ID: data-ing-001)
    │   └── resources...
    └── Project: data-analytics (ID: data-ana-001)
        └── resources...
```

### Level 1 — The Organization (root)

- Root node of the hierarchy, under which all other resources are created
- Access policies applied to the organization are automatically applied to **all** resources below it
- Projects belong to the organization (not to the user who created them) — they survive user deletion

### Level 2 — Folders (optional)

- Create **isolation boundaries** between projects
- Can contain projects and sub-folders
- Typical use cases: one folder per department, per team, or per environment

### Level 3 — Projects (base unit)

- All Google Cloud resources must belong to a project
- Resources within the same project can communicate via the **internal network**
- A project **cannot** access another project's resources (unless using Shared VPC or VPC Network Peering)

Each project has **3 identifiers**:

| Identifier | Who provides it | Note |
|------------|----------------|------|
| Project name | You | Freely chosen |
| Project ID | You or Google | Unique across all of GCP — never reusable after deletion |
| Project number | Google | Auto-generated |

### Level 4 — Service Resources

The actual resources: VMs, Cloud Storage buckets, Cloud SQL databases, GKE clusters, etc. They always live **inside a project**.

---

## 6. Key Concept: IAM Permission Inheritance

```
Organization
    ↓ inherited by
  Folder
    ↓ inherited by
    Project
      ↓ inherited by
      Service Resource
```

If you grant a permission to a user at the **Folder** level, they automatically inherit it on all **Projects** within that folder — no need to redefine it project by project.

---

## 7. App Hub — A Cross-Cutting Layer (Bonus)

The resource hierarchy is useful for organizing resources by team, geography, or compliance. However, cloud applications often combine resources that don't match the hierarchy. For example, an online shop might have UI components in one project and a database in another.

Google Cloud lets you group these resources as **App Hub applications**, without changing the structure of the resource hierarchy — to simplify monitoring, deployment, and troubleshooting.

---

## Summary in One Sentence

> Google Cloud exposes Google's infrastructure (worldwide datacenters, high-bandwidth networking) as on-demand services, organized into **projects** (isolated by default), grouped in **folders**, attached to an **organization**, with IAM permissions inherited top-down.