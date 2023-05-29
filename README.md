- create a cr to take user input, the controller will detect it and it will spin up a cr for snapshot

- to restore, we need create a pv where source has to be kept as snapshot, and then deployment has to modified and run again

- for restore we can do it in the same crd using boolean for now.
