package fact

import (
	"context"
	"strconv"
	"strings"
	"time"

	"gitlab.id.vin/gami/gami-common/adapters/cache"
	"gitlab.id.vin/gami/gami-common/adapters/database"
	"gitlab.id.vin/gami/gami-common/adapters/group"
	"gitlab.id.vin/gami/gami-common/logger"
	"gitlab.id.vin/gami/gami-service/adapter/voucher_ordering"
)

type Base struct {
	Time
	VoucherOrderingAdapter voucher_ordering.VoucherOrderingAdapter
	GroupAdapter           group.Adapter
	CacheAdapter           cache.CachedV2Adapter
	DBAdapter              database.DBAdapterV2
	Errors                 []error
}

func (f *Base) CheckUserGroup(userID, groupID int64) bool {
	res, err := f.GroupAdapter.IsInGroups(context.Background(), &group.IsInGroupsRequest{
		Value:     groupID,
		GroupIDs:  []int64{groupID},
		GroupType: group.UserGroupType,
	})
	if err != nil || res.Meta.Status != 200 {
		return false
	}

	return res.Data.IsInGroups
}

func (f *Base) CheckDataGroup(data interface{}, groupID int64) bool {
	res, err := f.GroupAdapter.IsInGroups(context.Background(), &group.IsInGroupsRequest{
		Value:     data,
		GroupIDs:  []int64{groupID},
		GroupType: group.DataGroupType,
	})

	if err != nil || res.Meta.Status != 200 {
		return false
	}

	return res.Data.IsInGroups
}

func (f *Base) GetUserLocation(id int64) string {
	ctx := context.Background()
	userID := strconv.FormatInt(id, 10)
	merchantInfo, err := f.VoucherOrderingAdapter.GetMerchantInfo(ctx, userID)

	if err != nil {
		logger.Context(ctx).Errorf("get merchant information for use_id %v got error : %v", userID, err)
		return ""
	}
	if merchantInfo.Meta.Code != 200 {
		logger.Context(ctx).Errorf("get merchant information for user_ud %v has invalid response %#v", userID, merchantInfo)
		return ""
	}
	if merchantInfo.Data.Address == nil {
		logger.Context(ctx).Errorf("get merchant information for user %v has invalid merchant address: %v", userID, merchantInfo)
		return ""
	}

	return strings.TrimSpace(merchantInfo.Data.Address.ProvinceCode)
}

type Time struct {
	Location *time.Location
}
