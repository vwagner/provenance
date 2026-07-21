package types

import (
	"errors"
)

// DefaultParams returns the default registry module parameters. By default the module defines no
// role authorization policies, which preserves the original chain behavior: every role change
// falls back to legacy NFT-owner authorization. Governance can install policies via
// MsgUpdateParams, and registry-class maintainers can define per-asset-class policies.
func DefaultParams() Params {
	return Params{}
}

// Validate validates the Params.
func (p Params) Validate() error {
	var errs []error
	errs = append(errs, validateRoleAuthorizations(p.RoleAuthorizations)...)
	if p.PendingChangeExpiry < 0 {
		errs = append(errs, NewErrCodeInvalidField("pending_change_expiry", "must be non-negative"))
	}
	return errors.Join(errs...)
}

// RoleAuthorizationMap returns a map of RegistryRole -> RoleAuthorization for the params' default
// policies, for fast lookup during authorization resolution.
func (p Params) RoleAuthorizationMap() map[RegistryRole]RoleAuthorization {
	return RoleAuthorizationMapFrom(p.RoleAuthorizations)
}
