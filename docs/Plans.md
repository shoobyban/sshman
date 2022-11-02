# Planned functionality for the next releases

## State Handling

I've been using the systems for a while and I've noticed that the current version gets out of sync quite easy.

The idea is to have 3 states: the host's read state (current files on host), the "stable" state (that's in theory we have in config) and a "staging" state (that we are clicking together).

In theory stable and staging will make sense on bigger systems, so we could just have two states.

I'm planning to have a Refresh (read state) and a Publish button / command.

## Edit multiple items

with checkboxes (removing them from master for now) when multiple items are selected an "Edit Selected" button will appear (implemented) by
Add {itemtype} button. The editor will display [ _multiple values_ ] value for non-uniform values, as soon as editing is made (after save) these values
will be updated for every edited item.

## Changeset based operations

Sync to host operations would gather a changeset and applying the changeset would be a separate operation. This would allow us to have an overview of what's going to happen listing in a "dry run" mode way.
  - [ ] should keep a list of todo ops (per server: user add or delete)
  - [ ] display the ops on frontend
  - [ ] ops should be grouped by hosts -> 1 op for host even if many user change
  - [ ] ops for same host-user pair (add + delete) would apply the latest change
  - [ ] apply button should run them, preparing undo op (cache old server authorized_keys files)
  - [ ] undo op to upload cached authorized_keys and restore changeset
  - [ ] on update or new host list new users on frontend and on cli
