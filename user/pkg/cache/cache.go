package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/swarit-pandey/e-commerce/common/cache"
	"github.com/swarit-pandey/e-commerce/user/pkg/config"
	"github.com/swarit-pandey/e-commerce/user/pkg/repository"
	"k8s.io/klog/v2"
)

const (
	cacheSet     = "entity_cache"
	userCacheBin = "user"
	addrCacheBin = "address"
)

type CacheService struct {
	cache         cache.Repository[any]
	userRepo      repository.User
	addrRepo      repository.UserAddress
	userInterface repository.UserRepo
	addrInterface repository.AddressRepo
}

func NewCacheService(conf *config.Cache, userRepo repository.User, addrRepo repository.UserAddress) (*CacheService, error) {
	cacheConfig := cache.Config{
		URL:       conf.Address,
		Port:      conf.Port,
		Namespace: conf.Namespace,
		LogLevel:  conf.LogLevel,
		Set:       cacheSet,
		Bin:       userCacheBin,
	}

	aeroCache, err := cache.NewAero[any](&cacheConfig)
	if err != nil {
		return nil, err
	}

	return &CacheService{
		cache:    aeroCache,
		userRepo: userRepo,
		addrRepo: addrRepo,
	}, nil
}

func (cs *CacheService) SetUser(ctx context.Context, user *repository.User) error {
	createdUser, err := cs.userInterface.Create(ctx, user)
	if err != nil {
		return err
	}

	err = cs.setUserInCache(ctx, *createdUser)
	if err != nil {
		klog.ErrorS(err, "failed to add a new user in cache")
		return err
	}

	return nil
}

func (cs *CacheService) GetUser(ctx context.Context, userID uint) (*repository.User, error) {
	// Scenario 1: Data found in cache
	user, err := cs.getUserFromCache(ctx, userID)
	if err != nil {
		return user, nil
	}

	// Scenario 2: Data not found in cache, get from repository
	user, err = cs.userInterface.Get(ctx, &repository.User{ID: userID})
	if err == nil {
		// Write-through
		err = cs.setUserInCache(ctx, *user)
		if err != nil {
			klog.ErrorS(err, "failed to write back the response onto the cache, cache debouncing might occur")
		}
		return user, nil
	}

	// Scenario 3: Data not found in cache and not in repository
	return nil, ErrEntityNotFound
}

func (cs *CacheService) getUserFromCache(ctx context.Context, userID uint) (*repository.User, error) {
	key := fmt.Sprintf("user_%d", userID)
	value, err := cs.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	user, ok := value.(*repository.User)
	if !ok {
		return nil, ErrCacheInvalidated
	}
	return user, nil
}

func (cs *CacheService) setUserInCache(ctx context.Context, user repository.User) error {
	userIDKey := fmt.Sprintf("user_%d", user.ID)
	usernameKey := fmt.Sprintf("user_username_%s", user.Username)

	err := cs.cache.Set(ctx, userIDKey, user, 2*time.Hour)
	if err != nil {
		return err
	}

	err = cs.cache.Set(ctx, usernameKey, user, 2*time.Hour)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CacheService) GetAddress(ctx context.Context, addrID uint) (*repository.UserAddress, error) {
	// Scenario 1: Data found in cache
	addr, err := cs.getAddressFromCache(ctx, addrID)
	if err != nil {
		return addr, nil
	}

	// Scenario 2: Data not found in cache, get from repository
	addr, err = cs.addrInterface.Get(ctx, &repository.UserAddress{ID: addrID})
	if err == nil {
		// Write-through
		err = cs.setAddressInCache(ctx, *addr)
		if err != nil {
			klog.ErrorS(err, "failed to write back the response onto the cache, cache debouncing might occur")
		}
		return addr, nil
	}

	// Scenario 3: Data not found in cache and not in repository
	return nil, ErrEntityNotFound
}

func (cs *CacheService) SetAddress(ctx context.Context, addr *repository.UserAddress) error {
	// Update the address data in the repository
	updatedAddr, err := cs.addrInterface.Update(ctx, addr)
	if err != nil {
		return err
	}

	// Update the address data in the cache (write-through)
	err = cs.setAddressInCache(ctx, *updatedAddr)
	if err != nil {
		klog.ErrorS(err, "failed to update the cache, cache debouncing might occur")
		return err
	}

	return nil
}

