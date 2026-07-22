* Add `MsgAssociateRegistryClass` to the registry module, allowing an existing legacy registry
  entry to be upgraded to use a registry class without modifying roles. Authorization accepts any
  of: CONTROLLER role address, scope data-owner party (for Provenance Metadata Scopes), or NFT
  owner — neither takes precedence.
