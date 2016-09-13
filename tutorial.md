# Writing a distributed application in Go using etcd/raft

The intuitive flow for a regular application: a simple cycle between data structures or storage holding the data.
<pre><code>
                    Write
        ┌─────────────────────────────┐
        │                             │
        │                             ▼
 ┌─────────────┐               ┌─────────────┐
 │             │               │             │
 │             │               │             │
 │ Application │               │   Storage   │
 │             │               │             │
 │             │               │             │
 └─────────────┘               └─────────────┘
        ▲                             │
        │                             │
        └─────────────────────────────┘
                     Read
</code></pre>



Raft is takes control of the writing steps of the storage, the application still reads from it - but data is modified by Raft alone.
<pre><code>
                              Propose
       ┌───────────────────────────────────────────────────────┐
       │                                                       │
       │                                                       ▼
┌─────────────┐           ┌─────────────┐               ┌─────────────┐
│             │           │             │               │             │
│             │           │             │     Commit    │             │
│ Application │           │   Storage   │◀──────────────│    Raft     │───────▶ Raft Cluster
│             │           │             │               │             │
│             │           │             │               │             │
└─────────────┘           └─────────────┘               └─────────────┘
       ▲                         │
       │                         │
       └─────────────────────────┘
                   Read
</code></pre>