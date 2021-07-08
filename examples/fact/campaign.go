package fact

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/buger/jsonparser"

	"gitlab.id.vin/gami/gami-common/adapters/group"
	"gitlab.id.vin/gami/gami-common/logger"
	"gitlab.id.vin/gami/gami-common/models"
	"gitlab.id.vin/gami/gami-service/utils"
)

const (
	RateLimited001            = "RATE_LIMIT_001"   // limit device_id per day
	RateLimited002            = "RATE_LIMIT_002"   // limit device_id per campaign
	LimitConditionByDeviceID  = "LIMIT_REWARD_002" // limit by device id
	RateLimitPerDayType       = "PER_DAY"
	RateLimitByCampaign       = "PER_CAMPAIGN"
	DeviceIdFieldMapping      = "user_attributes.audience.device_id"
	LimitRewardByCondition001 = "LIMIT_REWARD_001"
	LimitRewardByLocation001  = "LIMIT_LOCATION_001"
)

type CampaignFact struct {
	Base
	Redeem  *models.Redeem
	Message json.RawMessage
}

func NewCampaignFact(redeem *models.Redeem) *CampaignFact {
	return &CampaignFact{
		Redeem: redeem,
	}
}

func (f *CampaignFact) CheckDeviceIDPerDay(campaignID int64, limit int64) bool {
	ctx := context.Background()
	extraData, err := f.Redeem.ExtraData.MarshalJSON()
	if err != nil {
		logger.Context(ctx).Errorf("unmarshall extra data of %v error : %v", f.Redeem.UserID, err)
		return false
	}

	deviceID, err := jsonparser.GetString(extraData, strings.Split(DeviceIdFieldMapping, ".")...)
	if err != nil {
		logger.Context(ctx).Errorf("get device_id error: user_id %v", f.Redeem.UserID, err)
		return false
	}

	yy, mm, dd := time.Now().In(f.Location).Date()
	dateStr := fmt.Sprintf("%02d%02d%04d", dd, mm, yy)

	isExists := f.CacheAdapter.Exists(ctx, utils.GetRateLimitByDay(campaignID, dateStr))
	if !isExists {
		if err := f.CacheAdapter.Expire(ctx, utils.GetRateLimitByDay(campaignID, dateStr), time.Hour*24); err != nil {
			logger.Context(ctx).Errorf("set expire of key %v error : %v", deviceID, err)
			return false
		}
	}

	currentValue, err := f.CacheAdapter.HIncr(ctx, utils.GetRateLimitByDay(campaignID, dateStr), deviceID, 1)
	if err != nil {
		logger.Context(ctx).Errorf("increment key error : %v", err)
		return false
	}

	if currentValue > limit {
		logger.Context(ctx).Errorf("[PROMOTION_LOG][HandlePassRule][RateLimitedPerDay] user_id: %v - device_id: %v - campaign_id: %v", f.Redeem.UserID, deviceID, campaignID)
		return false
	}

	return true
}

func (f *CampaignFact) CheckDeviceIDByCampaign(campaignID int64, limit int64) bool {
	ctx := context.Background()
	extraData, err := f.Redeem.ExtraData.MarshalJSON()
	if err != nil {
		logger.Context(ctx).Errorf("unmarshall extra data of %v error : %v", f.Redeem.UserID, err)
		return false
	}

	deviceID, err := jsonparser.GetString(extraData, strings.Split(DeviceIdFieldMapping, ".")...)
	if err != nil {
		logger.Context(ctx).Errorf("get device_id error: user_id %v", f.Redeem.UserID, err)
		return false
	}
	currentValue, err := f.CacheAdapter.HIncr(ctx, utils.GetRateLimitByCampaign(campaignID), deviceID, 1)
	if err != nil {
		logger.Context(ctx).Errorf("increment key device_id %v in redis error : %v", deviceID, err)
		return false
	}

	if currentValue > limit {
		logger.Context(ctx).Errorf("[PROMOTION_LOG][HandlePassRule][RateLimitedCampaign] user_id: %v - device_id: %v - campaign_id: %v", f.Redeem.UserID, deviceID, campaignID)
		return false
	}

	return true
}

func (f *CampaignFact) LimitByDeviceID(groupID int64) bool {
	ctx := context.Background()
	extraData, err := f.Redeem.ExtraData.MarshalJSON()
	if err != nil {
		logger.Context(ctx).Errorf("unmarshall extra data of %v error : %v", f.Redeem.UserID, err)
		return true
	}

	deviceID, err := jsonparser.GetString(extraData, strings.Split(DeviceIdFieldMapping, ".")...)
	if err != nil {
		logger.Context(ctx).Errorf("get device_id error: user_id %v", f.Redeem.UserID, err)
		return true
	}

	req := &group.IsInGroupsRequest{
		Value:     deviceID,
		GroupIDs:  []int64{groupID},
		GroupType: group.DataGroupType,
	}

	res, err := f.GroupAdapter.IsInGroups(ctx, req)
	if err != nil {
		logger.Context(ctx).Errorf("verify device id in group for request %v error : %v", req, err)
		return true
	}
	if res.Meta.Status != 200 {
		logger.Context(ctx).Errorf("verify device id in group for request %v has resp is invalid", req, res)
		return true
	}

	return res.Data.IsInGroups
}

func (f *CampaignFact) Get(path string) interface{} {
	value, _, _, err := jsonparser.Get(f.Message, strings.Split(path, ".")...)
	if err != nil {
		return ""
	}

	var rt interface{}
	_ = json.Unmarshal(value, &rt)
	return rt
}

func (f *CampaignFact) Set(path string, value interface{}) {
	bytes, err := json.Marshal(value)
	if err != nil {
		f.Errors = append(f.Errors, err)
		return
	}

	_, err = jsonparser.Set(f.Message, bytes, strings.Split(path, ".")...)
	if err != nil {
		f.Errors = append(f.Errors, err)
	}
}
