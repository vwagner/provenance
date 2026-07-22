package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/provenance-io/provenance/x/registry/types"
)

// HasNFT checks if an NFT exists in either the metadata or nft module.
// If the assetClassId is a metadata scope, it will check if the scope exists.
// Otherwise, it will check if the NFT exists in the nft module.
func (k Keeper) HasNFT(ctx context.Context, assetClassID, nftID *string) bool {
	metadataAddress, isMetadataScope := types.MetadataScopeID(*nftID)
	if isMetadataScope {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		_, found := k.MetadataKeeper.GetScope(sdkCtx, metadataAddress)
		return found
	}
	return k.NFTKeeper.HasNFT(ctx, *assetClassID, *nftID)
}

// ValidateNFTExists returns nil if the described NFT exists, or an error otherwise.
func (k Keeper) ValidateNFTExists(ctx context.Context, assetClassID, nftID *string) error {
	hasNFT := k.HasNFT(ctx, assetClassID, nftID)
	if !hasNFT {
		return types.NewErrCodeNFTNotFound(*nftID)
	}
	return nil
}

// AssetClassExists checks if an asset class exists in either the metadata or nft module.
func (k Keeper) AssetClassExists(ctx context.Context, assetClassID *string) bool {
	metadataAddress, isMetadataScope := types.MetadataScopeSpecID(*assetClassID)
	if isMetadataScope {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		_, found := k.MetadataKeeper.GetScopeSpecification(sdkCtx, metadataAddress)
		return found
	}
	return k.NFTKeeper.HasClass(ctx, *assetClassID)
}

// GetNFTOwner returns the owner of an NFT.
// If the assetClassId is a metadata scope, it will return the owner of the scope.
// Otherwise, it will return the owner of the NFT from the nft module.
func (k Keeper) GetNFTOwner(ctx context.Context, assetClassID, nftID *string) sdk.AccAddress {
	metadataAddress, isMetadataScope := types.MetadataScopeID(*nftID)
	if isMetadataScope {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		// Use the value owner address as the owner of the scope.
		accAddr, err := k.MetadataKeeper.GetScopeValueOwner(sdkCtx, metadataAddress)
		if err != nil {
			return nil
		}
		return accAddr
	}
	return k.NFTKeeper.GetOwner(ctx, *assetClassID, *nftID)
}

// validateAssociateRegistryClassSigner checks whether signer is authorized at the NFT-ownership
// tier for AssociateRegistryClass. For Metadata Scopes the signer must be one of the scope's
// data-owner parties; for all other NFTs the signer must be the value/NFT owner.
func (k Keeper) validateAssociateRegistryClassSigner(ctx context.Context, assetClassID, nftID *string, signer string) error {
	metadataAddress, isMetadataScope := types.MetadataScopeID(*nftID)
	if isMetadataScope {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		scope, found := k.MetadataKeeper.GetScope(sdkCtx, metadataAddress)
		if !found {
			return types.NewErrCodeNFTNotFound(*nftID)
		}
		for _, party := range scope.Owners {
			if party.Address == signer {
				return nil
			}
		}
		return types.NewErrCodeUnauthorized("signer is not a data owner of the scope")
	}
	return k.ValidateNFTOwner(ctx, assetClassID, nftID, signer)
}

// Returns an error if owned by someone else, or if the NFT doesn't exist.
func (k Keeper) ValidateNFTOwner(ctx context.Context, assetClassID, nftID *string, expOwner string) error {
	nftOwner := k.GetNFTOwner(ctx, assetClassID, nftID)
	if len(nftOwner) == 0 || nftOwner.String() != expOwner {
		return types.NewErrCodeUnauthorized("signer does not own the NFT")
	}
	return nil
}