func (cs *CacheService) getAddressFromCache(ctx context.Context, addrID uint) (*repository.UserAddress, error) {
	key := fmt.Sprintf("addr_%d", addrID)
	value, err := cs.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	address, ok := value.(*repository.UserAddress)
	if !ok {
		return nil, ErrCacheInvalidated
	}
	return address, nil
}

func (cs *CacheService) setAddressInCache(ctx context.Context, addr repository.UserAddress) error {
	key := fmt.Sprintf("address_%d", addr.ID)
	err := cs.cache.Set(ctx, key, addr, 2*time.Hour)
	return err
}

func (cs *CacheService) DeleteUser(ctx context.Context, userID uint) error {
	ok, err := cs.userInterface.Exists(ctx, &repository.User{ID: userID})
	if err != nil {
		return err
	}

	if !ok {
		klog.Warning("the object that is trying to be deleted does not exists")
		return nil
	}

	user, err := cs.userInterface.Delete(ctx, &repository.User{ID: userID})
	if err != nil {
		klog.ErrorS(err, "failed to delete user from repository")
		return err
	}

	err = cs.deleteUserFromCache(ctx, user)
	if err != nil {
		klog.ErrorS(err, "failed to delete from cache")
		return err
	}
	return nil
}

func (cs *CacheService) deleteUserFromCache(ctx context.Context, user *repository.User) error {
	userIDKey := fmt.Sprintf("user_%d", user.ID)
	usernameKey := fmt.Sprintf("user_username_%s", user.Username)

	err := cs.cache.Delete(ctx, userIDKey)
	if err != nil {
		return err
	}

	err = cs.cache.Delete(ctx, usernameKey)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CacheService) DeleteAddress(ctx context.Context, addrID uint) error {
	ok, err := cs.addrInterface.Exists(ctx, &repository.UserAddress{ID: addrID})
	if err != nil {
		return err
	}

	if !ok {
		klog.Warning("the object we are trying to delete does not exists")
	}

	address, err := cs.addrInterface.Delete(ctx, &repository.UserAddress{ID: addrID})
	if err != nil {
		klog.ErrorS(err, "failed to delete user address from repository")
		return err
	}

	err = cs.deleteAddressFromCache(ctx, address.ID)
	if err != nil {
		klog.ErrorS(err, "failed to delete user address from cache")
		return err
	}
	return nil
}

func (cs *CacheService) deleteAddressFromCache(ctx context.Context, addrID uint) error {
	key := fmt.Sprintf("addr_%d", addrID)
	err := cs.cache.Delete(ctx, key)
	return err
}

func (cs *CacheService) GetUserByUsername(ctx context.Context, username string) (*repository.User, error) {
	// Scenario 1: Data found in cache
	user, err := cs.getUserFromCacheByUsername(ctx, username)
	if err != nil {
		return user, nil
	}

	// Scenario 2: Data not found in cache, get from repository
	user, err = cs.userInterface.GetByUsername(ctx, user)
	if err == nil {
		// Write-through
		err = cs.setUserInCache(ctx, *user)
		if err != nil {
			klog.ErrorS(err, "failed to write back the response onto the cache, cache debouncing might occur")
		}
		return user, nil
	}

	// Scenario 3: Data not found in cache and not in repository
	return nil, ErrEntityNotFound
}

func (cs *CacheService) getUserFromCacheByUsername(ctx context.Context, username string) (*repository.User, error) {
	key := fmt.Sprintf("user_username_%s", username)
	value, err := cs.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	user, ok := value.(*repository.User)
	if !ok {
		return nil, ErrCacheInvalidated
	}
	return user, nil
}

func (cs *CacheService) SetPasswordResetToken(ctx context.Context, userID uint, token string) error {
	key := fmt.Sprintf("reset_token_%d", userID)
	err := cs.cache.Set(ctx, key, token, 6*time.Hour)
	return err
}
